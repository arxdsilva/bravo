package http

import (
	"context"
	"fmt"

	"github.com/arxdsilva/bravo/internal/jwt"
	"github.com/arxdsilva/bravo/internal/service"
	"github.com/labstack/echo/v4"
)

type Server struct {
	server        *echo.Echo
	service       service.Orders
	tokenProvider *jwt.TokenProvider
	config        Config
}

func NewServer(svc service.Orders, tp *jwt.TokenProvider, cfg Config) Server {
	return Server{
		service:       svc,
		tokenProvider: tp,
		config:        cfg,
	}
}

func (s Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.stop()
	}()
	s.server = echo.New()
	RegisterMiddlewares(s.server)
	RouterRegister(s.server)
	return s.server.Start(fmt.Sprintf(":%v", s.config.Port))
}

func (s Server) stop() {
	s.server.Shutdown(context.Background())
}
