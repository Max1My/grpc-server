package note

import (
	"context"
	"di_container/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.Note, error) {
	note, err := s.noteRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}
