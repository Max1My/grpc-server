package env

import (
	"di_container/internal/config"
	"errors"
	"net"
	"os"
	"strconv"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName      = "GRPC_HOST"
	grpcPortEnvName      = "GRPC_PORT"
	grpcOtherPortEnvName = "GRPC_OTHER_PORT"
)

type grpcConfig struct {
	host      string
	port      string
	otherPort int64
}

func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	otherPort := os.Getenv(grpcOtherPortEnvName)
	if len(otherPort) == 0 {
		return nil, errors.New("grpc other port not found")
	}
	otherPortInt, err := strconv.ParseInt(otherPort, 10, 64)
	if err != nil {
		return nil, err
	}

	return &grpcConfig{
		host:      host,
		port:      port,
		otherPort: otherPortInt,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) OtherPort() int64 {
	return cfg.otherPort
}
