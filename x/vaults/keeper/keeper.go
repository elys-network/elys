package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tierkeeper "github.com/elys-network/elys/x/tier/keeper"
	"github.com/elys-network/elys/x/vaults/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority     string
		bk            types.BankKeeper
		tierKeeper    *tierkeeper.Keeper
		amm           types.AmmKeeper
		commitement   types.CommitmentKeeper
		accountKeeper types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	bk types.BankKeeper,
	tierKeeper *tierkeeper.Keeper,
	amm types.AmmKeeper,
	commitement types.CommitmentKeeper,
	accountKeeper types.AccountKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		authority:     authority,
		bk:            bk,
		tierKeeper:    tierKeeper,
		amm:           amm,
		commitement:   commitement,
		accountKeeper: accountKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
