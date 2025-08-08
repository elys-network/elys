package keeper

import (
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

type Keeper struct {
	cdc                   codec.BinaryCodec
	storeService          store.KVStoreService
	transientStoreService store.TransientStoreService
	authority             string

	bankKeeper   types.BankKeeper
	oracleKeeper types.OracleKeeper

	// Caches for performance
	marketCache *MarketCache
	priceCache  *PriceCache

	// In-memory orderbook for efficient matching
	memoryOrderBook *MemoryOrderBook
}

var _ types.MsgServer = Keeper{}

func (k Keeper) GetAuthority() string {
	return k.authority
}

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
		marketCache:           NewMarketCache(),
		priceCache:            NewPriceCache(),
		memoryOrderBook:       NewMemoryOrderBook(),
	}
}

// InitializeMemoryOrderBook initializes the in-memory orderbook from chain state
func (k *Keeper) InitializeMemoryOrderBook(ctx sdk.Context) {
	if !k.memoryOrderBook.IsInitialized() {
		k.memoryOrderBook.InitializeFromState(*k, ctx)
	}
}
