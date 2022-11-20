package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) RouterRegister(e *echo.Echo) {
	e.GET("/", HealthCheck)
	e.GET("/convertion/convert", s.Convert)
	// currency management
	e.GET("/currencies", s.GetCurrencies)           // get all allowed currencies
	e.POST("/currencies", s.AddCurrency)            // register new allowed currency
	e.GET("/currencies/{currency}", HealthCheck)    // gets a currency infos
	e.PUT("/currencies/{currency}", HealthCheck)    // register new allowed currency
	e.DELETE("/currencies/{currency}", HealthCheck) // remove allowed currency
	// currency rate management
	e.POST("/convertion/{currency}/rate", HealthCheck)
	e.PUT("/convertion/{currency}/rate", HealthCheck)
}

// todo: allow this to be configurable and to pass optional checks
// ex: db, services, connections ...
func HealthCheck(c echo.Context) (err error) {
	ok := struct {
		Service string `json:"service"`
	}{"ok"}
	return c.JSON(http.StatusOK, ok)
}
