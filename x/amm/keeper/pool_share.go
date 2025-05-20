package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/amm/types"
	atypes "github.com/elys-network/elys/v4/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

// MintPoolShareToAccount attempts to mint shares of a AMM denomination to the
// specified address returning an error upon failure. Shares are minted using
// the x/amm module account.
func (k Keeper) MintPoolShareToAccount(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress, amount math.Int) error {
	poolShareDenom := types.GetPoolShareDenom(pool.GetPoolId())
	amt := sdk.NewCoins(sdk.NewCoin(poolShareDenom, amount))

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amt)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, amt)
	if err != nil {
		return err
	}

	// All LP tokens minted should be committed to commitment module in order to make
	// the liquidity provider gets rewarded.
	// So deposit, and then commit
	// Before commit LP token to commitment module, we should first register the new denom
	// to assetProfile module.

	_, found := k.assetProfileKeeper.GetEntry(ctx, poolShareDenom)
	if !found {
		// Set an entity to assetprofile
		entry := atypes.Entry{
			Authority:                pool.Address,
			BaseDenom:                poolShareDenom,
			Decimals:                 ptypes.BASE_DECIMAL,
			Denom:                    poolShareDenom,
			Path:                     "",
			IbcChannelId:             "",
			IbcCounterpartyChannelId: "",
			DisplayName:              poolShareDenom,
			DisplaySymbol:            "",
			Network:                  "",
			Address:                  "",
			ExternalSymbol:           "",
			TransferLimit:            "",
			Permissions:              make([]string, 0),
			UnitDenom:                "",
			IbcCounterpartyDenom:     "",
			IbcCounterpartyChainId:   "",
			CommitEnabled:            true,
			WithdrawEnabled:          true,
		}

		k.assetProfileKeeper.SetEntry(ctx, entry)
	}

	// Commit LP token minted
	lockUntil := uint64(0)
	if pool.PoolParams.UseOracle {
		params := k.GetParams(ctx)
		lockUntil = uint64(ctx.BlockTime().Unix()) + params.LpLockupDuration
	}

	err = k.commitmentKeeper.CommitLiquidTokens(ctx, addr, poolShareDenom, amount, lockUntil)
	if err != nil {
		return err
	}

	return nil
}

// BurnPoolShareFromAccount burns `amount` of the given pool's shares held by `addr`.
func (k Keeper) BurnPoolShareFromAccount(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress, amount math.Int) error {
	coin := sdk.NewCoin(types.GetPoolShareDenom(pool.GetPoolId()), amount)
	coins := sdk.NewCoins(coin)

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	return nil
}
