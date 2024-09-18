package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateParams = "update_params"

var _ sdk.Msg = &MsgUpdateParams{}

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    &params,
	}
}

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate params
	if msg.Params.LimitProcessOrder <= 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "LimitProcessOrder must be greater than 0")
	}
	if msg.Params.RewardPercentage.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "RewardPercentage must be non-negative")
	}
	if msg.Params.MarginError.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "MarginError must be non-negative")
	}
	if msg.Params.MinimumDeposit.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "MinimumDeposit must be non-negative")
	}

	return nil
}
