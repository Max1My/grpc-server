package access

import (
	"di_container/internal/config/env"
	desc "di_container/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessV1Server
	config *env.TokenConfigData
}

func NewImplementation(config *env.TokenConfigData) *Implementation {
	return &Implementation{config: config}
}
