package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/elys-network/elys/x/accountedpool/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		amm        types.AmmKeeper
		margin     types.MarginKeeper
		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	amm types.AmmKeeper,
	margin types.MarginKeeper,
	bk types.BankKeeper,
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
		amm:        amm,
		margin:     margin,
		bankKeeper: bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Get accounted pool balance
func (k Keeper) GetAccountedBalance(ctx sdk.Context, poolId uint64, denom string) sdk.Int {
	pool, found := k.GetAccountedPool(ctx, poolId)
	if !found {
		return sdk.ZeroInt()
	}

	for _, asset := range pool.PoolAssets {
		if asset.Token.Denom == denom {
			return asset.Token.Amount
		}
	}

	return sdk.ZeroInt()
}
