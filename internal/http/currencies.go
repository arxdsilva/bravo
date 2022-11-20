package http

import (
	"net/http"

	"github.com/arxdsilva/bravo/internal/core"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// GetCurrencies retrieves currencies from DB and external exchange
//
// HTTP responses:
// 200 OK
// 500 Internal Server Error
func (s Server) GetCurrencies(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "GetCurrencies",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})
	currencies, err := s.service.GetCurrencies(c.Request().Context())
	if err != nil {
		lg.WithError(err).Error("service.GetCurrencies")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	lg.Info("success")
	return c.JSON(http.StatusOK, currencies)
}

// AddCurrency retrieves currencies from DB and external exchange
//
// HTTP responses:
// 200 OK
// 200 Bad Request
// 500 Internal Server Error
func (s Server) AddCurrency(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "AddCurrency",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})

	currency := &core.Currency{}

	if err = c.Bind(currency); err != nil {
		lg.WithError(err).Error("c.Bind")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = currency.Check(); err != nil {
		lg.WithError(err).Error("currency.Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.service.AddCurrency(c.Request().Context(), currency.Symbol, currency.Description)
	if err != nil {
		lg.WithError(err).Error("service.GetCurrencies")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	lg.Info("success")
	return c.JSON(http.StatusCreated, currency)
}
