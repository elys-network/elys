package keeper

import (
	"cosmossdk.io/core/store"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clock/types"
)

// Keeper of the clock store
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec

	contractKeeper wasmkeeper.PermissionedKeeper

	authority string
}

func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	contractKeeper wasmkeeper.PermissionedKeeper,
	authority string,
) *Keeper {

	return &Keeper{
		cdc:            cdc,
		storeService:   storeService,
		contractKeeper: contractKeeper,
		authority:      authority,
	}
}

// GetAuthority returns the x/clock module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// SetParams sets the x/clock module parameters.
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&p)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetParams returns the current x/clock module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

// GetContractKeeper returns the x/wasm module's contract keeper.
func (k Keeper) GetContractKeeper() wasmkeeper.PermissionedKeeper {
	return k.contractKeeper
}
