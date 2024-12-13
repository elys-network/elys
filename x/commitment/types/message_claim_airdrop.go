package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCancelVest{}

func NewMsgClaimAirdrop(creator string) *MsgClaimAirdrop {
	return &MsgClaimAirdrop{
		ClaimAddress: creator,
	}
}

func (msg *MsgClaimAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ClaimAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
