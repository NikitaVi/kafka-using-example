package app

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/NikitaVi/microservices_kafka/internal/config"
	"github.com/NikitaVi/microservices_kafka/internal/config/env"
	"github.com/NikitaVi/microservices_kafka/internal/repository"
	noteRepo "github.com/NikitaVi/microservices_kafka/internal/repository/note"
	"github.com/NikitaVi/microservices_kafka/internal/service"
	noteSaverConsumer "github.com/NikitaVi/microservices_kafka/internal/service/consumer"
	"github.com/NikitaVi/platform_shared/pkg/closer"
	"github.com/NikitaVi/platform_shared/pkg/db"
	"github.com/NikitaVi/platform_shared/pkg/db/pg"
	//"github.com/NikitaVi/platform_shared/pkg/kafka"
	"github.com/NikitaVi/microservices_kafka/internal/client/kafka"
	kafkaConsumer "github.com/NikitaVi/microservices_kafka/internal/client/kafka/consumer"
	"log"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	kafkaConsumerConfig config.KafkaConsumerConfig

	dbClient db.Client

	noteRepository repository.NoteRepository

	noteSaverConsumer service.ConsumerService

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPgConfig()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}

		closer.Add(cl.Close)
		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if s.noteRepository == nil {
		s.noteRepository = noteRepo.NewRepository(s.DBClient(ctx))
	}
	return s.noteRepository
}

func (s *serviceProvider) NoteSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.noteSaverConsumer == nil {
		s.noteSaverConsumer = noteSaverConsumer.NewService(
			s.NoteRepository(ctx),
			s.Consumer(),
		)
	}

	return s.noteSaverConsumer
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)

		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)

		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}
