package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/accountedpool/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService storetypes.KVStoreService
		bankKeeper   types.BankKeeper
		oracleKeeper types.OracleKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	bk types.BankKeeper,
	oracleKeeper types.OracleKeeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		bankKeeper:   bk,
		oracleKeeper: oracleKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Get accounted pool balance
func (k Keeper) GetAccountedBalance(ctx sdk.Context, poolId uint64, denom string) sdkmath.Int {
	pool, found := k.GetAccountedPool(ctx, poolId)
	if !found {
		return sdkmath.ZeroInt()
	}

	for _, token := range pool.TotalTokens {
		if token.Denom == denom {
			return token.Amount
		}
	}

	return sdkmath.ZeroInt()
}
