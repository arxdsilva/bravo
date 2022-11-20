package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) RouterRegister(e *echo.Echo) {
	e.GET("/", HealthCheck)
	e.GET("/convertion/convert", s.Convert)
	// currency management
	e.GET("/currencies", s.GetCurrencies)
	e.POST("/currencies", s.AddCurrency)
	e.GET("/currencies/:symbol", s.GetCurrency)
	e.PUT("/currencies/:symbol", s.UpdateCurrency)
	e.DELETE("/currencies/:symbol", s.RemoveCurrency)
	// currency rate management
	e.POST("/convertion/:symbol/rate", HealthCheck)
	e.PUT("/convertion/:symbol/rate", HealthCheck)
}

// todo: allow this to be configurable and to pass optional checks
// ex: db, services, connections ...
func HealthCheck(c echo.Context) (err error) {
	ok := struct {
		Service string `json:"service"`
	}{"ok"}
	return c.JSON(http.StatusOK, ok)
}
