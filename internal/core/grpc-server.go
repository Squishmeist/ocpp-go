package core

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	Grpc        *grpc.Server
	serviceName string
	listener    net.Listener
	port        string
}

type GrpcOption func(*GrpcServer)

func (g *GrpcServer) Validate() error {
	if g.serviceName == "" {
		return fmt.Errorf("service name is not set")
	}
	if g.port == "" {
		return fmt.Errorf("port is not set")
	}
	return nil
}

func WithGrpcServiceName(serviceName string) GrpcOption {
	return func(s *GrpcServer) {
		s.serviceName = serviceName
	}
}

func WithGrpcPort(port string) GrpcOption {
	return func(s *GrpcServer) {
		s.port = port
	}
}

func NewGrpcServer(opts ...GrpcOption) *GrpcServer {
	server := &GrpcServer{}

	for _, opt := range opts {
		opt(server)
	}

	if err := server.Validate(); err != nil {
		slog.Error("failed to validate server", "error", err)
		panic(err)
	}

	lis, err := net.Listen("tcp", server.port)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		panic(err)
	}
	server.listener = lis

	server.Grpc = grpc.NewServer()

	reflection.Register(server.Grpc)
	return server
}

func (s *GrpcServer) Start() {
	slog.Info(fmt.Sprintf("Server listening on %s", s.listener.Addr().String()))
	if err := s.Grpc.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
		slog.Error("failed to serve", "error", err)
		panic(err)
	}
}

func (s *GrpcServer) Shutdown() {
	if s.Grpc != nil {
		s.Grpc.Stop()
	}
	if s.listener != nil {
		s.listener.Close()
	}
}
