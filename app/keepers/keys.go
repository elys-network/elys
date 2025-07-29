package keepers

import (
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	accountedpoolmoduletypes "github.com/elys-network/elys/v7/x/accountedpool/types"
	ammmoduletypes "github.com/elys-network/elys/v7/x/amm/types"
	assetprofilemoduletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	burnermoduletypes "github.com/elys-network/elys/v7/x/burner/types"
	commitmentmoduletypes "github.com/elys-network/elys/v7/x/commitment/types"
	epochsmoduletypes "github.com/elys-network/elys/v7/x/epochs/types"
	estakingmoduletypes "github.com/elys-network/elys/v7/x/estaking/types"
	leveragelpmoduletypes "github.com/elys-network/elys/v7/x/leveragelp/types"
	masterchefmoduletypes "github.com/elys-network/elys/v7/x/masterchef/types"
	parametermoduletypes "github.com/elys-network/elys/v7/x/parameter/types"
	perpetualmoduletypes "github.com/elys-network/elys/v7/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/v7/x/stablestake/types"
	tiermoduletypes "github.com/elys-network/elys/v7/x/tier/types"
	tokenomicsmoduletypes "github.com/elys-network/elys/v7/x/tokenomics/types"
	tradeshieldmoduletypes "github.com/elys-network/elys/v7/x/tradeshield/types"
	vaultstypes "github.com/elys-network/elys/v7/x/vaults/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (appKeepers *AppKeepers) GenerateKeys() {
	// Define what keys will be used in the cosmos-sdk key/value store.
	// Cosmos-SDK modules each have a "key" that allows the application to reference what they've stored on the chain.
	appKeepers.keys = storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		crisistypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		ibcexported.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		icahosttypes.StoreKey,
		icacontrollertypes.StoreKey,
		capabilitytypes.StoreKey,
		ibchookstypes.StoreKey,
		feegrant.StoreKey,
		authz.ModuleName,
		group.StoreKey,
		consensusparamtypes.StoreKey,
		ccvconsumertypes.StoreKey,
		wasmTypes.StoreKey,
		packetforwardtypes.StoreKey,

		epochsmoduletypes.StoreKey,
		assetprofilemoduletypes.StoreKey,
		oracletypes.StoreKey,
		commitmentmoduletypes.StoreKey,
		tokenomicsmoduletypes.StoreKey,
		burnermoduletypes.StoreKey,
		accountedpoolmoduletypes.StoreKey,
		ammmoduletypes.StoreKey,
		parametermoduletypes.StoreKey,
		perpetualmoduletypes.StoreKey,
		stablestaketypes.StoreKey,
		leveragelpmoduletypes.StoreKey,
		masterchefmoduletypes.StoreKey,
		estakingmoduletypes.StoreKey,
		tiermoduletypes.StoreKey,
		tradeshieldmoduletypes.StoreKey,
		vaultstypes.StoreKey,
	)

	// Define transient store keys
	appKeepers.tkeys = storetypes.NewTransientStoreKeys(paramstypes.TStoreKey, ammmoduletypes.TStoreKey)

	// MemKeys are for information that is stored only in RAM.
	appKeepers.memKeys = storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
}

func (appKeepers *AppKeepers) GetKVStoreKey() map[string]*storetypes.KVStoreKey {
	return appKeepers.keys
}

func (appKeepers *AppKeepers) GetTransientStoreKey() map[string]*storetypes.TransientStoreKey {
	return appKeepers.tkeys
}

func (appKeepers *AppKeepers) GetMemoryStoreKey() map[string]*storetypes.MemoryStoreKey {
	return appKeepers.memKeys
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return appKeepers.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return appKeepers.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (appKeepers *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return appKeepers.memKeys[storeKey]
}
