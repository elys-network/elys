package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddCollateral = "add_collateral"

var _ sdk.Msg = &MsgAddCollateral{}

func NewMsgAddCollateral(creator string, id uint64, collateral sdk.Int) *MsgAddCollateral {
	return &MsgAddCollateral{
		Creator:    creator,
		Id:         id,
		Collateral: collateral,
	}
}

func (msg *MsgAddCollateral) Route() string {
	return RouterKey
}

func (msg *MsgAddCollateral) Type() string {
	return TypeMsgAddCollateral
}

func (msg *MsgAddCollateral) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !msg.Collateral.IsPositive() {
		return fmt.Errorf("collateral should be positive (%s)", msg.Collateral.String())
	}
	return nil
}
