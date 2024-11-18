package types_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/x/estaking/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	params := types.Params{
		StakeIncentives:         nil,
		EdenCommitVal:           "",
		EdenbCommitVal:          "",
		MaxEdenRewardAprStakers: sdkmath.LegacyNewDecWithPrec(3, 1), // 30%
		EdenBoostApr:            sdkmath.LegacyOneDec(),
		DexRewardsStakers: types.DexRewardsTracker{
			NumBlocks: 1,
			Amount:    sdkmath.LegacyZeroDec(),
		},
	}
	require.Equal(t, types.DefaultParams(), params)
}

func TestRewardPortionForLps(t *testing.T) {
	params := types.DefaultParams()
	params.EdenCommitVal = "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30"
	params.EdenbCommitVal = "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s"
	tests := []struct {
		name   string
		setter func()
		err    string
	}{
		{
			name: "success",
			setter: func() {
			},
			err: "",
		},
		{
			name: "invalid EdenCommitVal address",
			setter: func() {
				params.EdenCommitVal = "invalid address"
			},
			err: "invalid EdenCommitVal address",
		},
		{
			name: "invalid EdenBCommitVal address",
			setter: func() {
				params.EdenCommitVal = "cosmosvaloper1x8efhljzvs52u5xa6m7crcwes7v9u0nlwdgw30"
				params.EdenbCommitVal = "cosmos18ruzecmqj9pv8ac0gvkgryuc7u004te9xr2mcr"
			},
			err: "invalid EdenBCommitVal address",
		},
		{
			name: "stake incentive BlocksDistributed is -ve",
			setter: func() {
				params.EdenbCommitVal = "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s"
				params.StakeIncentives = &types.IncentiveInfo{
					BlocksDistributed: -1,
					EdenAmountPerYear: sdkmath.ZeroInt(),
				}
			},
			err: "StakeIncentives blocks distributed must be >= 0",
		},
		{
			name: "invalid eden amount per year",
			setter: func() {
				params.StakeIncentives.BlocksDistributed = 10
			},
			err: "invalid eden amount per year",
		},
		{
			name: "MaxEdenRewardAprStakers is nil",
			setter: func() {
				params.StakeIncentives.EdenAmountPerYear = sdkmath.OneInt()
				params.MaxEdenRewardAprStakers = sdkmath.LegacyDec{}
			},
			err: "MaxEdenRewardAprStakers must not be nil",
		},
		{
			name: "MaxEdenRewardAprStakers is -ve",
			setter: func() {
				params.MaxEdenRewardAprStakers = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			err: "MaxEdenRewardAprStakers cannot be negative",
		},
		{
			name: "EdenBoostApr is nil",
			setter: func() {
				params.MaxEdenRewardAprStakers = sdkmath.LegacyOneDec()
				params.EdenBoostApr = sdkmath.LegacyDec{}
			},
			err: "EdenBoostApr must not be nil",
		},
		{
			name: "EdenBoostApr is -ve",
			setter: func() {
				params.EdenBoostApr = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			err: "EdenBoostApr cannot be negative",
		},
		{
			name: "DexRewardsStakers DexRewardsStakers is nil",
			setter: func() {
				params.EdenBoostApr = sdkmath.LegacyOneDec()
				params.DexRewardsStakers.Amount = sdkmath.LegacyDec{}
			},
			err: "DexRewardsStakers amount must not be nil",
		},
		{
			name: "DexRewardsStakers DexRewardsStakers is -ve",
			setter: func() {
				params.DexRewardsStakers.Amount = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			err: "DexRewardsStakers amount cannot be -ve",
		},
		{
			name: "DexRewardsStakers NumBlocks is -ve",
			setter: func() {
				params.DexRewardsStakers.Amount = sdkmath.LegacyOneDec()
				params.DexRewardsStakers.NumBlocks = -1
			},
			err: "DexRewardsStakers NumBlocks cannot be -ve",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setter()
			err := params.Validate()
			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
