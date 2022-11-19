package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConversionAPI_Check(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		api           ConversionAPI
		wantErrFn     require.ErrorAssertionFunc
		wantErrEquals error
	}{
		{
			name:          "empty should return From err",
			api:           ConversionAPI{},
			wantErrFn:     require.Error,
			wantErrEquals: ErrInvalidFromCurrency,
		},
		{
			name: "should return To err",
			api: ConversionAPI{
				From: "USD",
			},
			wantErrFn:     require.Error,
			wantErrEquals: ErrInvalidToCurrency,
		},
		{
			name: "invalid To err",
			api: ConversionAPI{
				From: "USD",
				To:   "LALA",
			},
			wantErrFn:     require.Error,
			wantErrEquals: ErrInvalidToCurrency,
		},
		{
			name: "invalid From err",
			api: ConversionAPI{
				From: "USDa",
			},
			wantErrFn:     require.Error,
			wantErrEquals: ErrInvalidFromCurrency,
		},
		{
			name: "amount err",
			api: ConversionAPI{
				From: "USD",
				To:   "BRL",
			},
			wantErrFn:     require.Error,
			wantErrEquals: ErrAmountIsNotANumber,
		},
		{
			name: "amount err not float",
			api: ConversionAPI{
				From:   "USD",
				To:     "BRL",
				Amount: "amount",
			},
			wantErrFn:     require.Error,
			wantErrEquals: ErrAmountIsNotANumber,
		},
		{
			name: "no err",
			api: ConversionAPI{
				From:   "USD",
				To:     "BRL",
				Amount: "100",
			},
			wantErrFn:     require.NoError,
			wantErrEquals: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.api.Check()
			tt.wantErrFn(t, err)
			require.Equal(t, tt.wantErrEquals, err)
		})
	}
}
