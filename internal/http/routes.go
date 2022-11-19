package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RouterRegister(e *echo.Echo) {
	e.GET("/", HealthCheck)
}

// todo: allow this to be configurable and to pass optional checks
// ex: db, services, connections ...
func HealthCheck(c echo.Context) (err error) {
	ok := struct {
		Service string `json:"service"`
	}{"ok"}
	return c.JSON(http.StatusOK, ok)
}
