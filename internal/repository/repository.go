package repository

import (
	"context"
	"github.com/NikitaVi/microservices_kafka/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	GetById(ctx context.Context, id int64) (*model.Note, error)
}
