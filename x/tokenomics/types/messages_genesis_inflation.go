package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateGenesisInflation = "update_genesis_inflation"
)

var _ sdk.Msg = &MsgUpdateGenesisInflation{}

func NewMsgUpdateGenesisInflation(authority string, inflation InflationEntry, seedVesting uint64, strategicSalesVesting uint64) *MsgUpdateGenesisInflation {
	return &MsgUpdateGenesisInflation{
		Authority:             authority,
		Inflation:             &inflation,
		SeedVesting:           seedVesting,
		StrategicSalesVesting: strategicSalesVesting,
	}
}

func (msg *MsgUpdateGenesisInflation) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGenesisInflation) Type() string {
	return TypeMsgUpdateGenesisInflation
}

func (msg *MsgUpdateGenesisInflation) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateGenesisInflation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGenesisInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
