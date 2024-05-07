package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/utils"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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
	sender := msg.GetSigners()[0]

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return 0, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// If the fee denom is empty, set it to the base currency
	if msg.PoolParams.FeeDenom == "" {
		msg.PoolParams.FeeDenom = baseCurrency
	}

	// Get the next pool ID and increment the pool ID counter
	// Create the pool with the given pool ID
	poolId := k.GetNextPoolId(ctx)
	pool, err := types.NewBalancerPool(poolId, *msg.PoolParams, msg.PoolAssets, ctx.BlockTime())
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
	if err := k.InitializePool(ctx, &pool, sender); err != nil {
		return 0, err
	}

	// Send initial liquidity to the pool's address.
	initialPoolLiquidity := msg.InitialLiquidity()
	err = k.bankKeeper.SendCoins(ctx, sender, address, initialPoolLiquidity)
	if err != nil {
		return 0, err
	}

	// Increase liquidty amount
	for _, asset := range msg.PoolAssets {
		k.RecordTotalLiquidityIncrease(ctx, sdk.Coins{asset.Token})
	}

	// emitCreatePoolEvents(ctx, poolId, msg)
	return pool.GetPoolId(), nil
}

// This function:
// - saves the pool to state
// - Mints LP shares to the pool creator
// - Sets bank metadata for the LP denom
// - Records total liquidity increase
// - Calls the AfterPoolCreated hook
func (k Keeper) InitializePool(ctx sdk.Context, pool *types.Pool, sender sdk.AccAddress) (err error) {
	tvl, err := pool.TVL(ctx, k.oracleKeeper)
	if err != nil {
		return err
	}

	if tvl.IsPositive() {
		pool.TotalShares = sdk.NewCoin(pool.TotalShares.Denom, tvl.Mul(sdk.NewDecFromInt(types.OneShare)).RoundInt())
	}

	// Mint the initial pool shares share token to the sender
	err = k.MintPoolShareToAccount(ctx, *pool, sender, pool.GetTotalShares().Amount)
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

	if k.hooks != nil {
		k.hooks.AfterPoolCreated(ctx, sender, *pool)
	}
	return nil
}
