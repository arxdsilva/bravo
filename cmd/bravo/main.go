package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arxdsilva/bravo/internal/clients/exchange"
	"github.com/arxdsilva/bravo/internal/http"
	"github.com/arxdsilva/bravo/internal/logger"
	"github.com/arxdsilva/bravo/internal/option"
	"github.com/arxdsilva/bravo/internal/service"
	"github.com/arxdsilva/bravo/internal/storage/postgres"
	migration "github.com/gosidekick/migration/v3"
	"golang.org/x/sync/errgroup"

	_ "github.com/lib/pq"
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
	lg := log.WithFields(log.Fields{
		"pkg": "main",
		"fn":  "startup",
	})

	lg.Info("starting the service")

	// maybe add some tracer

	excg := exchange.New(cfg.Exchange)

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf(`could not connect to db %w`, err)
	}

	err = migrate(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf(`could not migrate db %w`, err)
	}

	svc := service.NewService(db, excg)

	go svc.Seed(ctx)

	srv := http.NewServer(svc, cfg.HTTP)

	// run seed svc

	errg, ctx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return srv.Run(ctx)
	})

	lg.Info("service started")

	if err := errg.Wait(); err != nil {
		lg.Error(err)
	}

	lg.Info("service stopped")
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

func migrate(ctx context.Context, cfg postgres.Config) error {
	n, executed, err := migration.Run(
		ctx, "migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User, cfg.Password,
			cfg.Host, cfg.Port,
			cfg.DB, "disable"), "up")
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"number": n, "executed": executed}).Info("migration")
	return nil
}
