package types

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFeedMultipleExternalLiquidity{}

func NewMsgFeedMultipleExternalLiquidity(sender string) *MsgFeedMultipleExternalLiquidity {
	return &MsgFeedMultipleExternalLiquidity{
		Sender: sender,
	}
}

func (msg *MsgFeedMultipleExternalLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	for _, liquidity := range msg.Liquidity {
		for _, depthInfo := range liquidity.AmountDepthInfo {
			if depthInfo.Asset == "" {
				return errors.New("asset cannot be empty")
			}
			if depthInfo.Depth.IsNil() || depthInfo.Depth.IsNegative() {
				return errors.New("depth cannot be negative or nil")
			}
			if depthInfo.Amount.IsNil() || depthInfo.Amount.IsNegative() {
				return errors.New("depth amount cannot be negative or nil")
			}
		}
	}
	return nil
}
