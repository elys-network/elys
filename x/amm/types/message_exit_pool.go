package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgExitPool = "exit_pool"

var _ sdk.Msg = &MsgExitPool{}

func NewMsgExitPool(sender string, poolId uint64, minAmountsOut sdk.Coins, shareAmountIn sdk.Int) *MsgExitPool {
	return &MsgExitPool{
		Sender:        sender,
		PoolId:        poolId,
		MinAmountsOut: minAmountsOut,
		ShareAmountIn: shareAmountIn,
	}
}

func (msg *MsgExitPool) Route() string {
	return RouterKey
}

func (msg *MsgExitPool) Type() string {
	return TypeMsgExitPool
}

func (msg *MsgExitPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgExitPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgExitPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
