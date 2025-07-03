package types

import (
	"errors"
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/assert"
)

func TestNewMsgOpen(t *testing.T) {

	accAdress := sample.AccAddress()
	got := NewMsgOpen(
		accAdress,
		Position_LONG,
		math.LegacyNewDec(200),
		1,
		sdk.NewCoin("uusdc", math.NewInt(2000)),
		math.LegacyNewDec(100),
		math.LegacyNewDec(0),
	)

	want := &MsgOpen{
		Creator:         accAdress,
		Position:        Position_LONG,
		Leverage:        math.LegacyNewDec(200),
		PoolId:          1,
		Collateral:      sdk.NewCoin("uusdc", math.NewInt(2000)),
		TakeProfitPrice: math.LegacyNewDec(100),
		StopLossPrice:   math.LegacyNewDec(0),
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
				Creator:  sample.AccAddress(),
				Position: Position_UNSPECIFIED,
			},
			want: ErrInvalidPosition,
		},
		{
			title: "leverage is nil",
			msg: MsgOpen{
				Creator:  sample.AccAddress(),
				Position: Position_LONG,
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "invalid leverage",
			msg: MsgOpen{
				Creator:  sample.AccAddress(),
				Position: Position_LONG,
				Leverage: math.LegacyNewDec(-200),
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "trading asset is empty",
			msg: MsgOpen{
				Creator:  sample.AccAddress(),
				Position: Position_LONG,
				Leverage: math.LegacyNewDec(200),
			},
			want: ErrInvalidTradingAsset,
		},
		{
			title: "take profit price is nil",
			msg: MsgOpen{
				Creator:  sample.AccAddress(),
				Position: Position_SHORT,
				Leverage: math.LegacyNewDec(200),
			},
			want: ErrInvalidTakeProfitPrice,
		},
		{
			title: "take profit price is negative",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_SHORT,
				TakeProfitPrice: math.LegacyNewDec(-10),
				Leverage:        math.LegacyNewDec(200),
			},
			want: ErrInvalidTakeProfitPrice,
		},
		{
			title: "stop loss price is nil",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_SHORT,
				Leverage:        math.LegacyNewDec(200),
				TakeProfitPrice: math.LegacyNewDec(10),
			},
			want: errorsmod.Wrapf(ErrInvalidPrice, "stopLossPrice cannot be nil"),
		},
		{
			title: "stop loss price is negative",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_SHORT,
				TakeProfitPrice: math.LegacyNewDec(10),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(-10),
			},
			want: errorsmod.Wrapf(ErrInvalidPrice, "stopLossPrice cannot be negative"),
		},
		{
			title: "take profit price is greater than stop loss price for short",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_SHORT,
				TakeProfitPrice: math.LegacyNewDec(110),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(100),
			},
			want: errors.New("TakeProfitPrice cannot be >= StopLossPrice for SHORT"),
		},
		{
			title: "take profit price is less than stop loss price for long",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_LONG,
				TakeProfitPrice: math.LegacyNewDec(90),
				Leverage:        math.LegacyNewDec(200),
				StopLossPrice:   math.LegacyNewDec(100),
			},
			want: errors.New("TakeProfitPrice cannot be <= StopLossPrice for LONG"),
		},
		{
			title: "successful",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_LONG,
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

func TestMsgOpen_GetCustodyAndLiabilitiesAsset(t *testing.T) {

	type Test struct {
		title string
		msg   MsgOpen
		want  error
	}

	tableTest := []Test{
		{
			title: "success: LONG",
			msg: MsgOpen{
				Position: Position_LONG,
			},
			want: nil,
		},
		{
			title: "success: SHORT",
			msg: MsgOpen{
				Position: Position_SHORT,
			},
			want: nil,
		},
		{
			title: "invalid position",
			msg: MsgOpen{
				Position: Position_UNSPECIFIED,
			},
			want: ErrInvalidPosition,
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			custodyAsset, liabilitiesAsset, err := test.msg.GetCustodyAndLiabilitiesAsset("uatom", "uusdc")

			if test.want != nil {
				require.Error(t, err, test.want)
			} else {
				assert.Nil(t, err)
				if test.msg.Position == Position_LONG {
					assert.Equal(t, "uatom", custodyAsset)
					assert.Equal(t, "uusdc", liabilitiesAsset)
				} else {
					assert.Equal(t, "uusdc", custodyAsset)
					assert.Equal(t, "uatom", liabilitiesAsset)
				}
			}
		})
	}

}

