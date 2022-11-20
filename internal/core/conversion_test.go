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
			name:          "empty should return min len err",
			api:           ConversionAPI{},
			wantErrFn:     require.Error,
			wantErrEquals: ErrSymbolMinLen,
		},
		{
			name:          "should return min len err",
			api:           ConversionAPI{From: "USD"},
			wantErrFn:     require.Error,
			wantErrEquals: ErrSymbolMinLen,
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

func Test_ConvertToService(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		c          ConversionAPI
		wantCs     ConversionSVC
		wantShould bool
		wantErrFn  require.ErrorAssertionFunc
	}{
		{
			name:       "empty should return amount err",
			c:          ConversionAPI{},
			wantCs:     ConversionSVC{},
			wantShould: false,
			wantErrFn:  require.Error,
		},
		{
			name: "same currencies should not error",
			c: ConversionAPI{
				From:   "BRL",
				To:     "BRL",
				Amount: "1234",
			},
			wantCs: ConversionSVC{
				From:   "BRL",
				To:     "BRL",
				Amount: 1234,
			},
			wantShould: false,
			wantErrFn:  require.NoError,
		},
		{
			name: "different currencies should not error, with should true",
			c: ConversionAPI{
				From:   "USD",
				To:     "BRL",
				Amount: "1234",
			},
			wantCs: ConversionSVC{
				From:   "USD",
				To:     "BRL",
				Amount: 1234,
			},
			wantShould: true,
			wantErrFn:  require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCs, gotShould, err := ConvertToService(tt.c)
			tt.wantErrFn(t, err)
			require.Equal(t, tt.wantCs, gotCs)
			require.Equal(t, tt.wantShould, gotShould)
		})
	}
}
