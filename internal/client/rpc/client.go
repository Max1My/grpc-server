package rpc

import (
	"context"

	"di_container/internal/model"
)

type OtherServiceClient interface {
	Get(ctx context.Context, id int64) (*model.Note, error)
}
