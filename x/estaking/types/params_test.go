package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	params := types.Params{
		StakeIncentives:         nil,
		EdenCommitVal:           "",
		EdenbCommitVal:          "",
		MaxEdenRewardAprStakers: sdk.NewDecWithPrec(3, 1), // 30%
		EdenBoostApr:            sdk.OneDec(),
		DexRewardsStakers: types.DexRewardsTracker{
			NumBlocks: 1,
			Amount:    sdk.ZeroDec(),
		},
	}
	require.Equal(t, types.DefaultParams(), params)
	output, err := yaml.Marshal(types.DefaultParams())
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams().String(), string(output))
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
					EdenAmountPerYear: sdk.ZeroInt(),
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
				params.StakeIncentives.EdenAmountPerYear = sdk.OneInt()
				params.MaxEdenRewardAprStakers = sdk.Dec{}
			},
			err: "MaxEdenRewardAprStakers must not be nil",
		},
		{
			name: "MaxEdenRewardAprStakers is -ve",
			setter: func() {
				params.MaxEdenRewardAprStakers = sdk.OneDec().MulInt64(-1)
			},
			err: "MaxEdenRewardAprStakers cannot be negative",
		},
		{
			name: "EdenBoostApr is nil",
			setter: func() {
				params.MaxEdenRewardAprStakers = sdk.OneDec()
				params.EdenBoostApr = sdk.Dec{}
			},
			err: "EdenBoostApr must not be nil",
		},
		{
			name: "EdenBoostApr is -ve",
			setter: func() {
				params.EdenBoostApr = sdk.OneDec().MulInt64(-1)
			},
			err: "EdenBoostApr cannot be negative",
		},
		{
			name: "DexRewardsStakers DexRewardsStakers is nil",
			setter: func() {
				params.EdenBoostApr = sdk.OneDec()
				params.DexRewardsStakers.Amount = sdk.Dec{}
			},
			err: "DexRewardsStakers amount must not be nil",
		},
		{
			name: "DexRewardsStakers DexRewardsStakers is -ve",
			setter: func() {
				params.DexRewardsStakers.Amount = sdk.OneDec().MulInt64(-1)
			},
			err: "DexRewardsStakers amount cannot be -ve",
		},
		{
			name: "DexRewardsStakers NumBlocks is -ve",
			setter: func() {
				params.DexRewardsStakers.Amount = sdk.OneDec()
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
