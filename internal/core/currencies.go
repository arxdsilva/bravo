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
