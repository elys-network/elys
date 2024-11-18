package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFeedMultiplePrices{}

func (msg *MsgFeedMultiplePrices) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Prices) == 0 {
		return fmt.Errorf("no prices provided")
	}

	for _, price := range msg.Prices {
		if price.Price.IsNil() {
			return errorsmod.Wrapf(ErrInvalidPrice, "price is nil")
		}

		if price.Price.IsNegative() {
			return errorsmod.Wrapf(ErrInvalidPrice, "price is negative")
		}
	}
	return nil
}
