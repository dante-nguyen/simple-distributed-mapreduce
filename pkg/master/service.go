package master

import (
	"context"
	"errors"
	"log"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	errInvalidConfig = errors.New("invalid config")
)

type Service struct {
	rpcv1.UnimplementedMasterServiceServer
	Config Config
}

func (s *Service) Register(_ context.Context, req *rpcv1.RegisterRequest) (*rpcv1.RegisterResponse, error) {
	log.Printf("registered worker %q at %q\n", req.Name, req.Address)
	return &rpcv1.RegisterResponse{Ok: true}, nil
}

func NewService(cfg Config) (*Service, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, errx.Chain(errInvalidConfig, err)
	}

	return &Service{Config: cfg}, nil
}
