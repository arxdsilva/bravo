package core

type ConvertClientResp struct {
	Success bool `json:"success"`
	Query   struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	} `json:"query"`
	Info struct {
		Rate float64 `json:"rate"`
	} `json:"info"`
	Historical bool    `json:"historical"`
	Date       string  `json:"date"`
	Result     float64 `json:"result"`
}

type SymbolsClientResp struct {
	Success bool              `json:"success"`
	Symbols map[string]Symbol `json:"symbols"`
}

type Symbol struct {
	Description string `json:"description"`
	Code        string `json:"code"`
}
