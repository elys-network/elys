package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdatePoolParams = "update_pool_params"

var _ sdk.Msg = &MsgUpdatePoolParams{}

func NewMsgUpdatePoolParams(sender string, poolId uint64, poolParams *PoolParams) *MsgUpdatePoolParams {
	return &MsgUpdatePoolParams{
		Sender:     sender,
		PoolId:     poolId,
		PoolParams: poolParams,
	}
}

func (msg *MsgUpdatePoolParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePoolParams) Type() string {
	return TypeMsgUpdatePoolParams
}

func (msg *MsgUpdatePoolParams) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdatePoolParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePoolParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.PoolParams == nil {
		return ErrPoolParamsShouldNotBeNil
	}

	if msg.PoolParams.SwapFee.IsNegative() {
		return ErrFeeShouldNotBeNegative
	}

	if msg.PoolParams.SwapFee.GT(sdk.NewDecWithPrec(2, 2)) { // >2%
		return ErrSwapFeeShouldNotExceedTwoPercent
	}

	if msg.PoolParams.ExitFee.IsNegative() {
		return ErrFeeShouldNotBeNegative
	}

	if msg.PoolParams.ExitFee.GT(sdk.NewDecWithPrec(2, 2)) { // >2%
		return ErrExitFeeShouldNotExceedTwoPercent
	}

	return nil
}
