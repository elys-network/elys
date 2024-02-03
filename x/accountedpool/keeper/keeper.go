package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/accountedpool/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	bk types.BankKeeper,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		bankKeeper: bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Get accounted pool balance
func (k Keeper) GetAccountedBalance(ctx sdk.Context, poolId uint64, denom string) math.Int {
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
