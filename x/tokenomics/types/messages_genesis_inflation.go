package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateGenesisInflation{}

func NewMsgUpdateGenesisInflation(authority string, inflation InflationEntry, seedVesting uint64, strategicSalesVesting uint64) *MsgUpdateGenesisInflation {
	return &MsgUpdateGenesisInflation{
		Authority:             authority,
		Inflation:             &inflation,
		SeedVesting:           seedVesting,
		StrategicSalesVesting: strategicSalesVesting,
	}
}

func (msg *MsgUpdateGenesisInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate Inflation is not nil and its fields are positive
	if msg.Inflation == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "inflation entry cannot be nil")
	}
	if msg.Inflation.LmRewards <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "lm rewards must be positive")
	}
	if msg.Inflation.IcsStakingRewards <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "ics staking rewards must be positive")
	}
	if msg.Inflation.CommunityFund <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "community fund must be positive")
	}
	if msg.Inflation.StrategicReserve <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "strategic reserve must be positive")
	}
	if msg.Inflation.TeamTokensVested <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "team tokens vested must be positive")
	}

	// Validate SeedVesting is positive
	if msg.SeedVesting <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "seed vesting must be positive")
	}

	// Validate StrategicSalesVesting is positive
	if msg.StrategicSalesVesting <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "strategic sales vesting must be positive")
	}

	return nil
}
