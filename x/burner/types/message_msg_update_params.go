package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateParams = "msg_update_params"

var _ sdk.Msg = &MsgUpdateParams{}

func NewMsgUpdateParams(creator string, authority string, params *Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Params == nil {
		return errorsmod.Wrapf(ErrInvalidParams, "params is nil")
	}

	if len(msg.Params.EpochIdentifier) == 0 {
		return ErrInvalidEpochIdentifier
	}

	return nil
}
