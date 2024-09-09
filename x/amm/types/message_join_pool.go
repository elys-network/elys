package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgJoinPool = "join_pool"

var _ sdk.Msg = &MsgJoinPool{}

func NewMsgJoinPool(sender string, poolId uint64, maxAmountsIn sdk.Coins, shareAmountOut math.Int) *MsgJoinPool {
	return &MsgJoinPool{
		Sender:         sender,
		PoolId:         poolId,
		MaxAmountsIn:   maxAmountsIn,
		ShareAmountOut: shareAmountOut,
	}
}

func (msg *MsgJoinPool) Route() string {
	return RouterKey
}

func (msg *MsgJoinPool) Type() string {
	return TypeMsgJoinPool
}

func (msg *MsgJoinPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgJoinPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgJoinPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.ShareAmountOut.IsNil() {
		return ErrInvalidShareAmountOut
	}

	if msg.ShareAmountOut.IsNegative() {
		return ErrInvalidShareAmountOut
	}

	return nil
}
