package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/master"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/server"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	port = flag.Int("port", 8000, "the port that master will listen on")
)

func run() int {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := server.NewConfig(*port)
	if err != nil {
		log.Printf("initialize config: %s", err)
		return 1
	}

	svr, err := server.New(cfg)
	if err != nil {
		log.Printf("initialize server: %s", err)
		return 1
	}
	defer svr.Close()

	svc, err := master.NewService()
	if err != nil {
		log.Printf("initialize service: %s", err)
		return 1
	}

	rpcv1.RegisterMasterServiceServer(svr.GrpcServer, svc)

	if err := svr.Serve(ctx); err != nil {
		log.Printf("server stopped with error: %s", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
