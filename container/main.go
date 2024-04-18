package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/94peter/log"
	"github.com/94peter/microservice"
	"github.com/94peter/microservice/di"
	"github.com/94peter/microservice/grpc_tool"
	"github.com/94peter/serial_number/grpc/pb"
	"github.com/94peter/serial_number/grpc/service"
	"github.com/94peter/serial_number/model"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	v         = flag.Bool("v", false, "version")
	Version   = "1.0.0"
	BuildTime = time.Now().Local().GoString()
)

func main() {
	flag.Parse()

	if *v {
		fmt.Println("Version: " + Version)
		fmt.Println("Build Time: " + BuildTime)
		return
	}

	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envFile := exePath + "/.env"
	if model.FileExists(envFile) {
		err := godotenv.Load(envFile)
		if err != nil {
			panic(errors.Wrap(err, "load .env file fail"))
		}
	}

	modelCfg, err := model.GetModelCfgFromEnv()
	if err != nil {
		panic(err)
	}

	ms, err := microservice.New(modelCfg, &mydi{})
	if err != nil {
		panic(err)
	}

	serv := newService(ms)
	microservice.RunService(
		serv.StartGrpc,
	)

	fmt.Println("Bye")
}

type mydi struct {
	di.CommonServiceDI
	*log.LoggerConf `yaml:"log"`
}

func (d *mydi) IsConfEmpty() error {
	if os.Getenv("LOG_TARGET") == "fluent" &&
		d.LoggerConf.FluentLog == nil {
		return errors.New("log.FluentLog no set")
	}
	return nil
}

type mService struct {
	microservice.MicroService[*model.Config, *mydi]
}

func newService(ms microservice.MicroService[*model.Config, *mydi]) *mService {
	return &mService{MicroService: ms}
}

func (s *mService) StartGrpc(ctx context.Context) {
	cfg, err := s.NewCfg("grpc")
	if err != nil {
		panic(err)
	}
	defer cfg.Close()
	grpcCfg, err := grpc_tool.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}
	grpcCfg.Logger = cfg.Log
	grpcCfg.SetRegisterServiceFunc(func(grpcServ *grpc.Server) {
		pb.RegisterGcpServiceServer(grpcServ, service.NewGcp(cfg))
	})

	err = grpc_tool.RunGrpcServ(ctx, grpcCfg)
	if err != nil {
		panic(fmt.Sprintf("grpc run fail: %v", err))
	}
}
