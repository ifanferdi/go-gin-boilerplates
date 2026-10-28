package grpcinfra

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Engine *grpc.Server
}

func NewServer(opts ...grpc.ServerOption) *Server {
	srv := grpc.NewServer(opts...)
	reflection.Register(srv)
	return &Server{Engine: srv}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC on port %s: %w", port, err)
	}

	log.Printf("🚀 gRPC Server listening on :%s\n", port)
	return s.Engine.Serve(lis)
}

func (s *Server) Stop() {
	if s.Engine != nil {
		s.Engine.GracefulStop()
	}
}
