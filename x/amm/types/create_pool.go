package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgCreatePool) CreatePool(ctx sdk.Context, poolID uint64) (*Pool, error) {
	pool, err := NewBalancerPool(poolID, *msg.PoolParams, msg.PoolAssets, ctx.BlockTime())
	return &pool, err
}
