package types

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePerpetualOpenOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePerpetualOpenOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePerpetualOpenOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePerpetualOpenOrder{
				OwnerAddress:    sample.AccAddress(),
				TriggerPrice:    math.LegacyNewDec(100),
				Collateral:      sdk.NewCoin("token", math.NewInt(1000)),
				TradingAsset:    "asset",
				Position:        PerpetualPosition_LONG,
				Leverage:        math.LegacyNewDec(2),
				TakeProfitPrice: math.LegacyNewDec(150),
				StopLossPrice:   math.LegacyNewDec(90),
				PoolId:          1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgCreatePerpetualCloseOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePerpetualCloseOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePerpetualCloseOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePerpetualCloseOrder{
				OwnerAddress: sample.AccAddress(),
				TriggerPrice: math.LegacyNewDec(100),
				PositionId:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdatePerpetualOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdatePerpetualOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdatePerpetualOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdatePerpetualOrder{
				OwnerAddress: sample.AccAddress(),
				TriggerPrice: math.LegacyNewDec(100),
				OrderId:      1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgCancelPerpetualOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCancelPerpetualOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCancelPerpetualOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCancelPerpetualOrder{
				OwnerAddress: sample.AccAddress(),
				OrderId:      1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgCancelPerpetualOrders_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCancelPerpetualOrders
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCancelPerpetualOrders{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCancelPerpetualOrders{
				OwnerAddress: sample.AccAddress(),
				OrderIds:     []uint64{1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
