package note

import (
	"di_container/internal/client/db"
	"di_container/internal/repository"
	"di_container/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
	txManger       db.TxManager
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
		txManger:       txManager,
	}
}

func NewMockService(deps ...interface{}) service.NoteService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.NoteRepository:
			srv.noteRepository = s
		}
	}

	return &srv
}
