package service

import (
	"context"

	"github.com/arxdsilva/bravo/internal/core"
)

type Resolver interface {
	Convert(ctx context.Context, conv core.ConversionSVC) (amount float64, source string, err error)
	GetCurrencies(ctx context.Context) (core.Currencies, error)
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
	// always get latest rate and update on repo
	resp, err := s.Exchange.Exchange(ctx, conv.From, conv.To, conv.Amount)
	if err != nil {
		return
	}
	// todo: store this search
	return resp.ConvertedAmount, resp.ConversionSource, err
}

// todo: on init try to seed the repo
func (s Service) GetCurrencies(ctx context.Context) (cs core.Currencies, err error) {
	// get from repo

	// if none in repo, fall back to external
	currencies, err := s.Exchange.GetCurrencies(ctx)
	if err != nil {
		return
	}
	for k, v := range currencies {
		cs = append(cs, core.Currency{
			Symbol:      k,
			Description: v,
			Source:      "exchange",
		})
	}
	return
}
