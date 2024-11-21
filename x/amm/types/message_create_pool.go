package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePool{}
var _ sdk.HasValidateBasic = &MsgCreatePool{}

func NewMsgCreatePool(sender string, poolParams PoolParams, poolAssets []PoolAsset) *MsgCreatePool {
	return &MsgCreatePool{
		Sender:     sender,
		PoolParams: poolParams,
		PoolAssets: poolAssets,
	}
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(msg.PoolAssets) != 2 {
		return ErrPoolAssetsMustBeTwo
	}

	if err = msg.PoolParams.Validate(); err != nil {
		return err
	}

	for _, asset := range msg.PoolAssets {
		if err = asset.Validate(); err != nil {
			return err
		}
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
