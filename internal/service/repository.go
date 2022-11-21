package service

import "context"

type Repository interface {
	CreateCurrency(ctx context.Context, symbol, description, source string) error
	CountCurrencies(ctx context.Context) (int, error)
}
