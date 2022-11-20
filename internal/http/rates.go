package http

import (
	"errors"
	"net/http"

	"github.com/arxdsilva/bravo/internal/core"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// GetRates retrieves all currency rates in DB
//
// HTTP responses:
// 200 OK
// 500 Internal Server Error
func (s Server) GetRates(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "GetRates",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})
	currencies, err := s.service.GetRates(c.Request().Context())
	if err != nil {
		lg.WithError(err).Error("service.GetRates")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	lg.Info("success")
	return c.JSON(http.StatusOK, currencies)
}

// CreateRate tries to create a currency rate into DB
//
// HTTP responses:
// 201 Created
// 400 Bad Request
// 500 Internal Server Error
func (s Server) CreateRate(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "CreateRate",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})

	rate := &core.CurrencyRate{}
	if err = c.Bind(rate); err != nil {
		lg.WithError(err).Error("c.Bind")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = rate.Check(); err != nil {
		lg.WithError(err).Error("rate.Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.service.CreateRate(
		c.Request().Context(), rate.From, rate.To, rate.Rate)
	if err != nil { // check not found err
		lg.WithError(err).Error("service.CreateRate")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lg.Info("success")
	return c.JSON(http.StatusCreated, rate)
}

// UpdateRate tries to create a currency rate into DB
//
// HTTP responses:
// 202 Accepted
// 400 Bad Request
// 500 Internal Server Error
func (s Server) UpdateRate(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "UpdateRate",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})

	rate := &core.CurrencyRate{}
	if err = c.Bind(rate); err != nil {
		lg.WithError(err).Error("c.Bind")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = rate.Check(); err != nil {
		lg.WithError(err).Error("rate.Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.service.UpdateRate(
		c.Request().Context(), rate.From, rate.To, rate.Rate)
	if err != nil && err != core.ErrNotFound { // check not found err
		lg.WithError(err).Error("service.UpdateRate")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if errors.Is(err, core.ErrNotFound) {
		lg.WithError(core.ErrCurrencyNotFound).Error("not found")
		return echo.NewHTTPError(http.StatusBadRequest, core.ErrCurrencyNotFound.Error())
	}

	lg.Info("success")
	return c.JSON(http.StatusAccepted, rate)
}

// RemoveRate tries to create a currency rate into DB
//
// HTTP responses:
// 204 No Content
// 400 Bad Request
// 500 Internal Server Error
func (s Server) RemoveRate(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "RemoveRate",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})

	rate := &core.CurrencyRate{}
	if err = c.Bind(rate); err != nil {
		lg.WithError(err).Error("c.Bind")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = rate.Check(); err != nil {
		lg.WithError(err).Error("rate.Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.service.RemoveRate(
		c.Request().Context(), rate.From, rate.To)
	if err != nil && err != core.ErrNotFound { // check not found err
		lg.WithError(err).Error("service.RemoveRate")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if errors.Is(err, core.ErrNotFound) {
		lg.WithError(core.ErrCurrencyNotFound).Error("not found")
		return echo.NewHTTPError(http.StatusBadRequest, core.ErrCurrencyNotFound.Error())
	}

	lg.Info("success")
	return c.NoContent(http.StatusNoContent)
}
