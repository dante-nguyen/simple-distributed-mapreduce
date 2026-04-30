package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/logx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/server"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/worker"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	name              = flag.String("name", "", "the worker's identity")
	port              = flag.Int("port", 5000, "the port to listen on")
	masterAddr        = flag.String("master-address", "", "master address")
	advertiseAddr     = flag.String("advertise-address", "", "advertise address")
	registerTimeout   = flag.Int("register-timeout", 5, "register timeout in seconds")
	heartbeatInterval = flag.Int("heartbeat-interval", 5, "heartbeat interval in seconds")
)

var (
	errHeartbeatFailure = errors.New("heartbeat failure")
)

func run() int {
	flag.Parse()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancelWithCause := context.WithCancelCause(signalCtx)
	defer cancelWithCause(nil)

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
		Name:            *name,
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

	go func() {
		err := periodicHeartbeat(ctx, svc, time.Duration(*heartbeatInterval)*time.Second)
		if err != nil {
			cancelWithCause(errx.Chain(errHeartbeatFailure, err))
		}
	}()

	if err := svr.Serve(ctx); err != nil {
		logx.Err("exited with error", err)
		if ctxErr, cause := ctx.Err(), context.Cause(ctx); ctxErr != nil && ctxErr != cause {
			logx.Err("cause", cause)
		}
		return 1
	}

	return 0
}

func periodicHeartbeat(ctx context.Context, svc *worker.Service, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// TODO make this timeout
			if err := svc.DoHeartbeat(); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func main() {
	os.Exit(run())
}
