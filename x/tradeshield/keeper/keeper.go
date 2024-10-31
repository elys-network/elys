package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		authority    string
		amm          types.AmmKeeper
		tier         types.TierKeeper
		perpetual    types.PerpetualKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	amm types.AmmKeeper,
	tier types.TierKeeper,
	perpetual types.PerpetualKeeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		amm:          amm,
		tier:         tier,
		perpetual:    perpetual,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetUserDiscount returns the discount of a user
func (k Keeper) GetUserDiscount(ctx sdk.Context, address string) (sdkmath.LegacyDec, error) {
	// Get user address
	user, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// Get discount tier (range from 0 to 100)
	_, _, discount := k.tier.GetMembershipTier(ctx, user)

	// Convert uint64 discount to math.LegacyDec
	discountDec := sdkmath.LegacyNewDec(int64(discount))

	// Normalize the discount to be between 0 and 1
	discountDec = discountDec.Quo(sdkmath.LegacyNewDec(100))

	// Return the discount
	return discountDec, nil
}

// GetAssetPriceFromDenomInToDenomOut returns the price of an asset from a denom to another denom
func (k Keeper) GetAssetPriceFromDenomInToDenomOut(ctx sdk.Context, denomIn, denomOut string) (sdkmath.LegacyDec, error) {
	priceIn := k.tier.CalculateUSDValue(ctx, denomIn, sdkmath.NewInt(1))
	priceOut := k.tier.CalculateUSDValue(ctx, denomOut, sdkmath.NewInt(1))

	// If the price of the asset is 0, return an error
	if priceIn.IsZero() || priceOut.IsZero() {
		return sdkmath.LegacyZeroDec(), types.ErrPriceNotFound
	}

	// Calculate the price of the asset from denomIn to denomOut
	price := priceIn.Quo(priceOut)

	return price, nil
}
