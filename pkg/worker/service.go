package worker

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/client"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/task"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/validate"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

type Service struct {
	rpcv1.UnimplementedWorkerServiceServer
	Config Config
	Name   string
	client *client.Client
	master rpcv1.MasterServiceClient
	// states
	curTask task.Task
}

func NewService(cfg Config) (*Service, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, errx.WithContext(err, "invalid config")
	}

	client, err := client.New(cfg.MasterAddr)
	if err != nil {
		return nil, errx.WithContext(err, "init master client")
	}

	master := rpcv1.NewMasterServiceClient(client.Conn)

	name := cfg.Name
	if len(name) == 0 {
		name = randomName()
	}

	return &Service{
		Config:  cfg,
		Name:    name,
		client:  client,
		master:  master,
		curTask: task.New(),
	}, nil
}

func (s *Service) Ping(context.Context, *rpcv1.PingRequest) (*rpcv1.PingResponse, error) {
	return &rpcv1.PingResponse{Message: "pong"}, nil
}

func (s *Service) Map(_ context.Context, req *rpcv1.MapRequest) (*rpcv1.MapResponse, error) {
	absPath := filepath.Join(s.Config.NfsRoot, req.NfsPath)
	if err := validate.EnsureIsFile(absPath); err != nil {
		return &rpcv1.MapResponse{Ok: false, Reason: errx.WithContext(err, fmt.Sprintf("find file %s", absPath)).Error()}, nil
	}

	if ok := s.curTask.SetMap(req.NfsPath); !ok {
		return &rpcv1.MapResponse{Ok: false, Reason: "worker is busy"}, nil
	}

	s.startMap()

	return &rpcv1.MapResponse{Ok: true}, nil
}

func (s *Service) Init(ctx context.Context) error {
	if err := s.register(ctx); err != nil {
		return errx.WithContext(err, "register to master")
	}

	return nil
}

func (s *Service) Close() error {
	return s.client.Close()
}

func (s *Service) DoHeartbeat(ctx context.Context) error {
	if _, err := s.master.Heartbeat(ctx, &rpcv1.HeartbeatRequest{Name: s.Name}); err != nil {
		return err
	}

	return nil
}

func (s *Service) register(ctx context.Context) error {
	req := rpcv1.RegisterRequest{Name: s.Name, Address: s.Config.AdvertiseAddr}
	_, err := s.master.Register(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) startMap() {
	// TODO
	mt, _ := s.curTask.GetMap()
	println("starting map task on", mt.Path)
}

func randomName() string {
	return fmt.Sprintf("worker-%03d", rand.Intn(1000))
}
