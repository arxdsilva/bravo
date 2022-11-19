package http

import (
	"net/http"

	"github.com/arxdsilva/bravo/internal/core"
	"github.com/labstack/echo/v4"
)

func (s Server) Convert(c echo.Context) (err error) {
	conv := core.ConversionAPI{
		From:   c.QueryParam("from"),
		To:     c.QueryParam("from"),
		Amount: c.QueryParam("amount"),
	}

	if err := conv.Check(); err != nil {
		// log
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	convService, shouldConvert, err := core.ConvertToService(conv)
	if err != nil {
		// log
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if !shouldConvert {
		// log
		return c.JSON(http.StatusOK, convService)
	}

	err = s.service.Convert(c.Request().Context(), convService)
	if err != nil {
		// log
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// log
	return c.JSON(http.StatusOK, "ok")
}
