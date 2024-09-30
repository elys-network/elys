package app

import (
	"fmt"
	"github.com/elys-network/elys/app/keepers"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"

	abci "github.com/cometbft/cometbft/abci/types"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
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
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	assetprofilemodule "github.com/elys-network/elys/x/assetprofile"
	assetprofilemoduletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentmodule "github.com/elys-network/elys/x/commitment"
	commitmentmoduletypes "github.com/elys-network/elys/x/commitment/types"
	epochsmodule "github.com/elys-network/elys/x/epochs"
	epochsmoduletypes "github.com/elys-network/elys/x/epochs/types"
	exdistr "github.com/elys-network/elys/x/estaking/modules/distribution"
	exstaking "github.com/elys-network/elys/x/estaking/modules/staking"
	incentivemodule "github.com/elys-network/elys/x/incentive"
	incentivemoduletypes "github.com/elys-network/elys/x/incentive/types"
	oraclemodule "github.com/elys-network/elys/x/oracle"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	tokenomicsmodule "github.com/elys-network/elys/x/tokenomics"
	tokenomicsmoduletypes "github.com/elys-network/elys/x/tokenomics/types"

	burnermodule "github.com/elys-network/elys/x/burner"
	burnermoduletypes "github.com/elys-network/elys/x/burner/types"

	ammmodule "github.com/elys-network/elys/x/amm"
	ammmoduletypes "github.com/elys-network/elys/x/amm/types"

	parametermodule "github.com/elys-network/elys/x/parameter"
	parametermoduletypes "github.com/elys-network/elys/x/parameter/types"

	perpetualmodule "github.com/elys-network/elys/x/perpetual"
	perpetualmoduletypes "github.com/elys-network/elys/x/perpetual/types"

	accountedpoolmodule "github.com/elys-network/elys/x/accountedpool"
	accountedpoolmoduletypes "github.com/elys-network/elys/x/accountedpool/types"

	"github.com/elys-network/elys/x/transferhook"
	transferhooktypes "github.com/elys-network/elys/x/transferhook/types"

	clockmodule "github.com/elys-network/elys/x/clock"
	clockmoduletypes "github.com/elys-network/elys/x/clock/types"

	wasmmodule "github.com/CosmWasm/wasmd/x/wasm"
	wasmmodulekeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmmoduletypes "github.com/CosmWasm/wasmd/x/wasm/types"

	stablestake "github.com/elys-network/elys/x/stablestake"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"

	leveragelpmodule "github.com/elys-network/elys/x/leveragelp"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"

	estakingmodule "github.com/elys-network/elys/x/estaking"
	estakingmoduletypes "github.com/elys-network/elys/x/estaking/types"
	masterchefmodule "github.com/elys-network/elys/x/masterchef"
	masterchefmoduletypes "github.com/elys-network/elys/x/masterchef/types"

	tiermodule "github.com/elys-network/elys/x/tier"
	tiermoduletypes "github.com/elys-network/elys/x/tier/types"

	ante "github.com/elys-network/elys/app/ante"

	tradeshieldmodule "github.com/elys-network/elys/x/tradeshield"
	tradeshieldmoduletypes "github.com/elys-network/elys/x/tradeshield/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	"github.com/elys-network/elys/docs"
)

const (
	AccountAddressPrefix = "elys"
	Name                 = "elys"

	// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
	ProposalsEnabled = "false"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""
)

