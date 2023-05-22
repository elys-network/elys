package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/elys-network/elys/x/amm/types"
)

// This function:
// - saves the pool to state
// - Mints LP shares to the pool creator
// - Sets bank metadata for the LP denom
// - Records total liquidity increase
// - Calls the AfterPoolCreated hook
func (k Keeper) InitializePool(ctx sdk.Context, pool *types.Pool, sender sdk.AccAddress) (err error) {
	// Mint the initial pool shares share token to the sender
	err = k.MintPoolShareToAccount(ctx, pool, sender, pool.GetTotalShares().Amount)
	if err != nil {
		return err
	}

	// Finally, add the share token's meta data to the bank keeper.
	poolShareBaseDenom := types.GetPoolShareDenom(pool.GetPoolId())
	poolShareDisplayDenom := fmt.Sprintf("AMM-%d", pool.GetPoolId())
	k.bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: fmt.Sprintf("The share token of the amm pool %d", pool.GetPoolId()),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    poolShareBaseDenom,
				Exponent: 0,
				Aliases: []string{
					"attopoolshare",
				},
			},
			{
				Denom:    poolShareDisplayDenom,
				Exponent: types.OneShareExponent,
				Aliases:  nil,
			},
		},
		Base:    poolShareBaseDenom,
		Display: poolShareDisplayDenom,
	})

	if err := k.SetPool(ctx, *pool); err != nil {
		return err
	}

	k.hooks.AfterPoolCreated(ctx, sender, pool.GetPoolId())
	// k.RecordTotalLiquidityIncrease(ctx, pool.GetTotalPoolLiquidity(ctx))
	return nil
}
