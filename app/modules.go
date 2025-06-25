package app

import (
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	ccvstaking "github.com/cosmos/interchain-security/v6/x/ccv/democracy/staking"
	accountedpoolmodule "github.com/elys-network/elys/v6/x/accountedpool"
	accountedpoolmoduletypes "github.com/elys-network/elys/v6/x/accountedpool/types"
	ammmodule "github.com/elys-network/elys/v6/x/amm"
	ammmoduletypes "github.com/elys-network/elys/v6/x/amm/types"
	assetprofilemodule "github.com/elys-network/elys/v6/x/assetprofile"
	assetprofilemoduletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	burnermodule "github.com/elys-network/elys/v6/x/burner"
	burnermoduletypes "github.com/elys-network/elys/v6/x/burner/types"
	commitmentmodule "github.com/elys-network/elys/v6/x/commitment"
	commitmentmoduletypes "github.com/elys-network/elys/v6/x/commitment/types"
	epochsmodule "github.com/elys-network/elys/v6/x/epochs"
	epochsmoduletypes "github.com/elys-network/elys/v6/x/epochs/types"
	estakingmodule "github.com/elys-network/elys/v6/x/estaking"
	exdistr "github.com/elys-network/elys/v6/x/estaking/modules/distribution"
	estakingmoduletypes "github.com/elys-network/elys/v6/x/estaking/types"
	leveragelpmodule "github.com/elys-network/elys/v6/x/leveragelp"
	leveragelpmoduletypes "github.com/elys-network/elys/v6/x/leveragelp/types"
	masterchefmodule "github.com/elys-network/elys/v6/x/masterchef"
	masterchefmoduletypes "github.com/elys-network/elys/v6/x/masterchef/types"
	oraclemodule "github.com/elys-network/elys/v6/x/oracle"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	parametermodule "github.com/elys-network/elys/v6/x/parameter"
	parametermoduletypes "github.com/elys-network/elys/v6/x/parameter/types"
	perpetualmodule "github.com/elys-network/elys/v6/x/perpetual"
	perpetualmoduletypes "github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/elys-network/elys/v6/x/stablestake"
	stablestaketypes "github.com/elys-network/elys/v6/x/stablestake/types"
	tiermodule "github.com/elys-network/elys/v6/x/tier"
	tiermoduletypes "github.com/elys-network/elys/v6/x/tier/types"
	tokenomicsmodule "github.com/elys-network/elys/v6/x/tokenomics"
	tokenomicsmoduletypes "github.com/elys-network/elys/v6/x/tokenomics/types"
	tradeshieldmodule "github.com/elys-network/elys/v6/x/tradeshield"
	tradeshieldmoduletypes "github.com/elys-network/elys/v6/x/tradeshield/types"
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:                    nil,
	distrtypes.ModuleName:                         nil,
	icatypes.ModuleName:                           nil,
	ibchookstypes.ModuleName:                      nil,
	stakingtypes.BondedPoolName:                   {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:                {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                           {authtypes.Burner},
	oracletypes.ModuleName:                        {authtypes.Minter},
	ibctransfertypes.ModuleName:                   {authtypes.Minter, authtypes.Burner},
	ccvconsumertypes.ConsumerRedistributeName:     {authtypes.Burner},
	ccvconsumertypes.ConsumerToSendToProviderName: nil,

	minttypes.ModuleName: {authtypes.Minter}, // Need in writing test cases to initialize accounts with balances, otherwise no use

	commitmentmoduletypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	burnermoduletypes.ModuleName:     {authtypes.Burner},
	ammmoduletypes.ModuleName:        {authtypes.Minter, authtypes.Burner, authtypes.Staking},
	stablestaketypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	masterchefmoduletypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	wasmTypes.ModuleName:             {authtypes.Burner},
}

func appModules(
	app *ElysApp,
	appCodec codec.Codec,
	txConfig client.TxEncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app,
			txConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)), // always be last to make sure that it checks for all invariants and not only part of them
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.ConsumerKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		// Important: The idea is that the rewards that needs to be sent to provider we will do so in estaking and masterchef EndBlocker by sending it to ConsumerToSendToProviderName.
		// And the one that needs to be distributed on Consumer we will do there only by sending it to ConsumerRedistributeName. This requires that our distribution module uses ConsumerRedistributeName
		// This needs consumer_redistribution_fraction MUST be 1 as we are totally controlling the rewards. Also, this needs  ccvconsumer EndBlocker to be run after estaking and masterchef as it will move funds from
		// FeeCollector to ConsumerRedistributeName
		exdistr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.CommitmentKeeper, app.EstakingKeeper, &app.AssetprofileKeeper, ccvconsumertypes.ConsumerRedistributeName, app.GetSubspace(distrtypes.ModuleName)),
		ccvstaking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ibctm.NewAppModule(),
		params.NewAppModule(app.ParamsKeeper),
		transfer.NewAppModule(*app.TransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		app.ConsumerModule, // Not defining it here directly because ConsumerModule needed for IBC router
		ibchooks.NewAppModule(app.AccountKeeper),
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmTypes.ModuleName)),
		epochsmodule.NewAppModule(appCodec, *app.EpochsKeeper),
		assetprofilemodule.NewAppModule(appCodec, app.AssetprofileKeeper, app.AccountKeeper, app.BankKeeper),
		oraclemodule.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		commitmentmodule.NewAppModule(appCodec, *app.CommitmentKeeper, app.AccountKeeper, app.BankKeeper),
		tokenomicsmodule.NewAppModule(appCodec, app.TokenomicsKeeper, app.AccountKeeper, app.BankKeeper),
		burnermodule.NewAppModule(appCodec, app.BurnerKeeper, app.AccountKeeper, app.BankKeeper),
		ammmodule.NewAppModule(appCodec, *app.AmmKeeper, app.AccountKeeper, app.BankKeeper),
		parametermodule.NewAppModule(appCodec, app.ParameterKeeper, app.AccountKeeper, app.BankKeeper),
		stablestake.NewAppModule(appCodec, *app.StablestakeKeeper, app.AccountKeeper, app.BankKeeper),
		accountedpoolmodule.NewAppModule(appCodec, app.AccountedPoolKeeper, app.AccountKeeper, app.BankKeeper),
		leveragelpmodule.NewAppModule(appCodec, *app.LeveragelpKeeper, app.AccountKeeper, app.BankKeeper),
		masterchefmodule.NewAppModule(appCodec, app.MasterchefKeeper, app.AccountKeeper, app.BankKeeper),
		estakingmodule.NewAppModule(appCodec, *app.EstakingKeeper, app.AccountKeeper, app.BankKeeper),
		perpetualmodule.NewAppModule(appCodec, app.PerpetualKeeper, app.AccountKeeper, app.BankKeeper),
		tiermodule.NewAppModule(appCodec, *app.TierKeeper, app.AccountKeeper, app.BankKeeper),
		tradeshieldmodule.NewAppModule(appCodec, app.TradeshieldKeeper, app.AccountKeeper, app.BankKeeper),
	}
}

