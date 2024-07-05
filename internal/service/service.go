package service

import (
	"context"
	"di_container/internal/model"
)

type NoteService interface {
	Create(context.Context, *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}

type OtherService interface {
	Get(ctx context.Context, id int64) (*model.Note, error)
}

type AuthService interface {
	Login(ctx context.Context)
	GetAccessToken(ctx context.Context)
	GetRefreshToken(ctx context.Context)
	Check(ctx context.Context)
}
