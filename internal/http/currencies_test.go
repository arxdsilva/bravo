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

func Test_UpdateCurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		wantToAdd    bool
		sentCurrency core.Currency

		updateCurrencyErr error

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
			updateCurrencyErr: nil,
			wantBody:          "",
			wantErrFn:         require.Error,
			wantCode:          http.StatusBadRequest,
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
			updateCurrencyErr: nil,
			wantBody:          "",
			wantErrFn:         require.Error,
			wantCode:          http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "currency symbol has to have 3 or more characters",
				Internal: nil,
			},
		},
		{
			name:      "update error",
			wantToAdd: true,
			sentCurrency: core.Currency{
				Symbol: "BRL",
			},
			updateCurrencyErr: errors.New("some err"),
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
			name:      "not found error",
			wantToAdd: true,
			sentCurrency: core.Currency{
				Symbol: "BRL",
			},
			updateCurrencyErr: core.ErrNotFound,
			wantBody:          "",
			wantErrFn:         require.Error,
			wantCode:          http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  core.ErrCurrencyNotFound.Error(),
				Internal: nil,
			},
		},
		{
			name:      "no error",
			wantToAdd: true,
			sentCurrency: core.Currency{
				Symbol: "BRL",
			},
			updateCurrencyErr: nil,
			wantBody:          "{\"symbol\":\"BRL\",\"description\":\"\",\"source\":\"\"}\n",
			wantErrFn:         require.NoError,
			wantCode:          http.StatusCreated,
			wantHTTPErr:       nil,
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
				mock.EXPECT().UpdateCurrency(gomock.Any(), tt.sentCurrency.Symbol, tt.sentCurrency.Description).
					Return(tt.updateCurrencyErr)
			}

			s := Server{service: mock}
			err = s.UpdateCurrency(ctx)
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

func Test_GetCurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		wantToAdd    bool
		sentCurrency string

		getCurrencyErr  error
		getCurrencyResp core.Currency

		wantBody    string
		wantErrFn   require.ErrorAssertionFunc
		wantCode    int
		wantHTTPErr *echo.HTTPError
	}{
		{
			name:            "check error - symbol too small",
			wantToAdd:       false,
			sentCurrency:    "BR",
			getCurrencyErr:  nil,
			getCurrencyResp: core.Currency{},
			wantBody:        "",
			wantErrFn:       require.Error,
			wantCode:        http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "currency symbol has to have 3 or more characters",
				Internal: nil,
			},
		},
		{
			name:            "add error",
			wantToAdd:       true,
			sentCurrency:    "BRL",
			getCurrencyResp: core.Currency{},
			getCurrencyErr:  errors.New("some err"),
			wantBody:        "",
			wantErrFn:       require.Error,
			wantCode:        http.StatusInternalServerError,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "some err",
				Internal: nil,
			},
		},
		{
			name:            "not found error",
			wantToAdd:       true,
			sentCurrency:    "BRL",
			getCurrencyResp: core.Currency{},
			getCurrencyErr:  core.ErrNotFound,
			wantBody:        "",
			wantErrFn:       require.Error,
			wantCode:        http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  core.ErrCurrencyNotFound.Error(),
				Internal: nil,
			},
		},
		{
			name:      "no error",
			wantToAdd: true,
			getCurrencyResp: core.Currency{
				Symbol: "BRL",
			},
			sentCurrency:   "BRL",
			getCurrencyErr: nil,
			wantBody:       "{\"symbol\":\"BRL\",\"description\":\"\",\"source\":\"\"}\n",
			wantErrFn:      require.NoError,
			wantCode:       http.StatusOK,
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
			ctx.SetPath("/currencies/:symbol")
			ctx.SetParamNames("symbol")
			ctx.SetParamValues(tt.sentCurrency)

			if tt.wantToAdd {
				mock.EXPECT().GetCurrency(gomock.Any(), tt.sentCurrency).
					Return(tt.getCurrencyResp, tt.getCurrencyErr)
			}

			s := Server{service: mock}
			err = s.GetCurrency(ctx)
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

func Test_RemoveCurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		wantToAdd    bool
		sentCurrency string

		rmvCurrencyErr error

		wantBody    string
		wantErrFn   require.ErrorAssertionFunc
		wantCode    int
		wantHTTPErr *echo.HTTPError
	}{
		{
			name:           "check error - symbol too small",
			wantToAdd:      false,
			sentCurrency:   "BR",
			rmvCurrencyErr: nil,
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
			name:           "add error",
			wantToAdd:      true,
			sentCurrency:   "BRL",
			rmvCurrencyErr: errors.New("some err"),
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
			name:           "not found error",
			wantToAdd:      true,
			sentCurrency:   "BRL",
			rmvCurrencyErr: core.ErrNotFound,
			wantBody:       "",
			wantErrFn:      require.Error,
			wantCode:       http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  core.ErrCurrencyNotFound.Error(),
				Internal: nil,
			},
		},
		{
			name:           "no error",
			wantToAdd:      true,
			sentCurrency:   "BRL",
			rmvCurrencyErr: nil,
			wantBody:       "",
			wantErrFn:      require.NoError,
			wantCode:       http.StatusNoContent,
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
			ctx.SetPath("/currencies/:symbol")
			ctx.SetParamNames("symbol")
			ctx.SetParamValues(tt.sentCurrency)

			if tt.wantToAdd {
				mock.EXPECT().RemoveCurrency(gomock.Any(), tt.sentCurrency).
					Return(tt.rmvCurrencyErr)
			}

			s := Server{service: mock}
			err = s.RemoveCurrency(ctx)
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
