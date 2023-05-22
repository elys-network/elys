package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgCreatePool) CreatePool(ctx sdk.Context, poolID uint64) (*Pool, error) {
	poolParams := PoolParams{
		SwapFee: msg.SwapFee,
		ExitFee: msg.ExitFee,
	}
	poolAssets := make([]PoolAsset, len(msg.InitialDeposit))
	pool, err := NewBalancerPool(poolID, poolParams, poolAssets, ctx.BlockTime())
	return &pool, err
}
