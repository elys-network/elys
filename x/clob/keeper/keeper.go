package keeper

import (
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/elys-network/elys/x/clob/types"
)

type Keeper struct {
	cdc                   codec.BinaryCodec
	storeService          store.KVStoreService
	transientStoreService store.TransientStoreService
	authority             string

	bankKeeper   types.BankKeeper
	oracleKeeper types.OracleKeeper
}

var _ types.MsgServer = Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	transientStoreService store.TransientStoreService,
	authority string,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
) *Keeper {
	return &Keeper{
		cdc:                   cdc,
		storeService:          storeService,
		transientStoreService: transientStoreService,
		authority:             authority,
		bankKeeper:            bankKeeper,
		oracleKeeper:          oracleKeeper,
	}
}
