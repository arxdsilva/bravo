package core

import "strconv"

var allowedCurrencies = map[string]bool{
	"USD": true,
	"BRL": true,
	"EUR": true,
	"BTC": true,
	"ETH": true,
}

type ConversionAPI struct {
	From   string
	To     string
	Amount string
}

type ConversionSVC struct {
	From   string
	To     string
	Amount float64
}

func (c ConversionAPI) Check() (err error) {
	_, ok := allowedCurrencies[c.From]
	if !ok {
		return ErrInvalidFromCurrency
	}
	_, ok = allowedCurrencies[c.To]
	if !ok {
		return ErrInvalidToCurrency
	}
	_, err = strconv.ParseFloat(c.Amount, 64)
	if err != nil {
		return ErrAmountIsNotANumber
	}
	return err
}

// ConvertToService assumes that Check has already been called
// and everything is ok to proceed
//
// if currencies are equal, returns dont convert command and return same amount
func ConvertToService(c ConversionAPI) (cs ConversionSVC, should bool, err error) {
	amount, err := strconv.ParseFloat(c.Amount, 64)
	if err != nil {
		return cs, false, ErrAmountIsNotANumber
	}
	if c.From == c.To {
		return ConversionSVC{
			From:   c.From,
			To:     c.To,
			Amount: amount,
		}, false, nil
	}
	return ConversionSVC{
		From:   c.From,
		To:     c.To,
		Amount: amount,
	}, true, err
}
