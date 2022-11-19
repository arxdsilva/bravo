package http

import (
	"net/http"

	"github.com/arxdsilva/bravo/internal/core"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (s Server) Convert(c echo.Context) (err error) {
	conv := core.ConversionAPI{
		From:   c.QueryParam("from"),
		To:     c.QueryParam("to"),
		Amount: c.QueryParam("amount"),
	}

	if err := conv.Check(); err != nil {
		log.WithError(err).Error("[Convert] Check")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	convService, shouldConvert, err := core.ConvertToService(conv)
	if err != nil {
		log.WithError(err).Error("[Convert] ConvertToService")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !shouldConvert {
		log.WithField("should_convert", shouldConvert).Info("[Convert] should_convert")
		return c.JSON(http.StatusOK, core.TransformSVCToResp(convService, convService.Amount, "no-edit"))
	}

	amount, source, err := s.service.Convert(c.Request().Context(), convService)
	if err != nil {
		log.WithError(err).Error("[Convert] service.Convert")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Info("[Convert] success")
	return c.JSON(http.StatusOK, core.TransformSVCToResp(convService, amount, source))
}
