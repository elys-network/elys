package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

var EdenBoostApr = sdk.NewDec(1)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		StakeIncentives:         nil,
		EdenCommitVal:           "",
		EdenbCommitVal:          "",
		MaxEdenRewardAprStakers: sdk.NewDecWithPrec(3, 1), // 30%
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
