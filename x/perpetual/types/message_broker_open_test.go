package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgBrokerOpen(t *testing.T) {

	accAdress := sample.AccAddress()
	owner := sample.AccAddress()
	got := NewMsgBrokerOpen(
		accAdress,
		Position_LONG,
		math.LegacyNewDec(200),
		"uatom",
		sdk.NewCoin("uusdc", sdk.NewInt(2000)),
		math.LegacyNewDec(100),
		owner,
		sdkmath.LegacyZeroDec(),
	)

	want := &MsgBrokerOpen{
		Creator:         accAdress,
		Position:        Position_LONG,
		Leverage:        math.LegacyNewDec(200),
		TradingAsset:    "uatom",
		Collateral:      sdk.NewCoin("uusdc", sdk.NewInt(2000)),
		TakeProfitPrice: math.LegacyNewDec(100),
		Owner:           owner,
		StopLossPrice:   sdkmath.LegacyZeroDec(),
	}

	assert.Equal(t, want, got)
}

func TestMsgBrokerOpen_ValidateBasic(t *testing.T) {

	type Test struct {
		title string
		msg   MsgBrokerOpen
		want  error
	}

	tableTest := []Test{
		{
			title: "invalid address",
			msg: MsgBrokerOpen{
				Creator: "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "invalid owner",
			msg: MsgBrokerOpen{
				Creator: sample.AccAddress(),
				Owner:   "invalid",
			},
			want: sdkerrors.ErrInvalidAddress,
		},
		{
			title: "invalid position",
			msg: MsgBrokerOpen{
				Creator:      sample.AccAddress(),
				Owner:        sample.AccAddress(),
				Position:     Position_UNSPECIFIED,
				TradingAsset: "",
			},
			want: ErrInvalidPosition,
		},
		{
			title: "leverage is nil",
			msg: MsgBrokerOpen{
				Creator:      sample.AccAddress(),
				Owner:        sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "uatom",
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "invalid leverage",
			msg: MsgBrokerOpen{
				Creator:      sample.AccAddress(),
				Owner:        sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "uatom",
				Leverage:     math.LegacyNewDec(-200),
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "trading asset is empty",
			msg: MsgBrokerOpen{
				Creator:      sample.AccAddress(),
				Owner:        sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "",
				Leverage:     math.LegacyNewDec(200),
			},
			want: ErrInvalidTradingAsset,
		},
		{
			title: "take profit price is nil",
			msg: MsgBrokerOpen{
				Creator:      sample.AccAddress(),
				Owner:        sample.AccAddress(),
				Position:     Position_SHORT,
				TradingAsset: "uatom",
				Leverage:     math.LegacyNewDec(200),
			},
			want: ErrInvalidTakeProfitPrice,
		},
		{
			title: "take profit price is negative",
			msg: MsgBrokerOpen{
				Creator:         sample.AccAddress(),
				Owner:           sample.AccAddress(),
				Position:        Position_SHORT,
				TradingAsset:    "uatom",
				TakeProfitPrice: math.LegacyNewDec(-10),
				Leverage:        math.LegacyNewDec(200),
			},
			want: ErrInvalidTakeProfitPrice,
		},
		{
			title: "stop loss price is nil",
			msg: MsgBrokerOpen{
				Creator:         sample.AccAddress(),
				Owner:           sample.AccAddress(),
				Position:        Position_SHORT,
				TradingAsset:    "uatom",
				Leverage:        math.LegacyNewDec(200),
				TakeProfitPrice: math.LegacyNewDec(10),
			},
			want: errorsmod.Wrapf(ErrInvalidPrice, "stopLossPrice cannot be nil"),
		},
		{
			title: "stop loss price is negative",
			msg: MsgBrokerOpen{
				Creator:         sample.AccAddress(),
				Owner:           sample.AccAddress(),
				Position:        Position_SHORT,
				TradingAsset:    "uatom",
				TakeProfitPrice: math.LegacyNewDec(10),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(-10),
			},
			want: errorsmod.Wrapf(ErrInvalidPrice, "stopLossPrice cannot be negative"),
		},
		{
			title: "take profit price is greater than stop loss price for short",
			msg: MsgBrokerOpen{
				Creator:         sample.AccAddress(),
				Owner:           sample.AccAddress(),
				Position:        Position_SHORT,
				TradingAsset:    "uatom",
				TakeProfitPrice: math.LegacyNewDec(110),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(100),
			},
			want: fmt.Errorf("TakeProfitPrice cannot be >= StopLossPrice for SHORT"),
		},
		{
			title: "take profit price is less than stop loss price for long",
			msg: MsgBrokerOpen{
				Creator:         sample.AccAddress(),
				Owner:           sample.AccAddress(),
				Position:        Position_LONG,
				TradingAsset:    "uatom",
				TakeProfitPrice: math.LegacyNewDec(90),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(100),
			},
			want: fmt.Errorf("TakeProfitPrice cannot be <= StopLossPrice for LONG"),
		},
		{
			title: "successful",
			msg: MsgBrokerOpen{
				Creator:         sample.AccAddress(),
				Owner:           sample.AccAddress(),
				Position:        Position_LONG,
				TradingAsset:    "uatom",
				TakeProfitPrice: math.LegacyNewDec(300),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(100),
			},
			want: nil,
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			got := test.msg.ValidateBasic()

			if got != nil {
				require.Error(t, got, test.want)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}

}
