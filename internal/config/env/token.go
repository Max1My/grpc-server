package env

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

const (
	refreshTokenSecretKeyEnvName  = "REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretKeyEnvName   = "ACCESS_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION"
	accessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"
	authPrefixEnvName             = "AUTH_PREFIX"
)

type TokenConfig interface {
}

type TokenConfigData struct {
	RefreshTokenSecretKey  string
	AccessTokenSecretKey   string
	RefreshTokenExpiration time.Duration
	AccessTokenExpiration  time.Duration
	AuthPrefix             string
}

func NewTokenConfig() (*TokenConfigData, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refresh token secret key not found")
	}
	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("access token secret key not found")
	}
	refreshTokenExpirationStr := os.Getenv(refreshTokenExpirationEnvName)
	if len(refreshTokenExpirationStr) == 0 {
		return nil, errors.New("refresh token expiration not found")
	}
	refreshTokenExpiration, err := strconv.ParseInt(refreshTokenExpirationStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid refresh token expiration value")
	}
	accessTokenExpirationStr := os.Getenv(accessTokenExpirationEnvName)
	if len(accessTokenExpirationStr) == 0 {
		return nil, errors.New("access token expiration not found")
	}
	accessTokenExpiration, err := strconv.ParseInt(accessTokenExpirationStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid access token expiration value")
	}
	authPrefix := os.Getenv(authPrefixEnvName)
	if len(authPrefix) == 0 {
		return nil, errors.New("auth prefix not found")
	}

	return &TokenConfigData{
		RefreshTokenSecretKey:  refreshTokenSecretKey,
		AccessTokenSecretKey:   accessTokenSecretKey,
		RefreshTokenExpiration: time.Duration(refreshTokenExpiration) * time.Minute,
		AccessTokenExpiration:  time.Duration(accessTokenExpiration) * time.Minute,
		AuthPrefix:             authPrefix,
	}, nil
}
