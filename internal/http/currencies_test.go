package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func Test_AddCurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		wantToAdd    bool
		sentCurrency core.Currency

		addCurrencyErr error

		wantBody    string
		wantErrFn   require.ErrorAssertionFunc
		wantCode    int
		wantHTTPErr *echo.HTTPError
	}{
		{
			name:      "check error - no symbol",
			wantToAdd: false,
			sentCurrency: core.Currency{
				Symbol: "",
			},
			addCurrencyErr: nil,
			wantBody:       "",
			wantErrFn:      require.Error,
			wantCode:       http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "currency needs a symbol",
				Internal: nil,
			},
		},
		{
			name:      "check error - symbol too small",
			wantToAdd: false,
			sentCurrency: core.Currency{
				Symbol: "BR",
			},
			addCurrencyErr: nil,
			wantBody:       "",
			wantErrFn:      require.Error,
			wantCode:       http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "currency symbol has to have 3 or more characters",
				Internal: nil,
			},
		},
		{
			name:      "add error",
			wantToAdd: true,
			sentCurrency: core.Currency{
				Symbol: "BRL",
			},
			addCurrencyErr: errors.New("some err"),
			wantBody:       "",
			wantErrFn:      require.Error,
			wantCode:       http.StatusInternalServerError,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "some err",
				Internal: nil,
			},
		},
		{
			name:      "no error",
			wantToAdd: true,
			sentCurrency: core.Currency{
				Symbol: "BRL",
			},
			addCurrencyErr: nil,
			wantBody:       "{\"symbol\":\"BRL\",\"description\":\"\",\"source\":\"\"}\n",
			wantErrFn:      require.NoError,
			wantCode:       http.StatusCreated,
			wantHTTPErr:    nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := rsv.NewMockResolver(ctrl)

			b, err := json.Marshal(tt.sentCurrency)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/currencies", strings.NewReader(string(b)))
			require.NoError(t, err)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)
			ctx.SetPath("/currencies")

			if tt.wantToAdd {
				mock.EXPECT().AddCurrency(gomock.Any(), tt.sentCurrency.Symbol, tt.sentCurrency.Description).
					Return(tt.addCurrencyErr)
			}

			s := Server{service: mock}
			err = s.AddCurrency(ctx)
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
