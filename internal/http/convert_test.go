package http

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/arxdsilva/bravo/internal/core"
	rsv "github.com/arxdsilva/bravo/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_Convert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		wantToConv bool
		from, to   string
		amount     string

		wantConvSVC core.ConversionSVC

		convAmount      float64
		convSource      string
		convCurrencyErr error

		wantBody    string
		wantErrFn   require.ErrorAssertionFunc
		wantCode    int
		wantHTTPErr *echo.HTTPError
	}{
		{
			name:            "check error - symbol too small",
			from:            "AB",
			to:              "AB",
			amount:          "10",
			wantToConv:      false,
			convCurrencyErr: nil,
			wantBody:        "",
			wantErrFn:       require.Error,
			wantCode:        http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "currency symbol has to have 3 or more characters",
				Internal: nil,
			},
		}, {
			name:            "check error - symbol too small",
			from:            "ABC",
			to:              "AB",
			amount:          "10",
			wantToConv:      false,
			convCurrencyErr: nil,
			wantBody:        "",
			wantErrFn:       require.Error,
			wantCode:        http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "currency symbol has to have 3 or more characters",
				Internal: nil,
			},
		}, {
			name:            "amount check error",
			from:            "ABC",
			to:              "ABD",
			amount:          "10a",
			wantToConv:      false,
			convCurrencyErr: nil,
			wantBody:        "",
			wantErrFn:       require.Error,
			wantCode:        http.StatusBadRequest,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "amount is not a number",
				Internal: nil,
			},
		}, {
			name:   "should not convert - no err",
			from:   "ABC",
			to:     "ABC",
			amount: "10",

			wantToConv:      false,
			convCurrencyErr: nil,

			wantBody:  "{\"from\":\"ABC\",\"to\":\"ABC\",\"original_amount\":10,\"converted_amount\":10,\"conversion_source\":\"no-edit\"}\n",
			wantErrFn: require.NoError,
			wantCode:  http.StatusOK,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusOK,
				Message:  "",
				Internal: nil,
			},
		},
		{
			name:   "convert error",
			from:   "ABC",
			to:     "ABD",
			amount: "10",

			wantToConv:      true,
			convAmount:      15,
			convSource:      "exchange",
			convCurrencyErr: errors.New("some err"),
			wantConvSVC:     core.ConversionSVC{From: "ABC", To: "ABD", Amount: 10.0},

			wantBody:  "",
			wantErrFn: require.Error,
			wantCode:  http.StatusInternalServerError,
			wantHTTPErr: &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "some err",
				Internal: nil,
			},
		},
		{
			name:   "no error",
			from:   "ABC",
			to:     "ABD",
			amount: "10",

			wantToConv:      true,
			convAmount:      15,
			convSource:      "exchange",
			convCurrencyErr: nil,
			wantConvSVC:     core.ConversionSVC{From: "ABC", To: "ABD", Amount: 10.0},

			wantBody:    "{\"from\":\"ABC\",\"to\":\"ABD\",\"original_amount\":10,\"converted_amount\":15,\"conversion_source\":\"exchange\"}\n",
			wantErrFn:   require.NoError,
			wantCode:    http.StatusOK,
			wantHTTPErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := rsv.NewMockResolver(ctrl)

			vals := url.Values{}
			vals.Set("from", tt.from)
			vals.Set("to", tt.to)
			vals.Set("amount", tt.amount)

			req := httptest.NewRequest(http.MethodGet, "/convertion/convert?"+vals.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)
			ctx.SetPath("/convertion/convert")

			if tt.wantToConv {
				mock.EXPECT().Convert(gomock.Any(), tt.wantConvSVC).
					Return(tt.convAmount, tt.convSource, tt.convCurrencyErr)
			}

			s := Server{service: mock}
			err := s.Convert(ctx)
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
