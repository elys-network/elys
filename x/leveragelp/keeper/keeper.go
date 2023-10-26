package keeper

import (
	"fmt"

	"math"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		memKey       storetypes.StoreKey
		authority    string
		amm          types.AmmKeeper
		bankKeeper   types.BankKeeper
		oracleKeeper ammtypes.OracleKeeper
		stableKeeper types.StableStakeKeeper
		commKeeper   types.CommitmentKeeper

		hooks types.LeveragelpHooks
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
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	keeper := &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		memKey:       memKey,
		authority:    authority,
		amm:          amm,
		bankKeeper:   bk,
		oracleKeeper: oracleKeeper,
		stableKeeper: stableKeeper,
		commKeeper:   commitmentKeeper,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWhitelistKey(address), []byte(address))
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWhitelistKey(address))
}

func (k Keeper) CheckIfWhitelisted(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWhitelistKey(address))
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (sdk.Int, error) {
	leveragelpEnabled := k.IsPoolEnabled(ctx, ammPool.PoolId)
	if !leveragelpEnabled {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrLeveragelpDisabled, "Leveragelp disabled pool")
	}

	tokensOut := sdk.Coins{tokenOutAmount}
	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	swapResult, err := k.amm.CalcInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, sdk.ZeroDec())

	if err != nil {
		return sdk.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return sdk.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) error {
	pool.Health = k.CalculatePoolHealth(ctx, pool)
	k.SetPool(ctx, *pool)

	return nil
}

func (k Keeper) CalculatePoolHealth(ctx sdk.Context, pool *types.Pool) sdk.Dec {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return sdk.ZeroDec()
	}

	if ammPool.TotalShares.Amount.IsZero() {
		return sdk.OneDec()
	}

	return sdk.NewDecFromBigInt(pool.LeveragedLpAmount.BigInt()).Quo(sdk.NewDecFromBigInt(ammPool.TotalShares.Amount.BigInt()))
}

func (k Keeper) TakeFundPayment(ctx sdk.Context, returnAmount sdk.Int, returnAsset string, takePercentage sdk.Dec, fundAddr sdk.AccAddress, ammPool *ammtypes.Pool) (sdk.Int, error) {
	returnAmountDec := sdk.NewDecFromBigInt(returnAmount.BigInt())
	takeAmount := sdk.NewIntFromBigInt(takePercentage.Mul(returnAmountDec).TruncateInt().BigInt())

	if !takeAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, sdk.NewIntFromBigInt(takeAmount.BigInt())))
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ammPool.Address, fundAddr, takeCoins)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	}
	return takeAmount, nil
}

func (k Keeper) GetWhitelistAddressIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WhitelistPrefix)
}

func (k Keeper) GetAllWhitelistedAddress(ctx sdk.Context) []string {
	var list []string
	iterator := k.GetWhitelistAddressIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, (string)(iterator.Value()))
	}

	return list
}

func (k Keeper) GetWhitelistedAddress(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error) {
	var list []string
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.WhitelistPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(prefixStore, pagination, func(key []byte, value []byte) error {
		list = append(list, string(value))
		return nil
	})

	return list, pageRes, err
}

// Set the leveragelp hooks.
func (k *Keeper) SetHooks(gh types.LeveragelpHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set leveragelp hooks twice")
	}

	k.hooks = gh

	return k
}
