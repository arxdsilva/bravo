package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/arxdsilva/bravo/internal/core"
	log "github.com/sirupsen/logrus"
)

type Exchanger interface {
	GetCurrencies() (map[string]string, error)
	Exchange(to, from string, amount float64) (core.ConversionResp, error)
}

type Exchange struct {
	APIKey string
	client http.Client
}

func New(cfg Config) Exchanger {
	return Exchange{
		APIKey: cfg.APIKey,
		client: http.Client{
			Timeout: time.Duration(time.Second * 10),
		},
	}
}

func (e Exchange) GetCurrencies() (l map[string]string, err error) {
	lg := log.WithField("pkg", "exchange")
	req, err := http.NewRequest(http.MethodGet, "https://api.apilayer.com/exchangerates_data/symbols", nil)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] NewRequest")
		return
	}
	req.Header.Add("apiKey", e.APIKey)
	resp, err := e.client.Do(req)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] client.Do")
		return
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] ReadAll")
		return
	}
	symbols := &core.SymbolsClientResp{
		Symbols: map[string]string{},
	}
	err = json.Unmarshal(b, symbols)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] Unmarshal")
		return
	}
	if !symbols.Success {
		lg.WithField("success", symbols.Success).Warn("[GetCurrencies] success")
		return
	}
	lg.Info("[GetCurrencies] ok")
	return symbols.Symbols, err
}

func (e Exchange) Exchange(to, from string, amount float64) (resp core.ConversionResp, err error) {
	lg := log.WithField("pkg", "exchange")
	lg.Info("[Exchange] ok")
	return
}
