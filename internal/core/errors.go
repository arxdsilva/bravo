package core

import "errors"

var (
	ErrInvalidFromCurrency = errors.New("invalid From currency")
	ErrInvalidToCurrency   = errors.New("invalid To currency")
	ErrAmountIsNotANumber  = errors.New("amount is not a number")
	// currencies errors
	ErrEmptySymbol  = errors.New("currency needs a symbol")
	ErrSymbolMinLen = errors.New("currency symbol has to have 3 or more characters")
)
