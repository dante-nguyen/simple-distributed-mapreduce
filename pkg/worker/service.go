package worker

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/client"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	errInvalidConfig    = errors.New("invalid config")
	errInitMasterClient = errors.New("initialize master client")
	errRegister         = errors.New("register to master")
)

type Service struct {
	rpcv1.UnimplementedWorkerServiceServer
	Config Config
	Name   string
	client *client.Client
	master rpcv1.MasterServiceClient
}

func (s *Service) Ping(context.Context, *rpcv1.PingRequest) (*rpcv1.PingResponse, error) {
	return &rpcv1.PingResponse{Message: "pong"}, nil
}

func NewService(cfg Config) (*Service, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, errx.Chain(errInvalidConfig, err)
	}

	client, err := client.New(cfg.MasterAddr)
	if err != nil {
		return nil, errx.Chain(errInitMasterClient, err)
	}

	master := rpcv1.NewMasterServiceClient(client.Conn)

	name := cfg.Name
	if len(name) == 0 {
		name = randomName()
	}

	return &Service{
		Config: cfg,
		Name:   name,
		client: client,
		master: master,
	}, nil
}

func (s *Service) Init() error {
	if err := s.register(); err != nil {
		return errx.Chain(errRegister, err)
	}

	return nil
}

func (s *Service) register() error {
	ctx, timeout := context.WithTimeout(context.Background(), s.Config.RegisterTimeout)
	defer timeout()

	req := rpcv1.RegisterRequest{Name: s.Name, Address: s.Config.AdvertiseAddr}
	_, err := s.master.Register(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Close() error {
	return s.client.Close()
}

func randomName() string {
	return fmt.Sprintf("worker-%03d", rand.Intn(1000))
}
