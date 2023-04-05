package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgFeedPrice = "feed_price"
)

var _ sdk.Msg = &MsgFeedPrice{}

func NewMsgFeedPrice(
	creator string,
	asset string,
	price sdk.Dec,
	source string,
) *MsgFeedPrice {
	return &MsgFeedPrice{
		Provider: creator,
		Asset:    asset,
		Price:    price,
		Source:   source,
	}
}

func (msg *MsgFeedPrice) Route() string {
	return RouterKey
}

func (msg *MsgFeedPrice) Type() string {
	return TypeMsgFeedPrice
}

func (msg *MsgFeedPrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFeedPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFeedPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
