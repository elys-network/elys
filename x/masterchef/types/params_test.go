package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	require.Equal(t, types.DefaultParams(), types.NewParams(nil, sdk.NewDecWithPrec(60, 2), sdk.NewDecWithPrec(25, 2), sdk.NewDecWithPrec(5, 1), authtypes.NewModuleAddress(govtypes.ModuleName).String()))
	output, err := yaml.Marshal(types.DefaultParams())
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams().String(), string(output))
}

func TestRewardPortionForLps(t *testing.T) {
	params := types.DefaultParams()
	params.ProtocolRevenueAddress = "cosmos1vjclnpz4hydg0nv5xn2xtfvg52dlnslnndyh0a"
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
			name: "RewardPortionForLps is nil",
			setter: func() {
				params.RewardPortionForLps = sdk.Dec{}
			},
			err: "reward percent for lp must not be nil",
		},
		{
			name: "RewardPortionForLps < 0",
			setter: func() {
				params.RewardPortionForLps = sdk.MustNewDecFromStr("-0.5")
			},
			err: "reward percent for lp must be positive",
		}, {
			name: "RewardPortionForLps > 1",
			setter: func() {
				params.RewardPortionForLps = sdk.OneDec().MulInt64(100)
			},
			err: "reward percent for lp too large:",
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

func TestRewardPortionForStakers(t *testing.T) {
	params := types.DefaultParams()
	params.ProtocolRevenueAddress = "cosmos1vjclnpz4hydg0nv5xn2xtfvg52dlnslnndyh0a"
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
			name: "RewardPortionForStakers is nil",
			setter: func() {
				params.RewardPortionForStakers = sdk.Dec{}
			},
			err: "reward percent for stakers must not be nil",
		},
		{
			name: "RewardPortionForStakers < 0",
			setter: func() {
				params.RewardPortionForStakers = sdk.MustNewDecFromStr("-0.5")
			},
			err: "reward percent for stakers must be positive",
		}, {
			name: "RewardPortionForLps > 1",
			setter: func() {
				params.RewardPortionForStakers = sdk.OneDec().MulInt64(100)
			},
			err: "reward percent for stakers too large:",
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

func TestLPIncentives(t *testing.T) {
	params := types.DefaultParams()
	params.ProtocolRevenueAddress = "cosmos1vjclnpz4hydg0nv5xn2xtfvg52dlnslnndyh0a"
	incentiveInfo := types.IncentiveInfo{
		EdenAmountPerYear: sdk.OneInt().MulRaw(1000),
		BlocksDistributed: 10,
	}
	params.LpIncentives = &incentiveInfo
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
			name: "invalid eden amount per year",
			setter: func() {
				params.LpIncentives.EdenAmountPerYear = sdk.OneInt().MulRaw(-1000)
			},
			err: "invalid eden amount per year",
		},
		{
			name: "invalid BlocksDistributed",
			setter: func() {
				params.LpIncentives.EdenAmountPerYear = sdk.OneInt().MulRaw(1000)
				params.LpIncentives.BlocksDistributed = -10
			},
			err: "invalid BlocksDistributed",
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

func TestMaxEdenRewardAprLps(t *testing.T) {
	params := types.DefaultParams()
	params.ProtocolRevenueAddress = "cosmos1vjclnpz4hydg0nv5xn2xtfvg52dlnslnndyh0a"
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
			name: "MaxEdenRewardAprLps is nil",
			setter: func() {
				params.MaxEdenRewardAprLps = sdk.Dec{}
			},
			err: "MaxEdenRewardAprLps must not be nil",
		},
		{
			name: "MaxEdenRewardAprLps is -ve",
			setter: func() {
				params.MaxEdenRewardAprLps = sdk.OneDec().MulInt64(-1)
			},
			err: "MaxEdenRewardAprLps must be positive",
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

func TestProtocolRevenueAddress(t *testing.T) {
	params := types.DefaultParams()
	params.ProtocolRevenueAddress = "cosmos1vjclnpz4hydg0nv5xn2xtfvg52dlnslnndyh0a"
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
			name: "ProtocolRevenueAddres is empty",
			setter: func() {
				params.ProtocolRevenueAddress = ""
			},
			err: "ProtocolRevenueAddres cannot be empty",
		},
		{
			name: "invalid ProtocolRevenueAddress",
			setter: func() {
				params.ProtocolRevenueAddress = "abcd"
			},
			err: "invalid ProtocolRevenueAddress",
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

func TestSupportedRewardDenoms(t *testing.T) {
	params := types.DefaultParams()
	params.ProtocolRevenueAddress = "cosmos1vjclnpz4hydg0nv5xn2xtfvg52dlnslnndyh0a"
	supportedRewardDenoms := []*types.SupportedRewardDenom{
		{
			"uusdc",
			sdk.OneInt(),
		},
	}
	params.SupportedRewardDenoms = supportedRewardDenoms
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
			name: "invalid reward denom",
			setter: func() {
				params.SupportedRewardDenoms[0].Denom = "%%%%"
			},
			err: "invalid reward denom",
		},
		{
			name: "reward denom minimum amount is nil",
			setter: func() {
				params.SupportedRewardDenoms[0].Denom = "uusdc"
				params.SupportedRewardDenoms[0].MinAmount = sdk.Int{}
			},
			err: "reward denom minimum amount cannot be nil",
		},
		{
			name: "reward denom minimum amount is -v",
			setter: func() {
				params.SupportedRewardDenoms[0].Denom = "uusdc"
				params.SupportedRewardDenoms[0].MinAmount = sdk.NewInt(-1000)
			},
			err: "minimum amount cannot be negative",
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
