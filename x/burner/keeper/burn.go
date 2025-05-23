package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/elys-network/elys/v5/x/burner/types"
)

// ShouldBurnTokens checks if tokens should be burned for the given epoch
func (k Keeper) ShouldBurnTokens(ctx sdk.Context, epochIdentifier string) bool {
	params := k.GetParams(ctx)
	return epochIdentifier == params.EpochIdentifier
}

// BurnTokensForAllDenoms burns tokens for all denominations
func (k Keeper) BurnTokensForAllDenoms(ctx sdk.Context) error {

	zeroAddress := types.GetZeroAddress()
	coinsToBurn := sdk.NewCoins()
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		// Get the balance of the zero address for this denom
		balance := k.bankKeeper.GetBalance(ctx, zeroAddress, metadata.Base)
		if balance.IsPositive() {
			coinsToBurn = coinsToBurn.Add(balance)
		}
		return false
	})

	if !coinsToBurn.IsZero() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, zeroAddress, types.ModuleName, coinsToBurn); err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Error sending coins for burning tokens %s, err: %s", coinsToBurn.String(), err.Error()))
			return err
		}
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToBurn); err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Error burning tokens %s, err: %s", coinsToBurn.String(), err.Error()))
			return err
		}
		k.Logger(ctx).Info("Burned tokens %s", coinsToBurn.String())

		// Record a history item
		history := types.History{
			Block:       uint64(ctx.BlockHeight()),
			BurnedCoins: coinsToBurn,
		}
		k.SetHistory(ctx, history)
	}

	return nil
}
