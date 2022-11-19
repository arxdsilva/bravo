package core

import "errors"

var (
	ErrInvalidFromCurrency = errors.New("invalid From currency")
	ErrInvalidToCurrency   = errors.New("invalid To currency")
	ErrAmountIsNotANumber  = errors.New("amount is not a number")
)