// this line is used by starport scaffolding # stargate/wasm/app/enabledProposals

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
	)

	return govProposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		solomachine.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},
		wasmmodule.AppModuleBasic{},
		epochsmodule.AppModuleBasic{},
		assetprofilemodule.AppModuleBasic{},
		oraclemodule.AppModuleBasic{},
		commitmentmodule.AppModuleBasic{},
		tokenomicsmodule.AppModuleBasic{},
		incentivemodule.AppModuleBasic{},
		burnermodule.AppModuleBasic{},
		ammmodule.AppModuleBasic{},
		parametermodule.AppModuleBasic{},
		perpetualmodule.AppModuleBasic{},
		accountedpoolmodule.AppModuleBasic{},
		transferhook.AppModuleBasic{},
		clockmodule.AppModuleBasic{},
		stablestake.AppModuleBasic{},
		leveragelpmodule.AppModuleBasic{},
		masterchefmodule.AppModuleBasic{},
		estakingmodule.AppModuleBasic{},
		tiermodule.AppModuleBasic{},
		tradeshieldmodule.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:       nil,
		distrtypes.ModuleName:            nil,
		icatypes.ModuleName:              nil,
		ibcfeetypes.ModuleName:           nil,
		minttypes.ModuleName:             {authtypes.Minter},
		stakingtypes.BondedPoolName:      {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:   {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:              {authtypes.Burner},
		ibctransfertypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		commitmentmoduletypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		burnermoduletypes.ModuleName:     {authtypes.Burner},
		incentivemoduletypes.ModuleName:  nil,
		ammmoduletypes.ModuleName:        {authtypes.Minter, authtypes.Burner, authtypes.Staking},
		wasmmoduletypes.ModuleName:       {authtypes.Burner},
		stablestaketypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		masterchefmoduletypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		// this line is used by starport scaffolding # stargate/app/maccPerms
	}
)

