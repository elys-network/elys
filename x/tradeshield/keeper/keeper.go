package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		authority    string
		bank         types.BankKeeper
		amm          types.AmmKeeper
		perpetual    types.PerpetualKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	bank types.BankKeeper,
	amm types.AmmKeeper,
	perpetual types.PerpetualKeeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		bank:         bank,
		amm:          amm,
		perpetual:    perpetual,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAssetPriceFromDenomInToDenomOut returns the price of an asset from a denom to another denom
func (k Keeper) GetAssetPriceFromDenomInToDenomOut(ctx sdk.Context, denomIn, denomOut string) (osmomath.BigDec, error) {
	priceIn := k.amm.CalculateUSDValue(ctx, denomIn, sdkmath.NewInt(1))
	priceOut := k.amm.CalculateUSDValue(ctx, denomOut, sdkmath.NewInt(1))

	// If the price of the asset is 0, return an error
	if priceIn.IsZero() || priceOut.IsZero() {
		return osmomath.ZeroBigDec(), types.ErrPriceNotFound
	}

	// Calculate the price of the asset from denomIn to denomOut
	price := priceIn.Quo(priceOut)

	return price, nil
}
