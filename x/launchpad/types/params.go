package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		InitialPrice:       sdk.NewDecWithPrec(75, 2),  // 0.75
		TotalReserve:       sdk.NewInt(4500000_000000), // 4.5 million
		SoldAmount:         sdk.NewInt(0),
		WithdrawAddress:    "",
		WithdrawnAmount:    sdk.NewCoins(),
		LaunchpadStarttime: 1710984623,  // 2024-03-21
		LaunchpadDuration:  86400 * 100, // 100 days
		ReturnDuration:     86400 * 180, // 6 months
		MaxReturnPercent:   50,          // 50%
		SpendingTokens:     []string{"uusdc", "uatom", "wei"},
		BonusInfo: Bonus{
			//   0-20% raise 100% bonus
			//   20-30% raise 90% bonus
			//   30-40% raise bonus 80%
			//   40-50% raise bonus 70%
			//   50-60% raise bonus 60%
			//   60-70% raise bonus 50%
			//   70-80% raise bonus 40%
			//   80-100% raise bonus 30%
			RaisePercents:   []uint64{20, 30, 40, 50, 60, 70, 80, 100},
			BonusPercents:   []uint64{100, 90, 80, 70, 60, 50, 40, 30},
			LockDuration:    86400 * 180, // 6 months
			VestingDuration: 86400 * 365, // 1 year
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