var (
	_ runtime.AppI            = (*ElysApp)(nil)
	_ servertypes.Application = (*ElysApp)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type ElysApp struct {
	*baseapp.BaseApp
	keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	txConfig          client.TxConfig

	// keys to access the substores

	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// mm is the module manager
	mm           *module.Manager
	ModuleBasics module.BasicManager

	// sm is the simulation manager
	sm           *module.SimulationManager
	configurator module.Configurator

	// Middleware wrapper
	Ics20WasmHooks   *ibchooks.WasmHooks
	HooksICS4Wrapper ibchooks.ICS4Middleware
}

// New returns a reference to an initialized blockchain app
func NewElysApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmmodule.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *ElysApp {

	encodingConfig := MakeEncodingConfig()
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		txConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)

	//
	// NOTE: This is a temporary fix to allow the version command to work with the current version
	// of Elys. This should be removed once the version command is refactored.
	if !strings.HasPrefix(version.Version, "v") {
		version.Version = "v" + version.Version
	}
	bApp.SetVersion(version.Version)

	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	// load state streaming if enabled
	if err := bApp.RegisterStreamingServices(appOpts, keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	app := &ElysApp{
		BaseApp:           bApp,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		txConfig:          txConfig,
	}

	/**** Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
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
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		exdistr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.CommitmentKeeper, &app.EstakingKeeper, &app.AssetprofileKeeper, authtypes.FeeCollectorName, app.GetSubspace(distrtypes.ModuleName)),
		exstaking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		wasmmodule.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmmoduletypes.ModuleName)),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),

		transferModule,
		icaModule,

		epochsmodule.NewAppModule(appCodec, *app.EpochsKeeper),
		assetprofilemodule.NewAppModule(appCodec, app.AssetprofileKeeper, app.AccountKeeper, app.BankKeeper),
		oraclemodule.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		commitmentmodule.NewAppModule(appCodec, *app.CommitmentKeeper, app.AccountKeeper, app.BankKeeper),
		tokenomicsmodule.NewAppModule(appCodec, app.TokenomicsKeeper, app.AccountKeeper, app.BankKeeper),
		incentivemodule.NewAppModule(appCodec, app.IncentiveKeeper, *app.EstakingKeeper, app.MasterchefKeeper, app.DistrKeeper, *app.CommitmentKeeper),
		burnermodule.NewAppModule(appCodec, app.BurnerKeeper, app.AccountKeeper, app.BankKeeper),
		ammmodule.NewAppModule(appCodec, *app.AmmKeeper, app.AccountKeeper, app.BankKeeper),
		parametermodule.NewAppModule(appCodec, app.ParameterKeeper, app.AccountKeeper, app.BankKeeper),
		stablestake.NewAppModule(appCodec, *app.StablestakeKeeper, app.AccountKeeper, app.BankKeeper),
		accountedpoolmodule.NewAppModule(appCodec, app.AccountedPoolKeeper, app.AccountKeeper, app.BankKeeper),
		transferhook.NewAppModule(appCodec, app.TransferhookKeeper),
		clockmodule.NewAppModule(appCodec, app.ClockKeeper),
		leveragelpmodule.NewAppModule(appCodec, *app.LeveragelpKeeper, app.AccountKeeper, app.BankKeeper),
		masterchefmodule.NewAppModule(appCodec, app.MasterchefKeeper, app.AccountKeeper, app.BankKeeper),
		estakingmodule.NewAppModule(appCodec, *app.EstakingKeeper, app.AccountKeeper, app.BankKeeper),
		perpetualmodule.NewAppModule(appCodec, app.PerpetualKeeper, app.AccountKeeper, app.BankKeeper),
		tiermodule.NewAppModule(appCodec, app.TierKeeper, app.AccountKeeper, app.BankKeeper),
		tradeshieldmodule.NewAppModule(appCodec, app.TradeshieldKeeper, app.AccountKeeper, app.BankKeeper),
		// this line is used by starport scaffolding # stargate/app/appModule
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at the beginning
		epochsmoduletypes.ModuleName,
		distrtypes.ModuleName,
		stablestaketypes.ModuleName,
		incentivemoduletypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
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
		wasmmodule.ModuleName,
		accountedpoolmoduletypes.ModuleName,
		transferhooktypes.ModuleName,
		clockmoduletypes.ModuleName,
		leveragelpmoduletypes.ModuleName,
		masterchefmoduletypes.ModuleName,
		estakingmoduletypes.ModuleName,
		tiermoduletypes.ModuleName,
		tradeshieldmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/beginBlockers
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the end
		epochsmoduletypes.ModuleName,
		clockmoduletypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stablestaketypes.ModuleName,
		incentivemoduletypes.ModuleName,
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
		wasmmoduletypes.ModuleName,
		accountedpoolmoduletypes.ModuleName,
		transferhooktypes.ModuleName,
		leveragelpmoduletypes.ModuleName,
		masterchefmoduletypes.ModuleName,
		estakingmoduletypes.ModuleName,
		tiermoduletypes.ModuleName,
		tradeshieldmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/endBlockers
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder := []string{
		parametermoduletypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		commitmentmoduletypes.ModuleName,
		distrtypes.ModuleName,
		epochsmoduletypes.ModuleName,
		stablestaketypes.ModuleName,
		incentivemoduletypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
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
		tokenomicsmoduletypes.ModuleName,
		burnermoduletypes.ModuleName,
		ammmoduletypes.ModuleName,
		perpetualmoduletypes.ModuleName,
		wasmmoduletypes.ModuleName,
		accountedpoolmoduletypes.ModuleName,
		transferhooktypes.ModuleName,
		clockmoduletypes.ModuleName,
		leveragelpmoduletypes.ModuleName,
		masterchefmoduletypes.ModuleName,
		estakingmoduletypes.ModuleName,
		tiermoduletypes.ModuleName,
		tradeshieldmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/initGenesis
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	// Uncomment if you want to set a custom migration order here.
	customOrder := module.DefaultMigrationsOrder(app.mm.ModuleNames())
	stablestakeIndex := -1
	leveragelpIndex := -1
	for i := range customOrder {
		if customOrder[i] == stablestaketypes.ModuleName {
			stablestakeIndex = i
		}
		if customOrder[i] == leveragelpmoduletypes.ModuleName {
			leveragelpIndex = i
		}
	}
	if stablestakeIndex != -1 && leveragelpIndex != -1 {
		customOrder[leveragelpIndex] = stablestaketypes.ModuleName
		customOrder[stablestakeIndex] = leveragelpmoduletypes.ModuleName
	}
	app.mm.SetOrderMigrations(customOrder...)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// create the simulation manager and define the order of the modules for deterministic simulations
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)
	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	SetupHandlers(app)
	app.setAnteHandler(encodingConfig.TxConfig, wasmConfig, keys[wasmmodule.StoreKey])

	// must be before Loading version
	// requires the snapshot store to be created and registered as a BaseAppOption
	// see cmd/wasmd/root.go: 206 - 214 approx
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmmodulekeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper

	// In v0.46, the SDK introduces _postHandlers_. PostHandlers are like
	// antehandlers, but are run _after_ the `runMsgs` execution. They are also
	// defined as a chain, and have the same signature as antehandlers.
	//
	// In baseapp, postHandlers are run in the same store branch as `runMsgs`,
	// meaning that both `runMsgs` and `postHandler` state will be committed if
	// both are successful, and both will be reverted if any of the two fails.
	//
	// The SDK exposes a default postHandlers chain, which comprises of only
	// one decorator: the Transaction Tips decorator. However, some chains do
	// not need it by default, so feel free to comment the next line if you do
	// not need tips.
	// To read more about tips:
	// https://docs.cosmos.network/main/core/tips.html
	//
	// Please note that changing any of the anteHandler or postHandler chain is
	// likely to be a state-machine breaking change, which needs a coordinated
	// upgrade.
	app.setPostHandler()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			logger.Error("error on loading last version", "err", err)
			os.Exit(1)
		}
		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
		parameters := app.ParameterKeeper.GetParams(ctx)
		wasmConfiguration(parameters)
		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}

	return app
}

func (app *ElysApp) setAnteHandler(txConfig client.TxConfig, wasmConfig wasmmoduletypes.WasmConfig, txCounterStoreKey storetypes.StoreKey) {
	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			HandlerOptions: sdkante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				FeegrantKeeper:  app.FeeGrantKeeper,
				SigGasConsumer:  sdkante.DefaultSigVerificationGasConsumer,
				TxFeeChecker:    ante.CheckTxFeeWithValidatorMinGasPrices,
			},
			StakingKeeper:     app.StakingKeeper,
			IBCKeeper:         app.IBCKeeper,
			BankKeeper:        app.BankKeeper,
			Cdc:               app.appCodec,
			ParameterKeeper:   app.ParameterKeeper,
			WasmConfig:        &wasmConfig,
			TXCounterStoreKey: txCounterStoreKey,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}
	app.SetAnteHandler(anteHandler)
}

func (app *ElysApp) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)
}

// Name returns the name of the App
func (app *ElysApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *ElysApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *ElysApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *ElysApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		panic(err)
	}

	response, err := app.mm.InitGenesis(ctx, app.appCodec, genesisState)
	if err != nil {
		panic(err)
	}

	return response, nil
}

// LoadHeight loads a particular height
func (app *ElysApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *ElysApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *ElysApp) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ElysApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ElysApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *ElysApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ElysApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ElysApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *ElysApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ElysApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *ElysApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *ElysApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *ElysApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService implements the Application.RegisterNodeService method.
func (app *ElysApp) RegisterNodeService(clientCtx client.Context, config config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), config)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// SimulationManager returns the app SimulationManager
func (app *ElysApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// TxConfig returns App's TxConfig.
func (app *ElysApp) TxConfig() client.TxConfig {
	return app.txConfig
}

// Configurator get app configurator
func (app *ElysApp) Configurator() module.Configurator {
	return app.configurator
}

// ModuleManager returns the app ModuleManager
func (app *ElysApp) ModuleManager() *module.Manager {
	return app.mm
}
