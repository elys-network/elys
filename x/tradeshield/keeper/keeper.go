package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/utils"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
		authority          string
		bank               types.BankKeeper
		amm                types.AmmKeeper
		perpetual          types.PerpetualKeeper
		assetProfileKeeper types.AssetProfileKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	bank types.BankKeeper,
	amm types.AmmKeeper,
	perpetual types.PerpetualKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeService:       storeService,
		authority:          authority,
		bank:               bank,
		amm:                amm,
		perpetual:          perpetual,
		assetProfileKeeper: assetProfileKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAssetPriceFromDenomInToDenomOut returns the price of an asset from a denom to another denom
func (k Keeper) GetAssetPriceFromDenomInToDenomOut(ctx sdk.Context, denomIn, denomOut string) (osmomath.BigDec, error) {
	// Get asset profile entries for both tokens to normalize decimals
	assetIn, found := k.assetProfileKeeper.GetEntryByDenom(ctx, denomIn)
	if !found {
		return osmomath.ZeroBigDec(), types.ErrPriceNotFound
	}

	assetOut, found := k.assetProfileKeeper.GetEntryByDenom(ctx, denomOut)
	if !found {
		return osmomath.ZeroBigDec(), types.ErrPriceNotFound
	}

	// Calculate USD value of 1 full token (10^decimals base units) for both tokens
	// This normalizes the decimal differences between tokens
	priceIn := k.amm.CalculateUSDValue(ctx, denomIn, sdkmath.NewInt(utils.Pow10Int64(assetIn.Decimals)))
	priceOut := k.amm.CalculateUSDValue(ctx, denomOut, sdkmath.NewInt(utils.Pow10Int64(assetOut.Decimals)))

	// If the price of the asset is 0, return an error
	if priceIn.IsZero() || priceOut.IsZero() {
		return osmomath.ZeroBigDec(), types.ErrPriceNotFound
	}

	// Calculate the price of the asset from denomIn to denomOut
	price := priceIn.Quo(priceOut)

	return price, nil
}
