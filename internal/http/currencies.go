package http

import (
	"errors"
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
// 201 Created
// 400 Bad Request
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

	err = s.service.AddCurrency(
		c.Request().Context(), currency.Symbol, currency.Description)
	if err != nil {
		lg.WithError(err).Error("service.AddCurrency")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lg.Info("success")
	return c.JSON(http.StatusCreated, currency)
}

// UpdateCurrency retrieves currencies from DB and external exchange
//
// HTTP responses:
// 204 No Content
// 400 Bad Request
// 500 Internal Server Error
func (s Server) UpdateCurrency(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "UpdateCurrency",
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

	err = s.service.UpdateCurrency(
		c.Request().Context(), currency.Symbol, currency.Description)
	if err != nil && err != core.ErrNotFound {
		lg.WithError(err).Error("service.UpdateCurrency")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if errors.Is(err, core.ErrNotFound) {
		lg.WithError(core.ErrCurrencyNotFound).Error("not found")
		return echo.NewHTTPError(http.StatusBadRequest, core.ErrCurrencyNotFound.Error())
	}

	lg.Info("success")
	return c.JSON(http.StatusCreated, currency)
}

// GetCurrency retrieves a currency from DB
//
// HTTP responses:
// 200 OK
// 400 Bad Request
// 404 Not Found
// 500 Internal Server Error
func (s Server) GetCurrency(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "GetCurrency",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})
	symbol := c.Param("symbol")

	cr := &core.Currency{Symbol: symbol}

	if err = cr.Check(); err != nil {
		lg.WithError(err).Error("currency.Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	svcCurrency, err := s.service.GetCurrency(c.Request().Context(), cr.Symbol)
	if err != nil && err != core.ErrNotFound {
		lg.WithError(err).Error("service.GetCurrency")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if errors.Is(err, core.ErrNotFound) {
		lg.WithError(core.ErrCurrencyNotFound).Error("not found")
		return echo.NewHTTPError(http.StatusBadRequest, core.ErrCurrencyNotFound.Error())
	}

	lg.Info("success")
	return c.JSON(http.StatusOK, svcCurrency)
}

// RemoveCurrency retrieves a currency from DB
//
// HTTP responses:
// 204 No Content
// 400 Bad Request
// 404 Not Found
// 500 Internal Server Error
func (s Server) RemoveCurrency(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "RemoveCurrency",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})
	symbol := c.Param("symbol")

	cr := &core.Currency{Symbol: symbol}

	if err = cr.Check(); err != nil {
		lg.WithError(err).Error("currency.Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.service.RemoveCurrency(c.Request().Context(), cr.Symbol)
	if err != nil && err != core.ErrNotFound {
		lg.WithError(err).Error("service.RemoveCurrency")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if errors.Is(err, core.ErrNotFound) {
		lg.WithError(core.ErrCurrencyNotFound).Error("not found")
		return echo.NewHTTPError(http.StatusBadRequest, core.ErrCurrencyNotFound.Error())
	}

	lg.Info("success")
	return c.NoContent(http.StatusNoContent)
}
