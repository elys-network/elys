package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgJoinPool{}

func NewMsgJoinPool(sender string, poolId uint64, maxAmountsIn sdk.Coins, shareAmountOut sdkmath.Int) *MsgJoinPool {
	return &MsgJoinPool{
		Sender:         sender,
		PoolId:         poolId,
		MaxAmountsIn:   maxAmountsIn,
		ShareAmountOut: shareAmountOut,
	}
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
