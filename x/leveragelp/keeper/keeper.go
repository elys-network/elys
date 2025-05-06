package keeper

import (
	"fmt"
	gomath "math"

	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeService        store.KVStoreService
		authority           string
		amm                 types.AmmKeeper
		bankKeeper          types.BankKeeper
		oracleKeeper        ammtypes.OracleKeeper
		stableKeeper        types.StableStakeKeeper
		commKeeper          types.CommitmentKeeper
		masterchefKeeper    types.MasterchefKeeper
		accountedPoolKeeper types.AccountedPoolKeeper

		hooks types.LeverageLpHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	amm types.AmmKeeper,
	bk types.BankKeeper,
	oracleKeeper ammtypes.OracleKeeper,
	stableKeeper types.StableStakeKeeper,
	commitmentKeeper types.CommitmentKeeper,
	masterchefKeeper types.MasterchefKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	keeper := &Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		authority:           authority,
		amm:                 amm,
		bankKeeper:          bk,
		oracleKeeper:        oracleKeeper,
		stableKeeper:        stableKeeper,
		commKeeper:          commitmentKeeper,
		masterchefKeeper:    masterchefKeeper,
		accountedPoolKeeper: accountedPoolKeeper,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetWhitelistKey(address), address)
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetWhitelistKey(address))
}

func (k Keeper) CheckIfWhitelisted(ctx sdk.Context, address sdk.AccAddress) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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
	snapshot := k.amm.GetPoolWithAccountedBalance(ctx, ammPool.PoolId)
	swapResult, _, err := k.amm.CalcInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, osmomath.ZeroBigDec())
	if err != nil {
		return math.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return math.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) {
	pool.Health = k.CalculatePoolHealth(ctx, pool).Dec()
	k.SetPool(ctx, *pool)
}

func (k Keeper) CalculatePoolHealth(ctx sdk.Context, pool *types.Pool) osmomath.BigDec {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return osmomath.ZeroBigDec()
	}

	if ammPool.TotalShares.Amount.IsZero() {
		return osmomath.OneBigDec()
	}

	return osmomath.BigDecFromSDKInt(ammPool.TotalShares.Amount.Sub(pool.LeveragedLpAmount)).
		Quo(osmomath.BigDecFromSDKInt(ammPool.TotalShares.Amount))
}

func (k Keeper) GetWhitelistAddressIterator(ctx sdk.Context) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return storetypes.KVStorePrefixIterator(store, types.WhitelistPrefix)
}

func (k Keeper) GetAllWhitelistedAddress(ctx sdk.Context) []sdk.AccAddress {
	var list []sdk.AccAddress
	iterator := k.GetWhitelistAddressIterator(ctx)
	defer func(iterator storetypes.Iterator) {
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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
