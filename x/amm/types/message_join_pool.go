package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgJoinPool = "join_pool"

var _ sdk.Msg = &MsgJoinPool{}

func NewMsgJoinPool(creator string, poolId uint64, maxAmountsIn sdk.Coins, shareAmountOut string) *MsgJoinPool {
	return &MsgJoinPool{
		Creator:        creator,
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
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgJoinPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgJoinPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
