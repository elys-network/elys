package keeper

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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

	entry, found := k.assetProfileKeeper.GetEntry(ctx, poolShareDenom)
	if !found {
		// Set an entity to assetprofile
		entry = atypes.Entry{
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
	msgServer := commitmentkeeper.NewMsgServerImpl(*k.commitmentKeeper)

	lockUntil := uint64(ctx.BlockTime().Unix())
	if pool.PoolParams.UseOracle {
		lockUntil += uint64(time.Hour.Seconds())
	}

	// Create a commit LP token message
	msgLiquidCommitLPToken := &ctypes.MsgCommitLiquidTokens{
		Creator:   addr.String(),
		Denom:     poolShareDenom,
		Amount:    amount,
		LockUntil: lockUntil,
	}

	// Commit LP token
	_, err = msgServer.CommitLiquidTokens(sdk.WrapSDKContext(ctx), msgLiquidCommitLPToken)
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
