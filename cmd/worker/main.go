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
	initTimeout       = flag.Int("init-timeout", 30, "init timeout in seconds")
	heartbeatInterval = flag.Int("heartbeat-interval", 5, "heartbeat interval in seconds")
	heartbeatTimeout  = flag.Int("heartbeat-timeout", 3, "heartbeat timeout in seconds")
)

var (
	errHeartbeatFailure = errors.New("heartbeat failure")
)

func run() int {
	flag.Parse()

	if err := validateFlags(); err != nil {
		logx.Err("configuring application", err)
		return 1
	}

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
		Name:          *name,
		MasterAddr:    *masterAddr,
		AdvertiseAddr: svr.Config.AdvertiseAddr,
	})
	if err != nil {
		logx.Err("configure service", err)
		return 1
	}
	defer svc.Close()

	initCtx, timeoutInit := context.WithTimeout(ctx, time.Duration(*initTimeout)*time.Second)
	defer timeoutInit()
	if err = svc.Init(initCtx); err != nil {
		logx.Err("initialize service", err)
		return 1
	}

	rpcv1.RegisterWorkerServiceServer(svr.GrpcServer, svc)

	go func() {
		err := periodicHeartbeat(ctx, svc, time.Duration(*heartbeatInterval)*time.Second, time.Duration(*heartbeatTimeout)*time.Second)
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

func periodicHeartbeat(ctx context.Context, svc *worker.Service, interval time.Duration, timeout time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := heartbeatWithTimeout(ctx, svc, timeout); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func heartbeatWithTimeout(parent context.Context, svc *worker.Service, timeout time.Duration) error {
	ctx, timeoutHeartbeat := context.WithTimeout(parent, timeout)
	defer timeoutHeartbeat()
	if err := svc.DoHeartbeat(ctx); err != nil {
		return err
	}

	return nil
}

func validateFlags() error {
	return nil
}

func main() {
	os.Exit(run())
}
