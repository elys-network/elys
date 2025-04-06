package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	epochsmoduletypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		StakeIncentives:                nil,
		EdenCommitVal:                  "",
		EdenbCommitVal:                 "",
		MaxEdenRewardAprStakers:        sdkmath.LegacyNewDecWithPrec(3, 1), // 30%
		EdenBoostApr:                   sdkmath.LegacyOneDec(),
		ProviderVestingEpochIdentifier: epochsmoduletypes.TenDaysEpochID,
		ProviderStakingRewardsPortion:  sdkmath.LegacyMustNewDecFromStr("0.25"),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if _, err := sdk.ValAddressFromBech32(p.EdenCommitVal); p.EdenCommitVal != "" && err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EdenCommitVal address (%s)", err)
	}
	if _, err := sdk.ValAddressFromBech32(p.EdenbCommitVal); p.EdenbCommitVal != "" && err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EdenBCommitVal address (%s)", err)
	}
	if p.StakeIncentives != nil {
		if p.StakeIncentives.BlocksDistributed < 0 {
			return fmt.Errorf("StakeIncentives blocks distributed must be >= 0")
		}
		if p.StakeIncentives.EdenAmountPerYear.LTE(sdkmath.ZeroInt()) {
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
	if p.ProviderVestingEpochIdentifier == "" {
		return fmt.Errorf("ProviderVestingEpochIdentifier must not be empty")
	}
	if p.ProviderStakingRewardsPortion.IsNil() {
		return fmt.Errorf("ProviderStakingRewardsPortion must not be nil")
	}
	if p.ProviderStakingRewardsPortion.IsNegative() {
		return fmt.Errorf("ProviderStakingRewardsPortion cannot be negative: %s", p.ProviderStakingRewardsPortion.String())
	}
	return nil
}

func (p Params) GetBigDecMaxEdenRewardAprStakers() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaxEdenRewardAprStakers)
}

func (p Params) GetBigDecEdenBoostApr() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.EdenBoostApr)
}
