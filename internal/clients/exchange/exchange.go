package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/arxdsilva/bravo/internal/core"
	log "github.com/sirupsen/logrus"
)

type Exchange struct {
	APIKey  string
	BaseURL string
	client  http.Client
}

func New(cfg Config) Exchange {
	return Exchange{
		APIKey:  cfg.APIKey,
		BaseURL: cfg.APIBaseURL,
		client: http.Client{
			Timeout: time.Duration(time.Second * 10),
		},
	}
}

// GetCurrencies tries to get all currencies from our currency provider
//
// receives a ctx so the request can be cancelled if the original request is also cancelled
func (e Exchange) GetCurrencies(ctx context.Context) (l map[string]string, err error) {
	lg := log.WithField("pkg", "exchange")
	url := fmt.Sprintf("%v/symbols", e.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] NewRequest")
		return
	}

	req = req.WithContext(ctx)
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

	url = fmt.Sprintf("%v/cryptocurrencies", e.BaseURL)
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] NewRequest")
		return
	}

	req = req.WithContext(ctx)
	resp, err = e.client.Do(req)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] client.Do")
		return
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] ReadAll")
		return
	}
	crypto := &core.CryptoClientResp{}
	err = json.Unmarshal(b, crypto)
	if err != nil {
		lg.WithError(err).Error("[GetCurrencies] Unmarshal")
		return
	}
	if !crypto.Success {
		lg.WithField("success", symbols.Success).Warn("[GetCurrencies] success")
		return
	}

	l = map[string]string{}
	for _, v := range crypto.Cryptocurrencies {
		l[v.Symbol] = v.Name
	}

	lg.Info("[GetCurrencies] ok")
	return l, err
}

// Exchange tries to get the exchange rate for the given currencies
//
// receives a ctx so the request can be cancelled if the original request is also cancelled
func (e Exchange) Exchange(ctx context.Context, from, to string, amount float64) (core.ConversionResp, error) {
	lg := log.WithField("pkg", "exchange")
	url := fmt.Sprintf("%v/convert", e.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		lg.WithError(err).Error("[Exchange] NewRequest")
		return core.ConversionResp{}, err
	}

	q := req.URL.Query()
	q.Add("from", from)
	q.Add("to", to)
	q.Add("amount", fmt.Sprintf("%f", amount))
	req.URL.RawQuery = q.Encode()

	req = req.WithContext(ctx)
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
		ConvertedAmount:  conv.Result,
		ConversionSource: "exchange",
	}, err
}
