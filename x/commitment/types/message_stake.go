package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStake = "stake"

var _ sdk.Msg = &MsgStake{}

func NewMsgStake(creator string, amount math.Int, asset string, validatorAddress string) *MsgStake {
	return &MsgStake{
		Creator:          creator,
		Amount:           amount,
		Asset:            asset,
		ValidatorAddress: validatorAddress,
	}
}

func (msg *MsgStake) Route() string {
	return RouterKey
}

func (msg *MsgStake) Type() string {
	return TypeMsgStake
}

func (msg *MsgStake) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
