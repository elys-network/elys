package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdatePoolMultipliers = "update_pool_multipliers"

var _ sdk.Msg = &MsgUpdatePoolMultipliers{}

func NewMsgUpdatePoolMultipliers(creator string, multipliers []PoolMultiplier) *MsgUpdatePoolMultipliers {
	return &MsgUpdatePoolMultipliers{
		Authority:       creator,
		PoolMultipliers: multipliers,
	}
}

func (msg *MsgUpdatePoolMultipliers) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePoolMultipliers) Type() string {
	return TypeMsgUpdatePoolMultipliers
}

func (msg *MsgUpdatePoolMultipliers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePoolMultipliers) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePoolMultipliers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.PoolMultipliers) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid number of pool multipliers (%d)", len(msg.PoolMultipliers))
	}
	return nil
}
