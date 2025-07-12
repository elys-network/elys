package types

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewBalancerPool returns a weighted CPMM pool with the provided parameters, and initial assets.
// Invariants that are assumed to be satisfied and not checked:
// (This is handled in ValidateBasic)
// * 2 <= len(assets) <= 8
// * FutureGovernor is valid
// * poolID doesn't already exist
func NewBalancerPool(poolId uint64, balancerPoolParams PoolParams, assets []PoolAsset, blockTime time.Time) (Pool, error) {
	poolAddr := NewPoolAddress(poolId)
	poolRebalanceTreasuryAddr := NewPoolRebalanceTreasury(poolId)

	// pool that's created up to ensuring the assets and params are valid.
	// We assume that FuturePoolGovernor is valid.
	pool := &Pool{
		PoolId:            poolId,
		Address:           poolAddr.String(),
		RebalanceTreasury: poolRebalanceTreasuryAddr.String(),
		PoolParams:        PoolParams{},
		TotalWeight:       math.ZeroInt(),
		TotalShares:       sdk.NewCoin(GetPoolShareDenom(poolId), InitPoolSharesSupply),
		PoolAssets:        nil,
	}

	err := pool.SetInitialPoolAssets(assets)
	if err != nil {
		return Pool{}, err
	}

	sortedPoolAssets := pool.GetAllPoolAssets()
	err = balancerPoolParams.Validate()
	if err != nil {
		return Pool{}, err
	}

	err = pool.setInitialPoolParams(balancerPoolParams, sortedPoolAssets, blockTime)
	if err != nil {
		return Pool{}, err
	}

	return *pool, nil
}
