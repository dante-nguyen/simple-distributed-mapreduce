package master

import (
	"context"

	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	workerregistry "github.com/nlduy0310/simple-distributed-mapreduce/master/worker-registry"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

type MasterService struct {
	rpcv1.UnimplementedMasterServiceServer
	workerRegistry *workerregistry.Registry
}

func NewMasterService() *MasterService {
	return &MasterService{
		workerRegistry: workerregistry.NewRegistry(),
	}
}

func (s *MasterService) RegisterWorker(ctx context.Context, req *rpcv1.RegisterWorkerRequest) (*rpcv1.RegisterWorkerResponse, error) {
	if err := s.workerRegistry.Register(req.Id); err != nil {
		return nil, errorsx.Wrap("can not register worker", err)
	}

	logger.Debug("registered worker %s", req.Id)
	return &rpcv1.RegisterWorkerResponse{}, nil
}
