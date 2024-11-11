package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgWhitelist{}

func NewMsgWhitelist(signer string, whitelistedAddress string) *MsgWhitelist {
	return &MsgWhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgWhitelist) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.WhitelistedAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid white listed address (%s)", err)
	}

	return nil
}
