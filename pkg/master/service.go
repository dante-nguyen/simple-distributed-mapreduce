package master

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/logx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	errInvalidConfig = errors.New("invalid config")
)

type Service struct {
	rpcv1.UnimplementedMasterServiceServer
	Config Config
	reg    *registry
}

func (s *Service) Register(_ context.Context, req *rpcv1.RegisterRequest) (*rpcv1.RegisterResponse, error) {
	if err := s.reg.register(req.Name, req.Address); err != nil {
		logx.Err(fmt.Sprintf("register worker %q at %q", req.Name, req.Address), err)
		return nil, err
	}

	log.Printf("successfully registered worker %q at %q", req.Name, req.Address)
	return &rpcv1.RegisterResponse{Ok: true}, nil
}

func (s *Service) Heartbeat(_ context.Context, req *rpcv1.HeartbeatRequest) (*rpcv1.HeartbeatResponse, error) {
	ts := time.Now()
	if err := s.reg.recordHeartbeat(req.Name, ts); err != nil {
		return nil, err
	}

	// log.Printf("received heartbeat from %q at %s", req.Name, ts)
	return &rpcv1.HeartbeatResponse{Ok: true}, nil
}

func NewService(cfg Config) (*Service, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, errx.Chain(errInvalidConfig, err)
	}

	return &Service{
		Config: cfg,
		reg:    newRegistry(),
	}, nil
}
