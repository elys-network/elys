package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/tradeshield/types"
)

type (
	Keeper struct {
		cdc       codec.BinaryCodec
		storeKey  storetypes.StoreKey
		memKey    storetypes.StoreKey
		authority string
		amm       types.AmmKeeper
		tier      types.TierKeeper
		perpetual types.PerpetualKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	authority string,
	amm types.AmmKeeper,
	tier types.TierKeeper,
	perpetual types.PerpetualKeeper,
) *Keeper {
	return &Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		memKey:    memKey,
		authority: authority,
		amm:       amm,
		tier:      tier,
		perpetual: perpetual,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetUserDiscount returns the discount of a user
func (k Keeper) GetUserDiscount(ctx sdk.Context, address string) (sdk.Dec, error) {
	// Get user address
	user, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// Get discount tier (range from 0 to 100)
	_, _, discount := k.tier.GetMembershipTier(ctx, user)

	// Convert uint64 discount to sdk.Dec
	discountDec := sdk.NewDec(int64(discount))

	// Normalize the discount to be between 0 and 1
	discountDec = discountDec.Quo(sdk.NewDec(100))

	// Return the discount
	return discountDec, nil
}

// GetAssetPriceFromDenomInToDenomOut returns the price of an asset from a denom to another denom
func (k Keeper) GetAssetPriceFromDenomInToDenomOut(ctx sdk.Context, denomIn, denomOut string) (sdk.Dec, error) {
	priceIn := k.tier.CalculateUSDValue(ctx, denomIn, sdk.NewInt(1))
	priceOut := k.tier.CalculateUSDValue(ctx, denomOut, sdk.NewInt(1))

	// If the price of the asset is 0, return an error
	if priceIn.IsZero() || priceOut.IsZero() {
		return sdk.ZeroDec(), types.ErrPriceNotFound
	}

	// Calculate the price of the asset from denomIn to denomOut
	price := priceIn.Quo(priceOut)

	return price, nil
}
