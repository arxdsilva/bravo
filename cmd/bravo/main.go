package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arxdsilva/bravo/internal/http"
	"github.com/arxdsilva/bravo/internal/jwt"
	"github.com/arxdsilva/bravo/internal/logger"
	"github.com/arxdsilva/bravo/internal/option"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupGracefulShutdown(cancel)

	cfg, err := option.FromEnv()
	if err != nil {
		return err
	}
	return startup(ctx, cfg)
}

func startup(ctx context.Context, cfg *option.Config) error {

	if err := logger.Setup(cfg.Log); err != nil {
		return fmt.Errorf(`log setup failed %w`, err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf(`config validation error %w`, err)
	}

	log.Info("starting the service")

	// maybe add some tracer

	provider := jwt.NewProvider(cfg.JWT) // take the provider and pass to the rest api

	// change this
	srv := http.NewServer(nil, &provider, cfg.HTTP)

	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return srv.Run(ctx)
	})

	log.Info("service started")

	if err := errg.Wait(); err != nil {
		log.Error(err)
	}

	log.Info("service stopped")
	return nil
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Info("received interrupt signal. Gracefully shutting down the service.")
		stop()
	}()
}
