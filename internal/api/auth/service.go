package auth

import (
	"di_container/internal/config/env"
	desc "di_container/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	config *env.TokenConfigData
}

func NewImplementation(config *env.TokenConfigData) *Implementation {
	return &Implementation{config: config}
}
