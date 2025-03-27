package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFeedPrice{}

func NewMsgFeedPrice(
	creator string,
	asset string,
	price sdkmath.LegacyDec,
	source string,
) *MsgFeedPrice {
	return &MsgFeedPrice{
		Provider: creator,
		FeedPrice: FeedPrice{
			Asset:  asset,
			Price:  price,
			Source: source,
		},
	}
}

func (price FeedPrice) Validate() error {
	if price.Price.IsNil() {
		return errorsmod.Wrapf(ErrInvalidPrice, "price is nil")
	}

	if price.Price.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidPrice, "price is negative")
	}

	if len(price.Source) == 0 {
		return errorsmod.Wrapf(ErrInvalidPrice, "source is empty")
	}

	return nil
}

func (msg *MsgFeedPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address (%s)", err)
	}

	if err = msg.FeedPrice.Validate(); err != nil {
		return errorsmod.Wrapf(ErrInvalidPrice, "invalid price (%s)", err)
	}

	return nil
}
