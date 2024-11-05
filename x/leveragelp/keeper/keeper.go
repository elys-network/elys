package keeper

import (
	"fmt"
	gomath "math"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		memKey              storetypes.StoreKey
		authority           string
		amm                 types.AmmKeeper
		bankKeeper          types.BankKeeper
		oracleKeeper        ammtypes.OracleKeeper
		stableKeeper        types.StableStakeKeeper
		commKeeper          types.CommitmentKeeper
		assetProfileKeeper  types.AssetProfileKeeper
		masterchefKeeper    types.MasterchefKeeper
		accountedPoolKeeper types.AccountedPoolKeeper

		hooks types.LeverageLpHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	authority string,
	amm types.AmmKeeper,
	bk types.BankKeeper,
	oracleKeeper ammtypes.OracleKeeper,
	stableKeeper types.StableStakeKeeper,
	commitmentKeeper types.CommitmentKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	masterchefKeeper types.MasterchefKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	keeper := &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		authority:           authority,
		amm:                 amm,
		bankKeeper:          bk,
		oracleKeeper:        oracleKeeper,
		stableKeeper:        stableKeeper,
		commKeeper:          commitmentKeeper,
		assetProfileKeeper:  assetProfileKeeper,
		masterchefKeeper:    masterchefKeeper,
		accountedPoolKeeper: accountedPoolKeeper,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWhitelistKey(address), address)
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWhitelistKey(address))
}

func (k Keeper) CheckIfWhitelisted(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWhitelistKey(address))
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (math.Int, error) {
	_, found := k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return math.Int{}, fmt.Errorf("pool %d not found", ammPool.PoolId)
	}

	tokensOut := sdk.NewCoins(tokenOutAmount)
	// Estimate swap
	snapshot := k.amm.GetAccountedPoolSnapshotOrSet(ctx, ammPool)
	swapResult, _, err := k.amm.CalcInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, sdk.ZeroDec())
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return sdk.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) {
	pool.Health = k.CalculatePoolHealth(ctx, pool)
	k.SetPool(ctx, *pool)
}

func (k Keeper) CalculatePoolHealth(ctx sdk.Context, pool *types.Pool) sdk.Dec {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return sdk.ZeroDec()
	}

	if ammPool.TotalShares.Amount.IsZero() {
		return sdk.OneDec()
	}

	return sdk.NewDecFromBigInt(ammPool.TotalShares.Amount.Sub(pool.LeveragedLpAmount).BigInt()).
		Quo(sdk.NewDecFromBigInt(ammPool.TotalShares.Amount.BigInt()))
}

func (k Keeper) GetWhitelistAddressIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WhitelistPrefix)
}

func (k Keeper) GetAllWhitelistedAddress(ctx sdk.Context) []sdk.AccAddress {
	var list []sdk.AccAddress
	iterator := k.GetWhitelistAddressIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Value())
	}

	return list
}

func (k Keeper) GetWhitelistedAddress(ctx sdk.Context, pagination *query.PageRequest) ([]sdk.AccAddress, *query.PageResponse, error) {
	var list []sdk.AccAddress
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.WhitelistPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: gomath.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(prefixStore, pagination, func(key []byte, value []byte) error {
		list = append(list, value)
		return nil
	})

	return list, pageRes, err
}

// SetHooks set the leveragelp hooks
func (k *Keeper) SetHooks(lh types.LeverageLpHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set leveragelp hooks twice")
	}

	k.hooks = lh

	return k
}

func (k Keeper) GetHooks() types.LeverageLpHooks {
	return k.hooks
}
