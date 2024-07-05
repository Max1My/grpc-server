package access

import (
	"context"
	"di_container/internal/model"
	"di_container/internal/utils"
	desc "di_container/pkg/access_v1"
	"errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

var accessibleRoles map[string]string

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("Metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("Authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], i.config.AuthPrefix) {
		return nil, errors.New("Invalid authorization")
	}

	accessToken := strings.TrimPrefix(authHeader[0], i.config.AuthPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(i.config.AccessTokenSecretKey))
	if err != nil {
		return nil, errors.New("Access token is invalid")
	}

	accessibleMap, err := i.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New("Failed to get accessible roles")
	}

	role, ok := accessibleMap[req.GetEndpointAddress()]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	if role == claims.Role {
		return &emptypb.Empty{}, nil
	}

	return nil, errors.New("Access denied")
}

func (i *Implementation) accessibleRoles(ctx context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		// Лезем в базу за данными о доступных ролях для каждого эндпоинта
		// Можно кешировать данные, чтобы не лезть в базу каждый раз

		// Например, для эндпоинта /note_v1.NoteV1/Get доступна только роль admin
		accessibleRoles[model.ExamplePath] = "admin"
	}

	return accessibleRoles, nil
}
