// a wrapper for the grpc.Server type
package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Config     Config
	listener   net.Listener
	GrpcServer *grpc.Server
}

func New(config Config) (*Server, error) {
	listener, err := net.Listen("tcp", listenAddr(config.Port))
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	return &Server{config, listener, grpcServer}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.GrpcServer.GracefulStop()
	}()

	log.Printf("server listening on port %d (advertises %s)", s.Config.Port, s.Config.AdvertiseAddr)
	return s.GrpcServer.Serve(s.listener)
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func listenAddr(port int) string {
	return fmt.Sprintf("0.0.0.0:%d", port)
}
