package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDepositTokens = "commit_liquid_tokens"

var _ sdk.Msg = &MsgCommitLiquidTokens{}

func NewMsgCommitLiquidTokens(creator string, amount sdk.Int, denom string) *MsgCommitLiquidTokens {
	return &MsgCommitLiquidTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgCommitLiquidTokens) Route() string {
	return RouterKey
}

func (msg *MsgCommitLiquidTokens) Type() string {
	return TypeMsgDepositTokens
}

func (msg *MsgCommitLiquidTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCommitLiquidTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCommitLiquidTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
