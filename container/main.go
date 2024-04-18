package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"serial_number/grpc/pb"
	"serial_number/grpc/service"
	"serial_number/model"
	"sync"
	"syscall"
	"time"

	"github.com/94peter/di"
	"github.com/94peter/grpc-tool/autorun"
	"github.com/94peter/log"
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

	mainDI := &mydi{}
	diCfg, err := di.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}
	err = di.InitServiceDIByCfg(diCfg, mainDI)
	if err != nil {
		panic(err)
	}
	if err = mainDI.IsConfEmpty(); err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	log, err := mainDI.NewLogger(mainDI.GetService(), "main")
	if err != nil {
		panic(err)
	}
	modelCfg, err := model.GetModelCfgFromEnv()
	if err != nil {
		panic(err)
	}
	modelCfg.Log = log
	defer modelCfg.Close()

	grpcCfg, err := autorun.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}
	grpcCfg.Logger = log

	grpcCfg.SetRegisterServiceFunc(func(grpcServer *grpc.Server) {
		pb.RegisterGcpServiceServer(grpcServer, service.NewGcp(modelCfg))
	})

	var wg sync.WaitGroup
	wg.Add(1)

	grpcCtx, grpcCancel := context.WithCancel(context.Background())
	go func(ctx context.Context, grpcfg *autorun.GrpcConfig) {
		defer wg.Done()
		autorun.RunGrpcServ(ctx, grpcfg)
	}(grpcCtx, grpcCfg)

	<-sig
	if grpcCancel != nil {
		grpcCancel()
	}
	wg.Wait()
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
