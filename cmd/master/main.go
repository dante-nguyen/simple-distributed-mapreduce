package main

import (
	"context"
	"fmt"
	"log"
	"net"

	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
	"google.golang.org/grpc"
)

const addr = "127.0.0.1:5000"

type server struct {
	rpcv1.UnimplementedMasterServiceServer
}

func (s *server) Greet(ctx context.Context, req *rpcv1.GreetRequest) (*rpcv1.GreetResponse, error) {
	log.Printf("received greet from %s\n", req.GetName())
	return &rpcv1.GreetResponse{Message: fmt.Sprintf("hello, %s", req.GetName())}, nil
}

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can not listen on %s: %s\n", addr, err)
	}
	defer listener.Close()

	var svr server
	grpcServer := grpc.NewServer()
	rpcv1.RegisterMasterServiceServer(grpcServer, &svr)

	log.Printf("server listening at %s\n", addr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("server stopped with error: %s\n", err)
	}
}
