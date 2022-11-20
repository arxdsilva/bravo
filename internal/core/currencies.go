package core

type Currencies []Currency

type Currency struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

func (c Currency) Check() error {
	if c.Symbol == "" {
		return ErrEmptySymbol
	}
	if len(c.Symbol) < 3 {
		return ErrSymbolMinLen
	}
	return nil
}

type CurrencyRate struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

func (c CurrencyRate) Check() (err error) {
	if len(c.From) < 3 {
		return ErrSymbolMinLen
	}
	if len(c.To) < 3 {
		return ErrSymbolMinLen
	}
	if c.Rate == 0 {
		return ErrRateIsZero
	}
	return err
}
