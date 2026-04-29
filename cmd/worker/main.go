package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/logx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/server"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/worker"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	port            = flag.Int("port", 5000, "the port to listen on")
	masterAddr      = flag.String("master-address", "", "master address")
	advertiseAddr   = flag.String("advertise-address", "", "advertise address")
	registerTimeout = flag.Int("register-timeout", 5, "register timeout in seconds")
)

func run() int {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	svrConfig, err := server.NewConfig(*port, *advertiseAddr)
	if err != nil {
		logx.Err("server config", err)
		return 1
	}

	svr, err := server.New(svrConfig)
	if err != nil {
		logx.Err("configure server", err)
		return 1
	}
	defer svr.Close()

	svc, err := worker.NewService(worker.Config{
		MasterAddr:      *masterAddr,
		AdvertiseAddr:   svr.Config.AdvertiseAddr,
		RegisterTimeout: time.Duration(*registerTimeout) * time.Second,
	})
	if err != nil {
		logx.Err("configure service", err)
		return 1
	}
	defer svc.Close()

	if err = svc.Init(); err != nil {
		logx.Err("initialize service", err)
		return 1
	}

	rpcv1.RegisterWorkerServiceServer(svr.GrpcServer, svc)

	if err := svr.Serve(ctx); err != nil {
		logx.Err("exited with error", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
