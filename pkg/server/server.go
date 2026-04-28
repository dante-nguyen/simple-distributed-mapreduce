// a wrapper for the grpc.Server type
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

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

	log.Printf("server listening on port %d (advertise %s)", s.Config.Port, s.AdvertiseAddr())
	return s.GrpcServer.Serve(s.listener)
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) AdvertiseAddr() string {
	var hostname string
	var err error

	if len(s.Config.HostOverride) > 0 {
		hostname = s.Config.HostOverride
	} else {
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "127.0.0.1"
		}
	}

	return fmt.Sprintf("%s:%d", hostname, s.Config.Port)
}

func listenAddr(port int) string {
	return fmt.Sprintf("0.0.0.0:%d", port)
}
