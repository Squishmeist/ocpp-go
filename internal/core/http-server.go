package core

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type HttpServer struct {
	ServiceName string
	e           *echo.Echo
	routeLock   sync.Mutex
}

type Option func(*HttpServer)

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

	log.Infof("server listening on %s", port)
	if err := s.e.Start(port); err != nil && err != http.ErrServerClosed {
		log.Errorf("failed to serve: %v", err)
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