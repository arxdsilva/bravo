package http

import (
	"net/http"

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
		"route": "convert",
	})
	currencies, err := s.service.GetCurrencies(c.Request().Context())
	if err != nil {
		lg.WithError(err).Error("service.GetCurrencies")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	lg.Info("success")
	return c.JSON(http.StatusOK, currencies)
}
