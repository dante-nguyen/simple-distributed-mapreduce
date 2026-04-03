package master

import (
	"fmt"
	"net"

	"github.com/nlduy0310/simple-distributed-mapreduce/cli"
	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
	"github.com/nlduy0310/simplelog"
	"google.golang.org/grpc"
)

type Server struct {
	cfg           config
	masterService *MasterService
}

var logger = simplelog.NewLogger("master server", simplelog.DEBUG)

func Setup(opts cli.MasterCLIOptions) (*Server, error) {
	cfg, err := autoConfig()
	if err != nil {
		return nil, errorsx.Wrap("can not initialize config", err)
	}

	masterService, err := NewMasterService(opts.FilePaths)
	if err != nil {
		return nil, errorsx.Wrap("can not initialize master service", err)
	}

	return &Server{
		cfg:           cfg,
		masterService: masterService,
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

	logger.Infof("listening on %s", s.cfg.ListenAddress)
	return grpcSvr.Serve(listener)
}
