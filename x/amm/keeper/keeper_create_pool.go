package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/utils"
)

// CreatePool attempts to create a pool returning the newly created pool ID or
// an error upon failure. The pool creation fee is used to fund the community
// pool. It will create a dedicated module account for the pool and sends the
// initial liquidity to the created module account.
//
// After the initial liquidity is sent to the pool's account, this function calls an
// InitializePool function from the source module. That module is responsible for:
// - saving the pool into its own state
// - Minting LP shares to pool creator
// - Setting metadata for the shares
func (k Keeper) CreatePool(ctx sdk.Context, msg *types.MsgCreatePool) (uint64, error) {
	// Send pool creation fee to community pool
	// params := k.GetParams(ctx)
	sender := msg.GetSigners()[0]
	// if err := k.communityPoolKeeper.FundCommunityPool(ctx, params.PoolCreationFee, sender); err != nil {
	// 	return 0, err
	// }

	// Get the next pool ID and increment the pool ID counter
	// Create the pool with the given pool ID
	poolId := k.GetNextPoolId(ctx)
	pool, err := msg.CreatePool(ctx, poolId)
	if err != nil {
		return 0, err
	}

	if err := pool.Validate(poolId); err != nil {
		return 0, err
	}

	address, err := sdk.AccAddressFromBech32(pool.GetAddress())
	if err != nil {
		return 0, fmt.Errorf("invalid pool address %s", pool.GetAddress())
	}

	// create and save the pool's module account to the account keeper
	if err := utils.CreateModuleAccount(ctx, k.accountKeeper, address); err != nil {
		return 0, fmt.Errorf("creating pool module account for id %d: %w", poolId, err)
	}

	// Run the initialization logic.
	if err := k.InitializePool(ctx, pool, sender); err != nil {
		return 0, err
	}

	// Send initial liquidity to the pool's address.
	initialPoolLiquidity := msg.InitialLiquidity()
	err = k.bankKeeper.SendCoins(ctx, sender, address, initialPoolLiquidity)
	if err != nil {
		return 0, err
	}

	// emitCreatePoolEvents(ctx, poolId, msg)
	return pool.GetPoolId(), nil
}
