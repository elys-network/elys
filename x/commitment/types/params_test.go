package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	require.Equal(t, types.DefaultParams(), types.Params{
		VestingInfos:        nil,
		TotalCommitted:      nil,
		NumberOfCommitments: 0,
	})
	output, err := yaml.Marshal(types.DefaultParams())
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams().String(), string(output))
}

func TestParamsValidation(t *testing.T) {
	vestingInfos := []*types.VestingInfo{
		{
			"uusdc",
			"uatom",
			12,
			sdk.OneInt(),
			123,
		},
	}
	params := types.DefaultParams()
	params.VestingInfos = vestingInfos
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
			name: "VestingInfo has invalid base denom",
			setter: func() {
				params.VestingInfos[0].BaseDenom = "@@@@@"
			},
			err: "invalid denom",
		},
		{
			name: "VestingInfo has invalid VestingDenom",
			setter: func() {
				params.VestingInfos[0].BaseDenom = "uusdc"
				params.VestingInfos[0].VestingDenom = "@@@@"
			},
			err: "invalid denom",
		},
		{
			name: "num_max_vestings is negative",
			setter: func() {
				params.VestingInfos[0].VestingDenom = "uusdc"
				params.VestingInfos[0].NumMaxVestings = -1
			},
			err: "num_max_vestings cannot be negative",
		},
		{
			name: "num_blocks is negative",
			setter: func() {
				params.VestingInfos[0].NumMaxVestings = 10
				params.VestingInfos[0].NumBlocks = -1
			},
			err: "num_blocks cannot be negative",
		},
		{
			name: "VestNowFactor is nil",
			setter: func() {
				params.VestingInfos[0].NumBlocks = 10
				params.VestingInfos[0].VestNowFactor = sdk.Int{}
			},
			err: "vesting now factor cannot be nil",
		},
		{
			name: "VestNowFactor is < 1",
			setter: func() {
				params.VestingInfos[0].VestNowFactor = sdk.ZeroInt()
			},
			err: "vesting now factor must be positive",
		},
		{
			name: "params.VestingInfos contains nil",
			setter: func() {
				params.VestingInfos[0].VestNowFactor = sdk.OneInt()
				params.VestingInfos = append(params.VestingInfos, nil)
			},
			err: "vesting info cannot be nil",
		},
		{
			name: "TotalCommitted is invalid",
			setter: func() {
				params.VestingInfos = vestingInfos
				params.TotalCommitted = sdk.Coins{sdk.Coin{"@@@@", sdk.OneInt()}}
			},
			err: "invalid denom",
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
