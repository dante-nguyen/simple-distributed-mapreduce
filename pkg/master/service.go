package master

import (
	"context"
	"log"

	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

type Service struct {
	rpcv1.UnimplementedMasterServiceServer
}

func (s *Service) Register(_ context.Context, req *rpcv1.RegisterRequest) (*rpcv1.RegisterResponse, error) {
	log.Printf("registered worker %q at %q\n", req.Name, req.Address)
	return &rpcv1.RegisterResponse{Ok: true}, nil
}

func NewService() (*Service, error) {
	res := Service{}
	return &res, nil
}
