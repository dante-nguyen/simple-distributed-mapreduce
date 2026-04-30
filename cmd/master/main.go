package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/logx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/master"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/server"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	port                = flag.Int("port", 8000, "the port that master will listen on")
	advertiseAddr       = flag.String("advertise-address", "", "advertise address")
	healthyDuration     = flag.Int("healthy-duration", 30, "maximum duration in seconds since last heartbeat of a healthy worker")
	healthcheckInterval = flag.Int("healthcheck-interval", 5, "interval in seconds between worker heartbeat checks")
	healthcheckTimeout  = flag.Int("healthcheck-timeout", 3, "timeout for each worker heartbeat check")
)

func run() int {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	svrConfig, err := server.NewConfig(*port, *advertiseAddr)
	if err != nil {
		logx.Err("initialize config", err)
		return 1
	}

	svr, err := server.New(svrConfig)
	if err != nil {
		logx.Err("configure server", err)
		return 1
	}
	defer svr.Close()

	svcConfig := master.Config{}
	svc, err := master.NewService(svcConfig)
	if err != nil {
		logx.Err("configure service", err)
		return 1
	}

	rpcv1.RegisterMasterServiceServer(svr.GrpcServer, svc)

	go func() {
		svc.PeriodicHealthcheck(
			ctx,
			time.Duration(*healthcheckInterval)*time.Second,
			time.Duration(*healthcheckTimeout)*time.Second,
			time.Duration(*healthyDuration)*time.Second,
		)
	}()

	if err := svr.Serve(ctx); err != nil {
		logx.Err("exited with error", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
