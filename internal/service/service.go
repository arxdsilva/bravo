package service

import (
	"context"

	"github.com/arxdsilva/bravo/internal/core"
)

type Resolver interface {
	Convert(ctx context.Context, conv core.ConversionSVC) (amount float64, source string, err error)
}

type Exchanger interface {
	GetCurrencies(ctx context.Context) (map[string]string, error)
	Exchange(ctx context.Context, from, to string, amount float64) (core.ConversionResp, error)
}

type Service struct {
	Repo     ConversionRepository
	Exchange Exchanger
}

func NewService(repo ConversionRepository, exchange Exchanger) Service {
	return Service{
		Repo:     repo,
		Exchange: exchange,
	}
}

func (s Service) Convert(ctx context.Context, conv core.ConversionSVC) (amount float64, source string, err error) {
	resp, err := s.Exchange.Exchange(ctx, conv.From, conv.To, conv.Amount)
	if err != nil {
		return
	}
	// todo: store this search
	return resp.ConvertedAmount, resp.ConversionSource, err
}
