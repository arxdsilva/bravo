package core

type SymbolsClientResp struct {
	Success bool              `json:"success"`
	Symbols map[string]string `json:"symbols"`
}
