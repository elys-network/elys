package types

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgOpen(t *testing.T) {

	accAdress := sample.AccAddress()
	got := NewMsgOpen(
		accAdress,
		Position_LONG,
		sdkmath.LegacyNewDec(200),
		"uatom",
		sdk.NewCoin("uusdc", sdkmath.NewInt(2000)),
		sdkmath.LegacyNewDec(100),
		sdkmath.LegacyNewDec(0),
	)

	want := &MsgOpen{
		Creator:         accAdress,
		Position:        Position_LONG,
		Leverage:        sdkmath.LegacyNewDec(200),
		TradingAsset:    "uatom",
		Collateral:      sdk.NewCoin("uusdc", sdkmath.NewInt(2000)),
		TakeProfitPrice: sdkmath.LegacyNewDec(100),
		StopLossPrice:   sdkmath.LegacyNewDec(0),
	}

	assert.Equal(t, want, got)
}

func TestMsgOpen_ValidateBasic(t *testing.T) {

	type Test struct {
		title string
		msg   MsgOpen
		want  error
	}

	tableTest := []Test{
		{
			title: "invalid address",
			msg: MsgOpen{
				Creator: "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "invalid position",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_UNSPECIFIED,
				TradingAsset: "",
			},
			want: ErrInvalidPosition,
		},
		{
			title: "leverage is nil",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "uatom",
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "invalid leverage",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "uatom",
				Leverage:     sdkmath.LegacyNewDec(-200),
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "trading asset is empty",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "",
				Leverage:     sdkmath.LegacyNewDec(200),
			},
			want: ErrTradingAssetIsEmpty,
		},
		{
			title: "take profit price is nil",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_SHORT,
				TradingAsset: "uatom",
				Leverage:     sdkmath.LegacyNewDec(200),
			},
			want: ErrInvalidTakeProfitPriceIsNegative,
		},
		{
			title: "take profit price is negative",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_SHORT,
				TradingAsset:    "uatom",
				TakeProfitPrice: sdkmath.LegacyNewDec(-10),
				Leverage:        sdkmath.LegacyNewDec(200),
			},
			want: ErrInvalidTakeProfitPriceIsNegative,
		},
		{
			title: "successful",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_LONG,
				TradingAsset:    "uatom",
				TakeProfitPrice: sdkmath.LegacyNewDec(10),
				Leverage:        sdkmath.LegacyNewDec(200),
			},
			want: nil,
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			got := test.msg.ValidateBasic()

			if got != nil {
				assert.ErrorIs(t, got, test.want)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}

}
