package http

import (
	"net/http"

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