func newBasicManagerFromManager(app *ElysApp) module.BasicManager {
	basicManager := module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		})
	basicManager.RegisterLegacyAminoCodec(app.legacyAmino)
	basicManager.RegisterInterfaces(app.interfaceRegistry)
	return basicManager
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulations
func simulationModules(
	app *ElysApp,
	appCodec codec.Codec,
	_ bool,
) []module.AppModuleSimulation {
	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, nil, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(*app.TransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		app.ConsumerModule, // Not defining it here directly because ConsumerModule needed for IBC router
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		epochsmodule.NewAppModule(appCodec, *app.EpochsKeeper),
		assetprofilemodule.NewAppModule(appCodec, app.AssetprofileKeeper, app.AccountKeeper, app.BankKeeper),
		oraclemodule.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		commitmentmodule.NewAppModule(appCodec, *app.CommitmentKeeper, app.AccountKeeper, app.BankKeeper),
		tokenomicsmodule.NewAppModule(appCodec, app.TokenomicsKeeper, app.AccountKeeper, app.BankKeeper),
		burnermodule.NewAppModule(appCodec, app.BurnerKeeper, app.AccountKeeper, app.BankKeeper),
		ammmodule.NewAppModule(appCodec, *app.AmmKeeper, app.AccountKeeper, app.BankKeeper),
		parametermodule.NewAppModule(appCodec, app.ParameterKeeper, app.AccountKeeper, app.BankKeeper),
		stablestake.NewAppModule(appCodec, *app.StablestakeKeeper, app.AccountKeeper, app.BankKeeper),
		accountedpoolmodule.NewAppModule(appCodec, app.AccountedPoolKeeper, app.AccountKeeper, app.BankKeeper),
		leveragelpmodule.NewAppModule(appCodec, *app.LeveragelpKeeper, app.AccountKeeper, app.BankKeeper),
		masterchefmodule.NewAppModule(appCodec, app.MasterchefKeeper, app.AccountKeeper, app.BankKeeper),
		estakingmodule.NewAppModule(appCodec, *app.EstakingKeeper, app.AccountKeeper, app.BankKeeper),
		perpetualmodule.NewAppModule(appCodec, app.PerpetualKeeper, app.AccountKeeper, app.BankKeeper),
		tiermodule.NewAppModule(appCodec, *app.TierKeeper, app.AccountKeeper, app.BankKeeper),
		tradeshieldmodule.NewAppModule(appCodec, app.TradeshieldKeeper, app.AccountKeeper, app.BankKeeper),
	}
}

