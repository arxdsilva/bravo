package core

import "errors"

var (
	ErrInvalidFromCurrency = errors.New("invalid From currency")
	ErrInvalidToCurrency   = errors.New("invalid To currency")
	ErrAmountIsNotANumber  = errors.New("amount is not a number")
	// currency errors
	ErrEmptySymbol      = errors.New("currency needs a symbol")
	ErrSymbolMinLen     = errors.New("currency symbol has to have 3 or more characters")
	ErrRateIsZero       = errors.New("currency convertion rate cannot be zero")
	ErrCurrencyNotFound = errors.New("currency not found")
	// general
	ErrNotFound = errors.New("not found")
)
