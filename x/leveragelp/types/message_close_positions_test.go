package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"

	"github.com/elys-network/elys/v7/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClosePositions(t *testing.T) {
	creator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	liquidatePositions := []*types.PositionRequest{}
	stoplossPositions := []*types.PositionRequest{}
	msg := types.NewMsgClosePositions(creator, liquidatePositions, stoplossPositions)
	tests := []struct {
		name   string
		setter func()
		errMsg string
	}{
		{
			name: "invalid address",
			setter: func() {
				msg.Creator = "invalid_address"
			},
			errMsg: "invalid creator address",
		},
		{
			name: "liquidate and stoploss position requests both are empty",
			setter: func() {
				msg.Creator = sample.AccAddress()
			},
			errMsg: "no liquidate or stoploss position requests",
		},
		{
			name: "invalid liquidate position requests",
			setter: func() {
				msg.Liquidate = append(msg.Liquidate, &types.PositionRequest{
					Address: "address",
					Id:      1,
				})
			},
			errMsg: "invalid liquidation address",
		},
		{
			name: "repeated liquidate positions",
			setter: func() {
				msg.Liquidate = liquidatePositions
				msg.Liquidate = append(msg.Liquidate,
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					},
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					})
			},
			errMsg: "repeated liquidation id",
		},
		{
			name: "invalid stoploss position requests",
			setter: func() {
				msg.Liquidate = liquidatePositions
				msg.StopLoss = append(msg.StopLoss, &types.PositionRequest{
					Address: "address",
					Id:      1,
				})
			},
			errMsg: "invalid stoploss address",
		},
		{
			name: "repeated stoploss positions",
			setter: func() {
				msg.StopLoss = stoplossPositions
				msg.StopLoss = append(msg.StopLoss,
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					},
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					})
			},
			errMsg: "repeated stoploss id",
		},
		{
			name: "repeated stoploss and liquidate positions",
			setter: func() {
				msg.Liquidate = liquidatePositions
				msg.Liquidate = append(msg.Liquidate,
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					})
				msg.StopLoss = stoplossPositions
				msg.StopLoss = append(msg.StopLoss,
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					})
			},
			errMsg: "repeated stoploss id",
		},
		{
			name: "success",
			setter: func() {
				msg.Creator = creator.String()
				msg.Liquidate = liquidatePositions
				msg.Liquidate = append(msg.Liquidate,
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      1,
					})
				msg.StopLoss = stoplossPositions
				msg.StopLoss = append(msg.StopLoss,
					&types.PositionRequest{
						Address: sample.AccAddress(),
						Id:      2,
					})
			},
			errMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := msg.ValidateBasic()
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
