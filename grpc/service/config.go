package service

import (
	"serial_number/model"

	"google.golang.org/grpc"
)

type GrpcConfig struct {
	Port           int  `env:"GRPC_PORT"`
	ReflectService bool `env:"GRPC_REFLECT"`

	Logger              Log
	registerServiceFunc func(grpcServer *grpc.Server)
}

func (c *GrpcConfig) SetRegisterServiceFunc(f func(grpcServer *grpc.Server)) {
	c.registerServiceFunc = f
}

func GetConfigFromEnv() (*GrpcConfig, error) {
	var cfg GrpcConfig
	err := model.GetFromEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Log interface {
	Infof(format string, a ...any)
	Fatalf(format string, a ...any)
}
