package worker

import (
	"context"

	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

type WorkerService struct {
	rpcv1.UnimplementedWorkerServiceServer
}

func NewWorkerService() *WorkerService {
	return &WorkerService{}
}

func (s *WorkerService) Map(ctx context.Context, req *rpcv1.MapRequest) (*rpcv1.MapResponse, error) {
	return &rpcv1.MapResponse{}, nil
}
