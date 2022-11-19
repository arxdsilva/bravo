package http

import (
	"github.com/arxdsilva/bravo/internal/jwt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return customGenerator()
		},
	}))
}

func customGenerator() string {
	return uuid.New().String()
}

func JWT(tp jwt.TokenProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := tp.ValidateUserUUID(c.Request()) // set this accordingly
			if err == nil {
				c.Set("user_uuid", user)
				return next(c)
			}
			log.Infof("[JWT] Middleware error: %v", err.Error())
			return &echo.HTTPError{
				Code:     middleware.ErrJWTInvalid.Code,
				Message:  middleware.ErrJWTInvalid.Message,
				Internal: err,
			}
		}
	}
}
