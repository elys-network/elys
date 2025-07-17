package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(
	lpIncentives *IncentiveInfo,
	rewardPortionForLps sdkmath.LegacyDec,
	rewardPortionForStakers sdkmath.LegacyDec,
	maxEdenRewardAprLps sdkmath.LegacyDec,
	protocolRevenueAddress string,
	takerManager string,
) Params {
	return Params{
		LpIncentives:            lpIncentives,
		RewardPortionForLps:     rewardPortionForLps,
		RewardPortionForStakers: rewardPortionForStakers,
		MaxEdenRewardAprLps:     maxEdenRewardAprLps,
		SupportedRewardDenoms:   nil,
		ProtocolRevenueAddress:  protocolRevenueAddress,
		TakerManager:            takerManager,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		nil,
		sdkmath.LegacyNewDecWithPrec(60, 2),
		sdkmath.LegacyNewDecWithPrec(25, 2),
		sdkmath.LegacyNewDecWithPrec(5, 1),
		authtypes.NewModuleAddress("protocol-revenue-address").String(), // TODO: Change it in genesis in mainnet launch
		authtypes.NewModuleAddress("taker-manager-address").String(),
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.RewardPortionForLps.IsNil() {
		return fmt.Errorf("reward percent for lp must not be nil")
	}
	if p.RewardPortionForLps.IsNegative() {
		return fmt.Errorf("reward percent for lp must be positive: %s", p.RewardPortionForLps.String())
	}
	if p.RewardPortionForLps.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("reward percent for lp too large: %s", p.RewardPortionForLps.String())
	}

	if p.RewardPortionForStakers.IsNil() {
		return fmt.Errorf("reward percent for stakers must not be nil")
	}
	if p.RewardPortionForStakers.IsNegative() {
		return fmt.Errorf("reward percent for stakers must be positive: %s", p.RewardPortionForStakers.String())
	}
	if p.RewardPortionForStakers.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("reward percent for stakers too large: %s", p.RewardPortionForStakers.String())
	}

	if p.RewardPortionForStakers.Add(p.RewardPortionForLps).GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("reward percent for stakers + lp too large: %s", p.RewardPortionForStakers.Add(p.RewardPortionForLps).String())
	}

	if p.MaxEdenRewardAprLps.IsNil() {
		return fmt.Errorf("MaxEdenRewardAprLps must not be nil")
	}
	if p.MaxEdenRewardAprLps.IsNegative() {
		return fmt.Errorf("MaxEdenRewardAprLps must be positive: %s", p.MaxEdenRewardAprLps.String())
	}

	if p.LpIncentives != nil && p.LpIncentives.EdenAmountPerYear.LTE(sdkmath.ZeroInt()) {
		return fmt.Errorf("invalid eden amount per year: %v", p.LpIncentives.String())
	}

	if p.LpIncentives != nil && p.LpIncentives.BlocksDistributed < 0 {
		return fmt.Errorf("invalid BlocksDistributed: %v", p.LpIncentives.String())
	}

	_, err := sdk.AccAddressFromBech32(p.ProtocolRevenueAddress)
	if err != nil {
		return fmt.Errorf("invalid ProtocolRevenueAddress %s: %s", p.ProtocolRevenueAddress, err.Error())
	}

	for _, supportedRewardDenom := range p.SupportedRewardDenoms {
		if err = sdk.ValidateDenom(supportedRewardDenom.Denom); err != nil {
			return fmt.Errorf("invalid reward denom %s: %v", supportedRewardDenom.Denom, err)
		}
		if supportedRewardDenom.MinAmount.IsNil() {
			return fmt.Errorf("reward denom minimum amount cannot be nil: %s", supportedRewardDenom.Denom)
		}
		if supportedRewardDenom.MinAmount.IsNegative() {
			return fmt.Errorf("reward denom(%s) minimum amount cannot be negative: %s", supportedRewardDenom.Denom, supportedRewardDenom.MinAmount.String())
		}
	}

	_, err = sdk.AccAddressFromBech32(p.TakerManager)
	if err != nil {
		return fmt.Errorf("invalid TakerManager %s: %s", p.TakerManager, err.Error())
	}
	return nil
}

// String implements the Stringer interface.
func (p LegacyParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p Params) GetBigDecMaxEdenRewardAprLps() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaxEdenRewardAprLps)
}
