package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const masterAddr = "127.0.0.1:5000"

func randomName() string {
	return fmt.Sprintf("worker-%d", rand.Int())
}

func main() {
	grpcClient, err := grpc.NewClient(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client to %s: %s\n", masterAddr, err)
	}
	defer grpcClient.Close()
	masterClient := rpcv1.NewMasterServiceClient(grpcClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := masterClient.Greet(ctx, &rpcv1.GreetRequest{Name: randomName()})
	if err != nil {
		log.Printf("request failed: %s\n", err)
	} else {
		log.Printf("received response from master: %s\n", resp.Message)
	}
}
