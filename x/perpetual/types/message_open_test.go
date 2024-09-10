package types

import (
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
		sdk.NewDec(200),
		"uatom",
		sdk.NewCoin("uusdc", sdk.NewInt(2000)),
		sdk.NewDec(100),
		sdk.NewDec(0),
	)

	want := &MsgOpen{
		Creator:         accAdress,
		Position:        Position_LONG,
		Leverage:        sdk.NewDec(200),
		TradingAsset:    "uatom",
		Collateral:      sdk.NewCoin("uusdc", sdk.NewInt(2000)),
		TakeProfitPrice: sdk.NewDec(100),
		StopLossPrice:   sdk.NewDec(0),
	}

	assert.Equal(t, want, got)
}

func TestMsgOpen_Route(t *testing.T) {
	msg := MsgOpen{}
	assert.Equal(t, "perpetual", msg.Route())
}

func TestMsgOpen_Type(t *testing.T) {
	msg := MsgOpen{}
	assert.Equal(t, "open", msg.Type())
}

func TestMsgOpen_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgOpen{Creator: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgOpen_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgOpen{Creator: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
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
				Leverage:     sdk.NewDec(-200),
			},
			want: ErrInvalidLeverage,
		},
		{
			title: "trading asset is empty",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_LONG,
				TradingAsset: "",
				Leverage:     sdk.NewDec(200),
			},
			want: ErrTradingAssetIsEmpty,
		},
		{
			title: "take profit price is nil",
			msg: MsgOpen{
				Creator:      sample.AccAddress(),
				Position:     Position_SHORT,
				TradingAsset: "uatom",
				Leverage:     sdk.NewDec(200),
			},
			want: ErrInvalidTakeProfitPriceIsNegative,
		},
		{
			title: "take profit price is negative",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_SHORT,
				TradingAsset:    "uatom",
				TakeProfitPrice: sdk.NewDec(-10),
				Leverage:        sdk.NewDec(200),
			},
			want: ErrInvalidTakeProfitPriceIsNegative,
		},
		{
			title: "successful",
			msg: MsgOpen{
				Creator:         sample.AccAddress(),
				Position:        Position_LONG,
				TradingAsset:    "uatom",
				TakeProfitPrice: sdk.NewDec(10),
				Leverage:        sdk.NewDec(200),
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
