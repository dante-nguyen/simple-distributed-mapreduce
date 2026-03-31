package master

import (
	"fmt"
	"net"

	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	"github.com/nlduy0310/simple-distributed-mapreduce/logging"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
	"google.golang.org/grpc"
)

type Server struct {
	cfg           config
	masterService *MasterService
}

var logger = logging.NewLogger("master server", logging.DEBUG)

func Setup() (*Server, error) {
	cfg, err := autoConfig()
	if err != nil {
		return nil, errorsx.Wrap("can not initialize config", err)
	}

	return &Server{
		cfg:           cfg,
		masterService: NewMasterService(),
	}, nil
}

func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.cfg.ListenAddress)
	if err != nil {
		return fmt.Errorf("can not listen on %s: %s", s.cfg.ListenAddress, err)
	}
	defer listener.Close()

	grpcSvr := grpc.NewServer()
	rpcv1.RegisterMasterServiceServer(grpcSvr, s.masterService)

	logger.Info("listening on %s", s.cfg.ListenAddress)
	return grpcSvr.Serve(listener)
}
