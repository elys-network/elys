package types

import (
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
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
					RewardPercentage:     sdkmath.LegacyMustNewDecFromStr("0.5"),
					MarginError:          sdkmath.LegacyMustNewDecFromStr("0.1"),
					MinimumDeposit:       sdkmath.NewInt(100),
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
					RewardPercentage:     sdkmath.LegacyMustNewDecFromStr("0.5"),
					MarginError:          sdkmath.LegacyMustNewDecFromStr("0.1"),
					MinimumDeposit:       sdkmath.NewInt(100),
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
					RewardPercentage:     sdkmath.LegacyMustNewDecFromStr("-0.5"),
					MarginError:          sdkmath.LegacyMustNewDecFromStr("0.1"),
					MinimumDeposit:       sdkmath.NewInt(100),
				},
			},
			err: errors.New("RewardPercentage is negative"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
