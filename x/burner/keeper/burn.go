package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/elys-network/elys/x/burner/types"
)

// ShouldBurnTokens checks if tokens should be burned for the given epoch
func (k Keeper) ShouldBurnTokens(ctx sdk.Context, epochIdentifier string) bool {
	params := k.GetParams(ctx)
	return epochIdentifier == params.EpochIdentifier
}

// BurnTokensForAllDenoms burns tokens for all denominations
func (k Keeper) BurnTokensForAllDenoms(ctx sdk.Context) error {
	balances := k.getPositiveBalances(ctx)
	for denom, balance := range balances {
		if err := k.burnTokensForDenom(ctx, balance, denom); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) getPositiveBalances(ctx sdk.Context) map[string]sdk.Coins {
	zeroAddress := types.GetZeroAddress()
	balances := make(map[string]sdk.Coins)
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		// Get the balance of the zero address for this denom
		balance := k.bankKeeper.GetBalance(ctx, zeroAddress, metadata.Base)
		if balance.IsPositive() {
			balances[metadata.Base] = sdk.NewCoins(balance)
		}
		return false
	})
	return balances
}

func (k Keeper) burnTokensForDenom(ctx sdk.Context, balance sdk.Coins, denom string) error {
	if err := k.sendCoinsFromZeroAddressToModule(ctx, balance); err != nil {
		k.Logger(ctx).Error("Error sending coins and burning tokens", "denom", denom, "error", err)
		return err
	}
	if err := k.burnCoins(ctx, balance); err != nil {
		k.Logger(ctx).Error("Error burning tokens", "denom", denom, "error", err)
		return err
	}
	k.Logger(ctx).Info("Burned tokens for denom", denom)

	// check if balance has at least one coin
	if len(balance) == 0 {
		return nil
	}

	// Record a history item
	history := types.History{
		Timestamp: ctx.BlockTime().String(),
		Denom:     denom,
		Amount:    balance[0].Amount.String(),
	}
	k.SetHistory(ctx, history)

	return nil
}

func (k Keeper) sendCoinsFromZeroAddressToModule(ctx sdk.Context, coins sdk.Coins) error {
	zeroAddress := types.GetZeroAddress()
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, zeroAddress, types.ModuleName, coins)
}

func (k Keeper) burnCoins(ctx sdk.Context, coins sdk.Coins) error {
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}
