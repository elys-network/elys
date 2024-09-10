package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
			NumBlocks: 1,
			Amount:    sdk.ZeroDec(),
		},
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if _, err := sdk.ValAddressFromBech32(p.EdenCommitVal); p.EdenCommitVal != "" && err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EdenCommitVal address (%s)", err)
	}
	if _, err := sdk.ValAddressFromBech32(p.EdenbCommitVal); p.EdenbCommitVal != "" && err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EdenBCommitVal address (%s)", err)
	}
	if p.StakeIncentives != nil {
		if p.StakeIncentives.BlocksDistributed < 0 {
			return fmt.Errorf("StakeIncentives blocks distributed must be >= 0")
		}
		if p.StakeIncentives.EdenAmountPerYear.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("invalid eden amount per year: %s", p.StakeIncentives.EdenAmountPerYear.String())
		}
	}
	if p.MaxEdenRewardAprStakers.IsNil() {
		return fmt.Errorf("MaxEdenRewardAprStakers must not be nil")
	}
	if p.MaxEdenRewardAprStakers.IsNegative() {
		return fmt.Errorf("MaxEdenRewardAprStakers cannot be negative: %s", p.MaxEdenRewardAprStakers.String())
	}
	if p.EdenBoostApr.IsNil() {
		return fmt.Errorf("EdenBoostApr must not be nil")
	}
	if p.EdenBoostApr.IsNegative() {
		return fmt.Errorf("EdenBoostApr cannot be negative: %s", p.EdenBoostApr.String())
	}
	if p.DexRewardsStakers.Amount.IsNil() {
		return fmt.Errorf("DexRewardsStakers amount must not be nil")
	}
	if p.DexRewardsStakers.Amount.IsNegative() {
		return fmt.Errorf("DexRewardsStakers amount cannot be -ve")
	}
	if p.DexRewardsStakers.NumBlocks < 0 {
		return fmt.Errorf("DexRewardsStakers NumBlocks cannot be -ve")
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p LegacyParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
