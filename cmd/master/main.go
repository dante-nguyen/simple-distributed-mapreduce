package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/fsx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/logx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/master"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/server"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	svrConfig, err := server.NewConfig(port, advertiseAddr)
	if err != nil {
		logx.Err(errx.WithContext(err, "initialize config"))
		return 1
	}

	svr, err := server.New(svrConfig)
	if err != nil {
		logx.Err(errx.WithContext(err, "configure server"))
		return 1
	}
	defer svr.Close()

	inputFiles, err := fsx.CollectPaths(inDir.Path, fsx.FilterFile)
	if err != nil {
		logx.Err(errx.WithContext(err, "list input files"))
		return 1
	}

	svcConfig := master.Config{InputFiles: inputFiles}
	svc, err := master.NewService(svcConfig)
	if err != nil {
		logx.Err(errx.WithContext(err, "configure service"))
		return 1
	}

	rpcv1.RegisterMasterServiceServer(svr.GrpcServer, svc)

	go func() {
		svc.PeriodicHealthcheck(ctx, healthcheckInterval, healthyDuration)
	}()

	if err := svr.Serve(ctx); err != nil {
		logx.Err(errx.WithContext(err, "exited with error"))
		return 1
	}

	return 0
}

func main() {
	prepArguments()
	os.Exit(run())
}
