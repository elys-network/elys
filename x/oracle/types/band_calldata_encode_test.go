package types

import (
	"testing"

	"github.com/bandprotocol/bandchain-packet/obi"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestBandPriceCallDataEncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		data BandPriceCallData
		err  error
	}{
		{
			name: "empty symbols",
			data: BandPriceCallData{
				Symbols:    []string{},
				Multiplier: 1,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "zero multiplier",
			data: BandPriceCallData{
				Symbols:    []string{},
				Multiplier: 0,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "normal data",
			data: BandPriceCallData{
				Symbols:    []string{"BTC", "ETH"},
				Multiplier: 18,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodedCalldata := obi.MustEncode(tt.data)
			var request BandPriceCallData
			err := obi.Decode(encodedCalldata, &request)
			require.NoError(t, err)
			require.Equal(t, request, tt.data)
		})
	}
}
