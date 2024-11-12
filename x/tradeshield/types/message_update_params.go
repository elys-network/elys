package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateParams{}

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    &params,
	}
}
func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate params
	if msg.Params.LimitProcessOrder <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "LimitProcessOrder must be greater than 0")
	}
	if msg.Params.RewardPercentage.IsNegative() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "RewardPercentage must be non-negative")
	}
	if msg.Params.MarginError.IsNegative() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "MarginError must be non-negative")
	}
	if msg.Params.MinimumDeposit.IsNegative() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "MinimumDeposit must be non-negative")
	}

	return nil
}
