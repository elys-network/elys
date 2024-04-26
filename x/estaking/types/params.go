package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		StakeIncentives:         nil,
		EdenCommitVal:           "",
		EdenbCommitVal:          "",
		MaxEdenRewardAprStakers: sdk.NewDecWithPrec(3, 1), // 30%
		EdenBoostApr:            sdk.OneDec(),
		DexRewardsStakers: DexRewardsTracker{
			NumBlocks: sdk.OneInt(),
			Amount:    sdk.ZeroDec(),
		},
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
