package converter

import (
	"github.com/NikitaVi/microservices_kafka/internal/model"
	modelRepo "github.com/NikitaVi/microservices_kafka/internal/repository/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {
	return &model.Note{
		ID: note.ID,
		Info: model.NoteInfo{
			Title:   note.Title,
			Content: note.Content,
		},
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
