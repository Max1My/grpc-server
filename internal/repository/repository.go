package repository

import (
	"context"

	"di_container/internal/model"
	// desc "di_container/pkg/note_v1"
)

type NoteRepository interface {
	Create(context.Context, *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}

type OtherNoteRepository interface {
	Get(ctx context.Context, id int64) (*model.Note, error)
}
