package note

import (
	"context"
	"di_container/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	var id int64

	err := s.txManger.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.noteRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.noteRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
