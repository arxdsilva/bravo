package http

import (
	"context"
	"fmt"

	"github.com/arxdsilva/bravo/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	server  *echo.Echo
	service service.Resolver
	config  Config
}

func NewServer(svc service.Resolver, cfg Config) Server {
	return Server{
		service: svc,
		config:  cfg,
	}
}

func (s Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.stop()
	}()
	s.server = echo.New()
	RegisterMiddlewares(s.server)
	s.RouterRegister(s.server)
	return s.server.Start(fmt.Sprintf(":%v", s.config.Port))
}

func (s Server) stop() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		log.Error("shutdown err: ", err)
	}
}
