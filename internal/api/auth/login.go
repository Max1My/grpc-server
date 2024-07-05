package auth

import (
	"context"
	"di_container/internal/model"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"di_container/internal/utils"
	desc "di_container/pkg/auth_v1"
)

func (i *Implementation) Login(_ context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	// Лезем в базу или кэш за данными пользователя
	// Сверяем хэши пароля

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: req.GetUsername(),
		// Это пример, в реальности роль должна браться из базы или кэша
		Role: "admin",
	},
		[]byte(i.config.RefreshTokenSecretKey),
		i.config.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("Failed to generate token")
	}

	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}

func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {

	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(i.config.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		// Это пример, в реальности роль должна браться из базы или кэша
		Role: "admin",
	},
		[]byte(i.config.RefreshTokenSecretKey),
		i.config.RefreshTokenExpiration,
	)

	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {

	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(i.config.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		// Это пример, в реальности роль должна браться из базы или кэша
		Role: "admin",
	},
		[]byte(i.config.AccessTokenSecretKey),
		i.config.AccessTokenExpiration,
	)

	if err != nil {
		return nil, err
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
