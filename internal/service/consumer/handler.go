package consumer

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/NikitaVi/microservices_kafka/internal/model"
	"log"
)

func (s *service) NoteSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	noteInfo := &model.NoteInfo{}

	log.Println("HERE")

	err := json.Unmarshal(msg.Value, noteInfo)
	if err != nil {
		return err
	}

	id, err := s.noteRepository.Create(ctx, noteInfo)
	if err != nil {
		return err
	}

	log.Printf("Note ID: %d", id)

	return nil
}
