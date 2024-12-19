package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgClaimKol(creator string, refund bool) *MsgClaimKol {
	return &MsgClaimKol{
		ClaimAddress: creator,
		Refund:       refund,
	}
}

func (msg *MsgClaimKol) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ClaimAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
