package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) GetTotalPoolLiquidity() sdk.Coins {
	return poolAssetsCoins(p.PoolAssets)
}
