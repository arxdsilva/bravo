package core

type SymbolsClientResp struct {
	Success bool              `json:"success"`
	Symbols map[string]string `json:"symbols"`
}

type ConvertClientResp struct {
	Date       string `json:"date"`
	Historical string `json:"historical"`
	Info       struct {
		Rate      float64 `json:"rate"`
		Timestamp int     `json:"timestamp"`
	} `json:"info"`
	Query struct {
		Amount int    `json:"amount"`
		From   string `json:"from"`
		To     string `json:"to"`
	} `json:"query"`
	Result  float64 `json:"result"`
	Success bool    `json:"success"`
}
