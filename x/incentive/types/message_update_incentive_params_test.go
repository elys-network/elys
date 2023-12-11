package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateIncentiveParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateIncentiveParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateIncentiveParams{
				Authority:                   "invalid_address",
				RewardPortionForLps:         sdk.NewDecWithPrec(60, 2),
				RewardPortionForStakers:     sdk.NewDecWithPrec(30, 2),
				MaxEdenRewardAprStakers:     sdk.NewDecWithPrec(3, 1),
				MaxEdenRewardAprLps:         sdk.NewDecWithPrec(3, 1),
				DistributionEpochForStakers: 10,
				DistributionEpochForLps:     10,
				ElysStakeTrackingRate:       10,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateIncentiveParams{
				Authority:                   sample.AccAddress(),
				RewardPortionForLps:         sdk.NewDecWithPrec(60, 2),
				RewardPortionForStakers:     sdk.NewDecWithPrec(30, 2),
				MaxEdenRewardAprStakers:     sdk.NewDecWithPrec(3, 1),
				MaxEdenRewardAprLps:         sdk.NewDecWithPrec(3, 1),
				DistributionEpochForStakers: 10,
				DistributionEpochForLps:     10,
				ElysStakeTrackingRate:       10,
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
