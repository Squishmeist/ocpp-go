package core

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	ServiceName string
	e           *echo.Echo
	routeLock   sync.Mutex
}

type Option func(*HttpServer)

func WithServiceName(serviceName string) Option {
	return func(s *HttpServer) {
		s.ServiceName = serviceName
	}
}

func NewHttpServer(opts ...Option) *HttpServer {
	server := &HttpServer{
		e: echo.New(),
		routeLock: sync.Mutex{},
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *HttpServer) Start(port string) {
	s.AddRoute(http.MethodGet, "/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Hello from %s", s.ServiceName))
	})

	fmt.Printf("server listening on %s\n", port)
	if err := s.e.Start(port); err != nil && err != http.ErrServerClosed {
		fmt.Printf("failed to serve: %v\n", err)
		panic(err)
	}
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	if err := s.e.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	return nil
}

func (s *HttpServer) AddRoute(method, path string, handler echo.HandlerFunc) {
	s.routeLock.Lock()
	defer s.routeLock.Unlock()
	s.e.Add(method, path, handler)
}