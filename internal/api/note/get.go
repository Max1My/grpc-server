package note

import (
	"context"
	"di_container/internal/converter"
	"di_container/internal/logger"
	"di_container/internal/sys"
	"di_container/internal/sys/validate"
	desc "di_container/pkg/note_v1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"math/rand"
	"time"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	err := validate.Validate(
		ctx,
		validate.ValidateID(req.GetId()),
		validate.OtherValidateID(req.GetId()),
	)
	if err != nil {
		return nil, err
	}

	if req.GetId() > 100 {
		return nil, sys.NewCommonError("id must be less than 100", codes.ResourceExhausted)
	}
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	logger.Info("Getting note...", zap.Int64("id", req.GetId()))
	noteObj, err := i.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		Note: converter.ToNoteFromService(noteObj),
	}, nil
}

//func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
//	if req.GetId() == 0 {
//		return nil, errors.Errorf("id is empty")
//	}
//
//	// rand.Intn(max - min) + min
//	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
//
//	span, ctx := opentracing.StartSpanFromContext(ctx, "get note")
//	defer span.Finish()
//
//	span.SetTag("id", req.GetId())
//
//	note, err := i.otherServiceClient.Get(ctx, 0)
//	if err != nil {
//		return nil, errors.WithMessage(err, "getting note")
//	}
//
//	var updatedAt *timestamppb.Timestamp
//	if note.UpdatedAt.Valid {
//		updatedAt = timestamppb.New(note.UpdatedAt.Time)
//	}
//
//	return &desc.GetResponse{
//		Note: &desc.Note{
//			Id: note.ID,
//			Info: &desc.NoteInfo{
//				Title:   note.Info.Title,
//				Content: note.Info.Content,
//			},
//			CreatedAt: timestamppb.New(note.CreatedAt),
//			UpdatedAt: updatedAt,
//		},
//	}, nil
//}
