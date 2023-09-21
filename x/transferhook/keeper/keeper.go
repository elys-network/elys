package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/transferhook/types"
)

type (
	Keeper struct {
		Cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace
		ammKeeper  ammkeeper.Keeper
	}
)

func NewKeeper(
	Cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ammKeeper ammkeeper.Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		Cdc:        Cdc,
		storeKey:   storeKey,
		paramstore: ps,
		ammKeeper:  ammKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