/*
orderBeginBlockers tells the app's module manager how to set the order of
BeginBlockers, which are run at the beginning of every block.

Interchain Security Requirements:
During begin block slashing happens after distr.BeginBlocker so that
there is nothing left over in the validator fee pool, so as to keep the
CanWithdrawInvariant invariant.
NOTE: staking module is required if HistoricalEntries param > 0
NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
*/

func orderBeginBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at the beginning
		epochsmoduletypes.ModuleName,
		distrtypes.ModuleName,
		stablestaketypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		ccvconsumertypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		assetprofilemoduletypes.ModuleName,
		oracletypes.ModuleName,
		commitmentmoduletypes.ModuleName,
		tokenomicsmoduletypes.ModuleName,
		burnermoduletypes.ModuleName,
		ammmoduletypes.ModuleName,
		parametermoduletypes.ModuleName,
		perpetualmoduletypes.ModuleName,
		accountedpoolmoduletypes.ModuleName,
		leveragelpmoduletypes.ModuleName,
		masterchefmoduletypes.ModuleName,
		estakingmoduletypes.ModuleName,
		tiermoduletypes.ModuleName,
		tradeshieldmoduletypes.ModuleName,
		wasmTypes.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
	}
}

/*
Interchain Security Requirements:
- provider.EndBlock gets validator updates from the staking module;
thus, staking.EndBlock must be executed before provider.EndBlock;
- creating a new consumer chain requires the following order,
CreateChildClient(), staking.EndBlock, provider.EndBlock;
thus, gov.EndBlock must be executed before staking.EndBlock
*/
func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the end
		epochsmoduletypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stablestaketypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		assetprofilemoduletypes.ModuleName,
		oracletypes.ModuleName,
		commitmentmoduletypes.ModuleName,
		tokenomicsmoduletypes.ModuleName,
		burnermoduletypes.ModuleName,
		ammmoduletypes.ModuleName,
		parametermoduletypes.ModuleName,
		perpetualmoduletypes.ModuleName,
		accountedpoolmoduletypes.ModuleName,
		leveragelpmoduletypes.ModuleName,
		masterchefmoduletypes.ModuleName,
		estakingmoduletypes.ModuleName,
		tiermoduletypes.ModuleName,
		tradeshieldmoduletypes.ModuleName,
		wasmTypes.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,

		// Must be called after estaking and masterchef
		ccvconsumertypes.ModuleName,
	}
}

/*
NOTE: The genutils module must occur after staking so that pools are
properly initialized with tokens from genesis accounts.
NOTE: The genutils module must also occur after auth so that it can access the params from auth.
NOTE: Capability module must occur first so that it can initialize any capabilities
so that other modules that want to create or claim capabilities afterwards in InitChain
can do so safely.
*/
func orderInitBlockers() []string {
	return []string{
		parametermoduletypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		commitmentmoduletypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		epochsmoduletypes.ModuleName,
		stablestaketypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		ibctransfertypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ccvconsumertypes.ModuleName,
		consensusparamtypes.ModuleName,
		assetprofilemoduletypes.ModuleName,
		oracletypes.ModuleName,
		tokenomicsmoduletypes.ModuleName,
		burnermoduletypes.ModuleName,
		ammmoduletypes.ModuleName,
		perpetualmoduletypes.ModuleName,
		accountedpoolmoduletypes.ModuleName,
		leveragelpmoduletypes.ModuleName,
		masterchefmoduletypes.ModuleName,
		estakingmoduletypes.ModuleName,
		tiermoduletypes.ModuleName,
		tradeshieldmoduletypes.ModuleName,
		// wasm after ibc transfer
		wasmTypes.ModuleName,
		// ibc_hooks after auth keeper
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		// crisis needs to be last so that the genesis state is consistent
		// when it checks invariants
		crisistypes.ModuleName,
	}
}
