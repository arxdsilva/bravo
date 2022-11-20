package service

import (
	"context"

	"github.com/arxdsilva/bravo/internal/core"
)

type Resolver interface {
	Convert(ctx context.Context, conv core.ConversionSVC) (amount float64, source string, err error)
	GetCurrencies(ctx context.Context) (core.Currencies, error)
	AddCurrency(ctx context.Context, symbol, description string) error
	UpdateCurrency(ctx context.Context, symbol, description string) error
	GetCurrency(ctx context.Context, symbol string) (core.Currency, error)
	RemoveCurrency(ctx context.Context, symbol string) error
	GetRates(ctx context.Context) (interface{}, error)
	CreateRate(ctx context.Context, from, to string, rate float64) error
	UpdateRate(ctx context.Context, from, to string, rate float64) error
	RemoveRate(ctx context.Context, from, to string) error
}

type Exchanger interface {
	GetCurrencies(ctx context.Context) (map[string]string, error)
	Exchange(ctx context.Context, from, to string, amount float64) (core.ConversionResp, error)
}

type Service struct {
	Repo     Repository
	Exchange Exchanger
}

func NewService(repo Repository, exchange Exchanger) Service {
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

func (s Service) AddCurrency(ctx context.Context, symbol, description string) (err error) {
	// check whether the currency already exists
	// if exists return ok without inserting

	// add currency to repo
	return
}

func (s Service) UpdateCurrency(ctx context.Context, symbol, description string) (err error) {
	// check whether the currency exists
	// if not exists return error

	// update currency on repo
	return
}

func (s Service) GetCurrency(ctx context.Context, symbol string) (cr core.Currency, err error) {
	// could use a cache system to reduce DB toll

	// get currency from repo
	return
}

func (s Service) RemoveCurrency(ctx context.Context, symbol string) (err error) {
	// get currency from repo

	// remove from repo
	return
}

func (s Service) GetRates(ctx context.Context) (rts interface{}, err error) {
	// get currency rates from repo
	return
}

func (s Service) CreateRate(ctx context.Context, from, to string, rate float64) (err error) {
	// ensure currencies exist and are stored

	// create rate

	// create reverse rate
	return
}

func (s Service) UpdateRate(ctx context.Context, from, to string, rate float64) (err error) {
	// ensure currencies exist and are stored

	// update rate

	// update reverse rate
	return
}

func (s Service) RemoveRate(ctx context.Context, from, to string) (err error) {
	// ensure currencies exist and are stored

	// remove rate

	// remove reverse rate
	return
}
