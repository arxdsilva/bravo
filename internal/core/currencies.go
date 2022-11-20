package core

type Currencies []Currency

type Currency struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Source      string `json:"source"`
}
