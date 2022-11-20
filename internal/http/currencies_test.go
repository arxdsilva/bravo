package http

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arxdsilva/bravo/internal/core"
	rsv "github.com/arxdsilva/bravo/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_GetCurrencies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		getCurrenciesResp core.Currencies
		getCurrenciesErr  error

		wantBody    string
		wantErrFn   require.ErrorAssertionFunc
		wantCode    int
		wantHTTPErr *echo.HTTPError
	}{
		{
			name:              "retrieve error",
			getCurrenciesResp: core.Currencies{},
			getCurrenciesErr:  errors.New("some err"),
			wantBody:          "",
			wantErrFn:         require.Error,
			wantCode:          http.StatusInternalServerError,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "some err",
				Internal: nil,
			},
		},
		{
			name: "no error",
			getCurrenciesResp: core.Currencies{
				core.Currency{Symbol: "BRL"},
			},
			getCurrenciesErr: nil,
			wantBody:         "[{\"symbol\":\"BRL\",\"description\":\"\",\"source\":\"\"}]\n",
			wantErrFn:        require.NoError,
			wantCode:         http.StatusOK,
			wantHTTPErr:      nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := rsv.NewMockResolver(ctrl)

			req, err := http.NewRequest(http.MethodGet, "/currencies", nil)
			require.NoError(t, err)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)
			ctx.SetPath("/currencies")

			mock.EXPECT().GetCurrencies(gomock.Any()).Return(tt.getCurrenciesResp, tt.getCurrenciesErr)

			s := Server{service: mock}
			err = s.GetCurrencies(ctx)
			tt.wantErrFn(t, err)
			if err == nil {
				require.Equal(t, tt.wantCode, rec.Code)
				b, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.Equal(t, tt.wantBody, string(b))
				return
			}
			require.Equal(t, tt.wantHTTPErr, err)
		})
	}
}
