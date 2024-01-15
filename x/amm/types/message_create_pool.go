package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(sender string, poolParams *PoolParams, poolAssets []PoolAsset) *MsgCreatePool {
	return &MsgCreatePool{
		Sender:     sender,
		PoolParams: poolParams,
		PoolAssets: poolAssets,
	}
}

func (msg *MsgCreatePool) Route() string {
	return RouterKey
}

func (msg *MsgCreatePool) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
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

func (msg *MsgCreatePool) InitialLiquidity() sdk.Coins {
	var coins sdk.Coins
	for _, asset := range msg.PoolAssets {
		coins = append(coins, asset.Token)
	}
	if coins == nil {
		panic("InitialLiquidity coins is equal to nil - this shouldn't happen")
	}
	coins = coins.Sort()
	return coins
}
