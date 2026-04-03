package worker

import (
	"errors"
	"fmt"
	"net"

	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
	"github.com/nlduy0310/simplelog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	logger              simplelog.Logger
	id                  string
	cfg                 config
	clientConn          *grpc.ClientConn
	masterServiceClient rpcv1.MasterServiceClient
	listener            net.Listener
	grpcServer          *grpc.Server
	workerService       *WorkerService
	registered          bool // will need rework
}

func Setup() (*Server, error) {
	cfg, err := autoConfig()
	if err != nil {
		return nil, errorsx.Wrap("can not initialize config", err)
	}

	// client setup
	grpcClient, err := grpc.NewClient(cfg.MasterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errorsx.Wrap(fmt.Sprintf("can not create client to %s", cfg.MasterAddr), err)
	}
	masterClient := rpcv1.NewMasterServiceClient(grpcClient)

	// server setup
	listener, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		return nil, errorsx.Wrap(fmt.Sprintf("can not listen on %s", cfg.ListenAddr), err)
	}
	grpcServer := grpc.NewServer()
	workerService := NewWorkerService()
	rpcv1.RegisterWorkerServiceServer(grpcServer, workerService)

	id := genId()
	return &Server{
		logger:              simplelog.NewLogger(id, simplelog.DEBUG),
		id:                  id,
		cfg:                 cfg,
		clientConn:          grpcClient,
		masterServiceClient: masterClient,
		listener:            listener,
		grpcServer:          grpcServer,
		workerService:       workerService,
	}, nil
}

func (s *Server) Serve() error {
	s.registerIfNeeded()
	s.logger.Info("listening on %s", s.cfg.ListenAddr)
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Close() error {
	return errors.Join(
		s.clientConn.Close(),
		s.listener.Close(),
	)
}

func (s *Server) registerIfNeeded() {
	if !s.registered {
		s.initialRegister()
	}
}