func TestMsgOpen_ValidatePosition(t *testing.T) {

	type Test struct {
		title string
		msg   MsgOpen
		want  error
	}

	tableTest := []Test{
		{
			title: "LONG",
			msg: MsgOpen{
				Position:   Position_LONG,
				Collateral: sdk.NewCoin("uelys", math.NewInt(1)),
			},
			want: errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid operation: collateral asset has to be either trading asset or base currency for long"),
		},
		{
			title: "SHORT",
			msg: MsgOpen{
				Position:   Position_SHORT,
				Collateral: sdk.NewCoin("uelys", math.NewInt(1)),
			},
			want: errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid collateral: collateral asset for short position must be the base currency"),
		},
		{
			title: "invalid position",
			msg: MsgOpen{
				Position: Position_UNSPECIFIED,
			},
			want: ErrInvalidPosition,
		},
		{
			title: "success: LONG",
			msg: MsgOpen{
				Position:   Position_LONG,
				Collateral: sdk.NewCoin("uusdc", math.NewInt(1)),
			},
			want: nil,
		},
		{
			title: "success: LONG",
			msg: MsgOpen{
				Position:   Position_LONG,
				Collateral: sdk.NewCoin("uatom", math.NewInt(1)),
			},
			want: nil,
		},
		{
			title: "success: SHORT",
			msg: MsgOpen{
				Position:   Position_SHORT,
				Collateral: sdk.NewCoin("uusdc", math.NewInt(1)),
			},
			want: nil,
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			err := test.msg.ValidatePosition("uatom", "uusdc")

			if test.want != nil {
				require.Error(t, err, test.want)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}

func TestMsgOpen_ValidateTakeProfitAndStopLossPrice(t *testing.T) {

	params := DefaultParams()
	price := math.LegacyNewDec(5)

	type Test struct {
		title string
		msg   MsgOpen
		want  error
	}

	tableTest := []Test{
		{
			title: "success: 0 price",
			msg: MsgOpen{
				TakeProfitPrice: math.LegacyZeroDec(),
				StopLossPrice:   math.LegacyZeroDec(),
			},
			want: nil,
		},
		{
			title: "success: LONG",
			msg: MsgOpen{
				Position:        Position_LONG,
				TakeProfitPrice: price.Mul(params.MaximumLongTakeProfitPriceRatio.Add(params.MinimumLongTakeProfitPriceRatio).QuoInt64(2)),
				StopLossPrice:   price.QuoInt64(2),
			},
			want: nil,
		},
		{
			title: "success: SHORT",
			msg: MsgOpen{
				Position:        Position_SHORT,
				TakeProfitPrice: price.QuoInt64(2),
				StopLossPrice:   price.MulInt64(2),
			},
			want: nil,
		},
		{
			title: "fail: LONG take profit",
			msg: MsgOpen{
				Position:        Position_LONG,
				TakeProfitPrice: price.QuoInt64(2),
			},
			want: fmt.Errorf("take profit price should be between %s and %s times of current market price for long (current ratio: %s)", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String(), price.QuoInt64(2).Quo(price).String()),
		},
		{
			title: "fail: LONG stop loss price",
			msg: MsgOpen{
				Position:        Position_LONG,
				TakeProfitPrice: price.MulInt64(2),
				StopLossPrice:   price.MulInt64(2),
			},
			want: fmt.Errorf("stop loss price cannot be greater than equal to tradingAssetPrice for long (Stop loss: %s, asset price: %s)", price.MulInt64(2).String(), price.String()),
		},
		{
			title: "fail: SHORT take profit",
			msg: MsgOpen{
				Position:        Position_SHORT,
				TakeProfitPrice: price.MulInt64(2),
			},
			want: fmt.Errorf("take profit price should be less than %s times of current market price for short (current ratio: %s)", params.MaximumShortTakeProfitPriceRatio.String(), price.MulInt64(2).Quo(price).String()),
		},
		{
			title: "fail: SHORT stop loss price",
			msg: MsgOpen{
				Position:        Position_SHORT,
				TakeProfitPrice: price.QuoInt64(2),
				StopLossPrice:   price.QuoInt64(2),
			},
			want: fmt.Errorf("stop loss price cannot be less than equal to tradingAssetPrice for short (Stop loss: %s, asset price: %s)", price.QuoInt64(2).String(), price.String()),
		},
	}

	for _, test := range tableTest {
		t.Run(test.title, func(t *testing.T) {
			err := test.msg.ValidateTakeProfitAndStopLossPrice(params, price)

			if test.want != nil {
				require.Error(t, err, test.want)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
