package master

import (
	"context"

	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	inputregistry "github.com/nlduy0310/simple-distributed-mapreduce/master/input-registry"
	workerregistry "github.com/nlduy0310/simple-distributed-mapreduce/master/worker-registry"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

type MasterService struct {
	rpcv1.UnimplementedMasterServiceServer
	workerRegistry *workerregistry.Registry
	inputRegistry  *inputregistry.Registry
}

func NewMasterService(filePaths []string) (*MasterService, error) {
	inputRegistry, err := inputregistry.FromPaths(filePaths)
	if err != nil {
		return nil, errorsx.Wrap("can not initialize input registry", err)
	}
	logger.Debugf("initialized input registry of %d files", inputRegistry.Size())

	return &MasterService{
		workerRegistry: workerregistry.EmptyRegistry(),
		inputRegistry:  inputRegistry,
	}, nil
}

func (s *MasterService) RegisterWorker(ctx context.Context, req *rpcv1.RegisterWorkerRequest) (*rpcv1.RegisterWorkerResponse, error) {
	if err := s.workerRegistry.Register(req.Id); err != nil {
		return nil, errorsx.Wrap("can not register worker", err)
	}

	logger.Debugf("registered worker %s", req.Id)
	return &rpcv1.RegisterWorkerResponse{}, nil
}
