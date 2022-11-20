package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/arxdsilva/bravo/internal/core"
	log "github.com/sirupsen/logrus"
)

type Exchanger interface {
	GetCurrencies() (map[string]string, error)
	Exchange(from, to string, amount float64) (core.ConversionResp, error)
}

type Exchange struct {
	APIKey  string
	BaseURL string
	client  http.Client
}

func New(cfg Config) Exchanger {
	return Exchange{
		APIKey:  cfg.APIKey,
		BaseURL: cfg.APIBaseURL,
		client: http.Client{
			Timeout: time.Duration(time.Second * 10),
		},
	}
}

func (e Exchange) GetCurrencies() (l map[string]string, err error) {
	lg := log.WithField("pkg", "exchange")
	url := fmt.Sprintf("%v/symbols", e.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
	symbols := &core.SymbolsClientResp{}
	err = json.Unmarshal(b, symbols)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] Unmarshal")
		return
	}
	if !symbols.Success {
		lg.WithField("success", symbols.Success).Warn("[GetCurrencies] success")
		return
	}
	l = map[string]string{}
	for k, v := range symbols.Symbols {
		l[k] = v.Description
	}
	lg.Info("[GetCurrencies] ok")
	return l, err
}

func (e Exchange) Exchange(from, to string, amount float64) (core.ConversionResp, error) {
	lg := log.WithField("pkg", "exchange")
	url := fmt.Sprintf("%v/convert", e.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		lg.WithError(err).Error("[Exchange] NewRequest")
		return core.ConversionResp{}, err
	}
	req.Header.Add("apiKey", e.APIKey)

	q := req.URL.Query()
	q.Add("from", from)
	q.Add("to", to)
	q.Add("amount", fmt.Sprintf("%f", amount))
	req.URL.RawQuery = q.Encode()

	resp, err := e.client.Do(req)
	if err != nil {
		lg.WithError(err).Error("[Exchange] client.Do")
		return core.ConversionResp{}, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		lg.WithError(err).Error("[Exchange] ReadAll")
		return core.ConversionResp{}, err
	}

	conv := &core.ConvertClientResp{}

	if err = json.Unmarshal(b, conv); err != nil {
		lg.WithError(err).Error("[Exchange] Unmarshal")
		return core.ConversionResp{}, err
	}

	if !conv.Success {
		lg.WithField("success", conv.Success).Warn("[Exchange] success")
		return core.ConversionResp{}, err
	}

	lg.Info("[Exchange] ok")
	return core.ConversionResp{
		From:             from,
		To:               to,
		OriginalAmount:   amount,
		ConvertedAmount:  amount * conv.Result,
		ConversionSource: "exchange",
	}, err
}
