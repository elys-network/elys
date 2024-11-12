package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTimeBasedInflation{}

func NewMsgCreateTimeBasedInflation(
	authority string,
	startBlockHeight uint64,
	endBlockHeight uint64,
	description string,
	inflation *InflationEntry,
) *MsgCreateTimeBasedInflation {
	return &MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
		Description:      description,
		Inflation:        inflation,
	}
}

func (msg *MsgCreateTimeBasedInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate EndBlockHeight is positive and after StartBlockHeight
	if msg.EndBlockHeight <= msg.StartBlockHeight {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "end block height must be after start block height")
	}

	// Validate Description is not empty
	if len(msg.Description) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "description cannot be empty")
	}

	if err := validateInflationEntry(msg.Inflation); err != nil {
		return err
	}

	return nil
}

var _ sdk.Msg = &MsgUpdateTimeBasedInflation{}

func NewMsgUpdateTimeBasedInflation(
	authority string,
	startBlockHeight uint64,
	endBlockHeight uint64,
	description string,
	inflation *InflationEntry,
) *MsgUpdateTimeBasedInflation {
	return &MsgUpdateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
		Description:      description,
		Inflation:        inflation,
	}
}

func (msg *MsgUpdateTimeBasedInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate StartBlockHeight is positive
	if msg.StartBlockHeight <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "start block height must be positive")
	}

	// Validate EndBlockHeight is positive and after StartBlockHeight
	if msg.EndBlockHeight <= msg.StartBlockHeight {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "end block height must be after start block height")
	}

	// Validate Description is not empty
	if len(msg.Description) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "description cannot be empty")
	}

	// Validate Inflation is not nil and its fields are positive
	if msg.Inflation == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "inflation entry cannot be nil")
	}
	if err := validateInflationEntry(msg.Inflation); err != nil {
		return err
	}

	return nil
}

var _ sdk.Msg = &MsgDeleteTimeBasedInflation{}

func NewMsgDeleteTimeBasedInflation(
	authority string,
	startBlockHeight uint64,
	endBlockHeight uint64,
) *MsgDeleteTimeBasedInflation {
	return &MsgDeleteTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
	}
}

func (msg *MsgDeleteTimeBasedInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate StartBlockHeight is positive
	if msg.StartBlockHeight <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "start block height must be positive")
	}

	// Validate EndBlockHeight is positive and after StartBlockHeight
	if msg.EndBlockHeight <= msg.StartBlockHeight {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "end block height must be after start block height")
	}

	return nil
}

func validateInflationEntry(inflation *InflationEntry) error {
	// Validate Inflation is not nil and its fields are positive
	if inflation == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "inflation entry cannot be nil")
	}

	if inflation.LmRewards <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "lm rewards must be positive")
	}

	if inflation.IcsStakingRewards <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "ics staking rewards must be positive")
	}

	if inflation.CommunityFund <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "community fund must be positive")
	}

	if inflation.StrategicReserve <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "strategic reserve must be positive")
	}

	if inflation.TeamTokensVested <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "team tokens vested must be positive")
	}

	return nil
}
