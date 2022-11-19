package service

import (
	"context"

	"github.com/arxdsilva/bravo/internal/core"
)

type Conversion interface {
	Convert(ctx context.Context, conv core.ConversionSVC) (amount float64, source string, err error)
}

type Service struct {
	Repo ConversionRepository
}

func NewService(repo ConversionRepository) Service {
	return Service{repo}
}

func (s Service) Convert(ctx context.Context, conv core.ConversionSVC) (amount float64, source string, err error) {

	return
}
