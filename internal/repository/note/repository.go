package note

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/NikitaVi/microservices_kafka/internal/model"
	"github.com/NikitaVi/microservices_kafka/internal/repository"
	"github.com/NikitaVi/microservices_kafka/internal/repository/converter"
	modelRepo "github.com/NikitaVi/microservices_kafka/internal/repository/model"
	"github.com/NikitaVi/platform_shared/pkg/db"
)

const (
	tableName       = "note"
	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.NoteRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING " + idColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "note_repository.Create",
		QueryRow: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetById(ctx context.Context, id int64) (*model.Note, error) {
	builder := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.GetById",
		QueryRow: query,
	}

	var note modelRepo.Note
	err = r.db.DB().ScanOneContext(ctx, &note, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrorNoteNotFound
		}

		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
