package note

import (
	"di_container/internal/client/rpc"
	"di_container/internal/service"
	desc "di_container/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedNoteV1Server
	noteService        service.NoteService
	otherServiceClient rpc.OtherServiceClient
}

func NewImplementation(noteService service.NoteService, otherServiceClient rpc.OtherServiceClient) *Implementation {
	return &Implementation{
		noteService:        noteService,
		otherServiceClient: otherServiceClient,
	}
}
