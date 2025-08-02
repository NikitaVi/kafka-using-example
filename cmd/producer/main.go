package main

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/NikitaVi/microservices_kafka/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"log"
	"strings"
	"time"
)

const (
	brokerAddress = "localhost:9092, localhost:9093, localhost:9094"
	topicName     = "notes_topic"
)

func main() {
	producer, err := newSyncProducer(strings.Split(brokerAddress, ","))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("faield to close producer: %v", err)
		}
	}()

	for i := 0; i < 10; i++ {
		info := model.NoteInfo{
			Title:   gofakeit.Username(),
			Content: gofakeit.BeerAlcohol(),
		}

		data, err := json.Marshal(info)
		if err != nil {
			log.Fatalf("failed to marshal note info: %v", err)
		}

		msg := &sarama.ProducerMessage{
			Topic: topicName,
			Value: sarama.StringEncoder(data),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("failed to send message: %v", err)
			return
		}

		log.Printf("message sent to partition %d at offset %d", partition, offset)
		time.Sleep(time.Second / 2)
	}
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
