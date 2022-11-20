package http

import (
	"net/http"

	"github.com/arxdsilva/bravo/internal/core"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Convert retrieves a conversion from two currencies
//
// HTTP responses:
// 200 OK
// 400 Bad request
// 500 Internal Server Error
func (s Server) Convert(c echo.Context) (err error) {
	lg := log.WithFields(log.Fields{
		"pkg":   "http",
		"route": "convert",
		"cid":   c.Response().Header().Get(echo.HeaderXRequestID),
	})
	conv := core.ConversionAPI{
		From:   c.QueryParam("from"),
		To:     c.QueryParam("to"),
		Amount: c.QueryParam("amount"),
	}

	if err := conv.Check(); err != nil {
		lg.WithError(err).Error("check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	convService, shouldConvert, err := core.ConvertToService(conv)
	if err != nil {
		lg.WithError(err).Error("convertToService")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !shouldConvert {
		lg.WithField("should_convert", shouldConvert).Info("should_convert")
		return c.JSON(http.StatusOK, core.TransformSVCToResp(convService, convService.Amount, "no-edit"))
	}

	amount, source, err := s.service.Convert(c.Request().Context(), convService)
	if err != nil {
		lg.WithError(err).Error("service.Convert")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	lg.Info("success")
	return c.JSON(http.StatusOK, core.TransformSVCToResp(convService, amount, source))
}
