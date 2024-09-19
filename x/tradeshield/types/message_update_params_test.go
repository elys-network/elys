package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
				Params: &Params{
					MarketOrderEnabled:   true,
					StakeEnabled:         true,
					ProcessOrdersEnabled: true,
					SwapEnabled:          true,
					PerpetualEnabled:     true,
					RewardEnabled:        true,
					LeverageEnabled:      true,
					LimitProcessOrder:    1,
					RewardPercentage:     sdk.MustNewDecFromStr("0.5"),
					MarginError:          sdk.MustNewDecFromStr("0.1"),
					MinimumDeposit:       sdk.NewInt(100),
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					MarketOrderEnabled:   true,
					StakeEnabled:         true,
					ProcessOrdersEnabled: true,
					SwapEnabled:          true,
					PerpetualEnabled:     true,
					RewardEnabled:        true,
					LeverageEnabled:      true,
					LimitProcessOrder:    1,
					RewardPercentage:     sdk.MustNewDecFromStr("0.5"),
					MarginError:          sdk.MustNewDecFromStr("0.1"),
					MinimumDeposit:       sdk.NewInt(100),
				},
			},
		},
		{
			name: "negative RewardPercentage",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					MarketOrderEnabled:   true,
					StakeEnabled:         true,
					ProcessOrdersEnabled: true,
					SwapEnabled:          true,
					PerpetualEnabled:     true,
					RewardEnabled:        true,
					LeverageEnabled:      true,
					LimitProcessOrder:    1,
					RewardPercentage:     sdk.MustNewDecFromStr("-0.5"),
					MarginError:          sdk.MustNewDecFromStr("0.1"),
					MinimumDeposit:       sdk.NewInt(100),
				},
			},
			err: sdkerrors.ErrInvalidRequest,
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
