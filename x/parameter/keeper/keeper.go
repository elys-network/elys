package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetAnteHandlerParam set a specific anteHandlerParam in the store from its index
func (k Keeper) SetAnteHandlerParam(ctx sdk.Context, anteHandlerParam types.AnteHandlerParam) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AnteHandlerParamKeyPrefix))
	b := k.cdc.MustMarshal(&anteHandlerParam)
	store.Set(types.AnteHandlerParamKey(
		types.AnteStoreKey,
	), b)
}

// GetAnteHandlerParam returns a anteHandlerParam from its index
func (k Keeper) GetAnteHandlerParam(
	ctx sdk.Context,

) (val types.AnteHandlerParam, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AnteHandlerParamKeyPrefix))

	b := store.Get(types.AnteHandlerParamKey(
		types.AnteStoreKey,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
