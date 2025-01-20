package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgExitPool{}

func NewMsgExitPool(sender string, poolId uint64, minAmountsOut sdk.Coins, shareAmountIn sdkmath.Int) *MsgExitPool {
	return &MsgExitPool{
		Sender:        sender,
		PoolId:        poolId,
		MinAmountsOut: minAmountsOut,
		ShareAmountIn: shareAmountIn,
	}
}

func (msg *MsgExitPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	for _, coin := range msg.MinAmountsOut {
		if err = coin.Validate(); err != nil {
			return err
		}
	}

	if msg.TokenOutDenom == "" {
		return errors.New("exit in single token not allowed")
	}

	if msg.ShareAmountIn.IsNil() {
		return ErrInvalidShareAmountOut
	}

	if msg.ShareAmountIn.IsNegative() {
		return ErrInvalidShareAmountOut
	}

	return nil
}
