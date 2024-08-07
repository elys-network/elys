package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	wasmbindingsclient "github.com/elys-network/elys/wasmbindings/client"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	assetprofilemodule "github.com/elys-network/elys/x/assetprofile"
	assetprofilemodulekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	assetprofilemoduletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentmodule "github.com/elys-network/elys/x/commitment"
	commitmentmodulekeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmentmoduletypes "github.com/elys-network/elys/x/commitment/types"
	epochsmodule "github.com/elys-network/elys/x/epochs"
	epochsmodulekeeper "github.com/elys-network/elys/x/epochs/keeper"
	epochsmoduletypes "github.com/elys-network/elys/x/epochs/types"
	exdistr "github.com/elys-network/elys/x/estaking/modules/distribution"
	exstaking "github.com/elys-network/elys/x/estaking/modules/staking"
	oraclemodule "github.com/elys-network/elys/x/oracle"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cast"

	tokenomicsmodule "github.com/elys-network/elys/x/tokenomics"
	tokenomicsmodulekeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicsmoduletypes "github.com/elys-network/elys/x/tokenomics/types"

	incentivemodule "github.com/elys-network/elys/x/incentive"
	incentivemodulekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivemoduletypes "github.com/elys-network/elys/x/incentive/types"

	burnermodule "github.com/elys-network/elys/x/burner"
	burnermodulekeeper "github.com/elys-network/elys/x/burner/keeper"
	burnermoduletypes "github.com/elys-network/elys/x/burner/types"

	ammmodule "github.com/elys-network/elys/x/amm"
	ammmodulekeeper "github.com/elys-network/elys/x/amm/keeper"
	ammmoduletypes "github.com/elys-network/elys/x/amm/types"

	parametermodule "github.com/elys-network/elys/x/parameter"
	parametermodulekeeper "github.com/elys-network/elys/x/parameter/keeper"
	parametermoduletypes "github.com/elys-network/elys/x/parameter/types"

	perpetualmodule "github.com/elys-network/elys/x/perpetual"
	perpetualmodulekeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualmoduletypes "github.com/elys-network/elys/x/perpetual/types"

	accountedpoolmodule "github.com/elys-network/elys/x/accountedpool"
	accountedpoolmodulekeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	accountedpoolmoduletypes "github.com/elys-network/elys/x/accountedpool/types"

	"github.com/elys-network/elys/x/transferhook"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
	transferhooktypes "github.com/elys-network/elys/x/transferhook/types"

	clockmodule "github.com/elys-network/elys/x/clock"
	clockmodulekeeper "github.com/elys-network/elys/x/clock/keeper"
	clockmoduletypes "github.com/elys-network/elys/x/clock/types"

	wasmmodule "github.com/CosmWasm/wasmd/x/wasm"
	wasmmodulekeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmmoduletypes "github.com/CosmWasm/wasmd/x/wasm/types"

	stablestake "github.com/elys-network/elys/x/stablestake"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"

	leveragelpmodule "github.com/elys-network/elys/x/leveragelp"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"

	estakingmodule "github.com/elys-network/elys/x/estaking"
	estakingmodulekeeper "github.com/elys-network/elys/x/estaking/keeper"
	estakingmoduletypes "github.com/elys-network/elys/x/estaking/types"
	masterchefmodule "github.com/elys-network/elys/x/masterchef"
	masterchefmodulekeeper "github.com/elys-network/elys/x/masterchef/keeper"
	masterchefmoduletypes "github.com/elys-network/elys/x/masterchef/types"

	tiermodule "github.com/elys-network/elys/x/tier"
	tiermodulekeeper "github.com/elys-network/elys/x/tier/keeper"
	tiermoduletypes "github.com/elys-network/elys/x/tier/types"

	ante "github.com/elys-network/elys/app/ante"

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
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasmmodule.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasmmodule.EnableAllProposals
		}
		return wasmmodule.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasmmodule.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
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

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	txConfig          client.TxConfig

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	AuthzKeeper           authzkeeper.Keeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper          ibcfeekeeper.Keeper
	IBCHooksKeeper        *ibchookskeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	WasmKeeper            wasmmodulekeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCFeeKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	EpochsKeeper       epochsmodulekeeper.Keeper
	AssetprofileKeeper assetprofilemodulekeeper.Keeper
	ScopedOracleKeeper capabilitykeeper.ScopedKeeper
	OracleKeeper       oraclekeeper.Keeper
	CommitmentKeeper   commitmentmodulekeeper.Keeper
	TokenomicsKeeper   tokenomicsmodulekeeper.Keeper
	IncentiveKeeper    incentivemodulekeeper.Keeper
	BurnerKeeper       burnermodulekeeper.Keeper
	AmmKeeper          ammmodulekeeper.Keeper
	ParameterKeeper    parametermodulekeeper.Keeper
	PerpetualKeeper    perpetualmodulekeeper.Keeper
	TransferhookKeeper transferhookkeeper.Keeper
	ContractKeeper     *wasmmodulekeeper.PermissionedKeeper
	ClockKeeper        clockmodulekeeper.Keeper

	AccountedPoolKeeper accountedpoolmodulekeeper.Keeper

	StablestakeKeeper stablestakekeeper.Keeper

	LeveragelpKeeper leveragelpmodulekeeper.Keeper

	MasterchefKeeper masterchefmodulekeeper.Keeper

	EstakingKeeper estakingmodulekeeper.Keeper

	TierKeeper tiermodulekeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// mm is the module manager
	mm *module.Manager

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
	enabledProposals []wasmmodule.ProposalType,
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

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, authz.ModuleName, banktypes.StoreKey, stakingtypes.StoreKey, crisistypes.StoreKey,
		slashingtypes.StoreKey, govtypes.StoreKey,
		paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey, evidencetypes.StoreKey,
		ibctransfertypes.StoreKey, icahosttypes.StoreKey, capabilitytypes.StoreKey, group.StoreKey,
		ibcfeetypes.StoreKey,
		icacontrollertypes.StoreKey,
		distrtypes.StoreKey,
		wasmmodule.StoreKey,
		consensusparamtypes.StoreKey,
		epochsmoduletypes.StoreKey,
		assetprofilemoduletypes.StoreKey,
		oracletypes.StoreKey,
		commitmentmoduletypes.StoreKey,
		tokenomicsmoduletypes.StoreKey,
		incentivemoduletypes.StoreKey,
		burnermoduletypes.StoreKey,
		accountedpoolmoduletypes.StoreKey,
		ammmoduletypes.StoreKey,
		parametermoduletypes.StoreKey,
		perpetualmoduletypes.StoreKey,
		transferhooktypes.StoreKey,
		clockmoduletypes.StoreKey,
		stablestaketypes.StoreKey,
		leveragelpmoduletypes.StoreKey,
		masterchefmoduletypes.StoreKey,
		estakingmoduletypes.StoreKey,
		tiermoduletypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, ammmoduletypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	app := &ElysApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		txConfig:          txConfig,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec, keys[consensusparamtypes.StoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	app.ParameterKeeper = *parametermodulekeeper.NewKeeper(
		appCodec,
		keys[parametermoduletypes.StoreKey],
		keys[parametermoduletypes.MemStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasmmodule.ModuleName)
	// this line is used by starport scaffolding # stargate/app/scopedKeeper

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32PrefixAccAddr,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authz.ModuleName],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.BlockedModuleAccountAddrs(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.AssetprofileKeeper = *assetprofilemodulekeeper.NewKeeper(
		appCodec,
		keys[assetprofilemoduletypes.StoreKey],
		keys[assetprofilemoduletypes.MemStoreKey],
		app.GetSubspace(assetprofilemoduletypes.ModuleName),
		&app.TransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	assetprofileModule := assetprofilemodule.NewAppModule(appCodec, app.AssetprofileKeeper, app.AccountKeeper, app.BankKeeper)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	commitmentKeeper := *commitmentmodulekeeper.NewKeeper(
		appCodec,
		keys[commitmentmoduletypes.StoreKey],
		keys[commitmentmoduletypes.MemStoreKey],
		app.GetSubspace(commitmentmoduletypes.ModuleName),

		app.BankKeeper,
		app.StakingKeeper,
		app.AssetprofileKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.TokenomicsKeeper = *tokenomicsmodulekeeper.NewKeeper(
		appCodec,
		keys[tokenomicsmoduletypes.StoreKey],
		keys[tokenomicsmoduletypes.MemStoreKey],
		app.GetSubspace(tokenomicsmoduletypes.ModuleName),
		&app.CommitmentKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	tokenomicsModule := tokenomicsmodule.NewAppModule(appCodec, app.TokenomicsKeeper, app.AccountKeeper, app.BankKeeper)

	app.EstakingKeeper = *estakingmodulekeeper.NewKeeper(
		appCodec,
		keys[estakingmoduletypes.StoreKey],
		keys[estakingmoduletypes.MemStoreKey],
		app.ParameterKeeper,
		app.StakingKeeper,
		&app.CommitmentKeeper,
		&app.DistrKeeper,
		app.AssetprofileKeeper,
		app.TokenomicsKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	estakingModule := estakingmodule.NewAppModule(appCodec, app.EstakingKeeper, app.AccountKeeper, app.BankKeeper)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		keys[slashingtypes.StoreKey],
		app.EstakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// ... other modules keepers

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// IBC Fee Module keeper
	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec, keys[ibcfeetypes.StoreKey],
		app.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper, app.AccountKeeper, app.BankKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	icaControllerKeeper := icacontrollerkeeper.NewKeeper(
		appCodec,
		keys[icacontrollertypes.StoreKey],
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)
	icaModule := ica.NewAppModule(&icaControllerKeeper, &app.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	scopedOracleKeeper := app.CapabilityKeeper.ScopeToModule(oracletypes.ModuleName)
	app.ScopedOracleKeeper = scopedOracleKeeper
	app.OracleKeeper = *oraclekeeper.NewKeeper(
		appCodec,
		keys[oracletypes.StoreKey],
		keys[oracletypes.MemStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.GetSubspace(oracletypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedOracleKeeper,
	)
	oracleModule := oraclemodule.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper)

	oracleIBCModule := oraclemodule.NewIBCModule(app.OracleKeeper)

	app.EpochsKeeper = *epochsmodulekeeper.NewKeeper(
		appCodec,
		keys[epochsmoduletypes.StoreKey],
	)

	app.AccountedPoolKeeper = *accountedpoolmodulekeeper.NewKeeper(
		appCodec,
		keys[accountedpoolmoduletypes.StoreKey],
		keys[accountedpoolmoduletypes.MemStoreKey],
		app.BankKeeper,
	)
	accountedPoolModule := accountedpoolmodule.NewAppModule(appCodec, app.AccountedPoolKeeper, app.AccountKeeper, app.BankKeeper)

	app.AmmKeeper = *ammmodulekeeper.NewKeeper(
		appCodec,
		keys[ammmoduletypes.StoreKey],
		tkeys[ammmoduletypes.TStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		&app.ParameterKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		app.OracleKeeper,
		&app.CommitmentKeeper,
		app.AssetprofileKeeper,
		app.AccountedPoolKeeper,
	)

	app.StablestakeKeeper = *stablestakekeeper.NewKeeper(
		appCodec,
		keys[stablestaketypes.StoreKey],
		keys[stablestaketypes.MemStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.GetSubspace(stablestaketypes.ModuleName),
		app.BankKeeper,
		&app.CommitmentKeeper,
		app.AssetprofileKeeper,
	)
	app.CommitmentKeeper = *commitmentKeeper.SetHooks(
		commitmentmodulekeeper.NewMultiCommitmentHooks(
			app.EstakingKeeper.CommitmentHooks(),
		),
	)
	commitmentModule := commitmentmodule.NewAppModule(appCodec, app.CommitmentKeeper, app.AccountKeeper, app.BankKeeper)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.AccountKeeper,
		app.CommitmentKeeper,
		app.EstakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.MasterchefKeeper = *masterchefmodulekeeper.NewKeeper(
		appCodec,
		keys[masterchefmoduletypes.StoreKey],
		keys[masterchefmoduletypes.MemStoreKey],
		app.GetSubspace(masterchefmoduletypes.ModuleName),
		app.ParameterKeeper,
		app.CommitmentKeeper,
		&app.AmmKeeper,
		app.OracleKeeper,
		app.AssetprofileKeeper,
		app.AccountedPoolKeeper,
		app.StablestakeKeeper,
		app.TokenomicsKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	masterchefModule := masterchefmodule.NewAppModule(appCodec, app.MasterchefKeeper, app.AccountKeeper, app.BankKeeper)

	app.IncentiveKeeper = *incentivemodulekeeper.NewKeeper(
		appCodec,
		keys[incentivemoduletypes.StoreKey],
		keys[incentivemoduletypes.MemStoreKey],
		app.ParameterKeeper,
		commitmentKeeper,
		app.StakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		&app.AmmKeeper,
		app.OracleKeeper,
		app.AssetprofileKeeper,
		app.AccountedPoolKeeper,
		app.StablestakeKeeper,
		app.TokenomicsKeeper,
		&app.MasterchefKeeper,
		&app.EstakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	incentiveModule := incentivemodule.NewAppModule(appCodec, app.IncentiveKeeper, app.EstakingKeeper, app.MasterchefKeeper, app.DistrKeeper, app.CommitmentKeeper)

	app.BurnerKeeper = *burnermodulekeeper.NewKeeper(
		appCodec,
		keys[burnermoduletypes.StoreKey],
		keys[burnermoduletypes.MemStoreKey],
		app.GetSubspace(burnermoduletypes.ModuleName),

		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	burnerModule := burnermodule.NewAppModule(appCodec, app.BurnerKeeper, app.AccountKeeper, app.BankKeeper)

	parameterModule := parametermodule.NewAppModule(appCodec, app.ParameterKeeper, app.AccountKeeper, app.BankKeeper)

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasmmodule.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	availableCapabilities := strings.Join(wasmbindingstypes.AllCapabilities(), ",")
	bankKeeper := app.BankKeeper.(bankkeeper.BaseKeeper)
	wasmOpts = append(
		wasmbindingsclient.RegisterCustomPlugins(
			&app.AccountedPoolKeeper,
			&app.AmmKeeper,
			&app.AssetprofileKeeper,
			&app.AccountKeeper,
			&bankKeeper,
			&app.BurnerKeeper,
			&app.ClockKeeper,
			&app.CommitmentKeeper,
			&app.EpochsKeeper,
			&app.IncentiveKeeper,
			&app.LeveragelpKeeper,
			&app.PerpetualKeeper,
			&app.OracleKeeper,
			&app.ParameterKeeper,
			&app.StablestakeKeeper,
			app.StakingKeeper,
			&app.TokenomicsKeeper,
			&app.TransferhookKeeper,
			&app.MasterchefKeeper,
			&app.EstakingKeeper,
			&app.TierKeeper,
		),
		wasmOpts...,
	)

	app.WasmKeeper = wasmmodule.NewKeeper(
		appCodec,
		keys[wasmmodule.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	app.TransferhookKeeper = *transferhookkeeper.NewKeeper(
		appCodec,
		keys[transferhooktypes.StoreKey],
		app.GetSubspace(transferhooktypes.ModuleName),
		app.AmmKeeper)
	transferhookModule := transferhook.NewAppModule(appCodec, app.TransferhookKeeper)
	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		keys[ibchookstypes.StoreKey],
	)
	app.IBCHooksKeeper = &hooksKeeper

	wasmHooks := ibchooks.NewWasmHooks(app.IBCHooksKeeper, &app.WasmKeeper, AccountAddressPrefix) // The contract keeper needs to be set later
	app.Ics20WasmHooks = &wasmHooks
	app.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		app.IBCKeeper.ChannelKeeper,
		app.Ics20WasmHooks,
	)

	// set the contract keeper for the Ics20WasmHooks
	app.ContractKeeper = wasmmodulekeeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	app.Ics20WasmHooks.ContractKeeper = &app.WasmKeeper

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	// The gov proposal types can be individually enabled
	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasmmoduletypes.RouterKey, wasmmodulekeeper.NewWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}
	govKeeper.SetLegacyRouter(govRouter)

	app.PerpetualKeeper = *perpetualmodulekeeper.NewKeeper(
		appCodec,
		keys[perpetualmoduletypes.StoreKey],
		keys[perpetualmoduletypes.MemStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		&app.AmmKeeper,
		app.BankKeeper,
		app.OracleKeeper,
		app.AssetprofileKeeper,
		&app.ParameterKeeper,
	)

	app.ClockKeeper = *clockmodulekeeper.NewKeeper(
		keys[clockmoduletypes.StoreKey],
		appCodec,
		app.GetSubspace(clockmoduletypes.ModuleName),
		*app.ContractKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	clockModule := clockmodule.NewAppModule(appCodec, app.ClockKeeper)

	app.LeveragelpKeeper = *leveragelpmodulekeeper.NewKeeper(
		appCodec,
		keys[leveragelpmoduletypes.StoreKey],
		keys[leveragelpmoduletypes.MemStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		&app.AmmKeeper,
		app.BankKeeper,
		app.OracleKeeper,
		app.StablestakeKeeper,
		app.CommitmentKeeper,
		app.AssetprofileKeeper,
		app.MasterchefKeeper,
	)
	leveragelpModule := leveragelpmodule.NewAppModule(appCodec, app.LeveragelpKeeper, app.AccountKeeper, app.BankKeeper)

	app.TierKeeper = *tiermodulekeeper.NewKeeper(
		appCodec,
		keys[tiermoduletypes.StoreKey],
		keys[tiermoduletypes.MemStoreKey],
		app.GetSubspace(tiermoduletypes.ModuleName),
		app.BankKeeper,
		app.OracleKeeper,
		app.AssetprofileKeeper,
		&app.AmmKeeper,
		app.EstakingKeeper,
		app.MasterchefKeeper,
		app.CommitmentKeeper,
		app.StakingKeeper,
		app.PerpetualKeeper,
		app.LeveragelpKeeper,
		app.StablestakeKeeper,
	)
	tierModule := tiermodule.NewAppModule(appCodec, app.TierKeeper, app.AccountKeeper, app.BankKeeper)

	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	/**** IBC Routing ****/

	// Sealing prevents other modules from creating scoped sub-keepers
	app.CapabilityKeeper.Seal()

	// Create fee enabled wasm ibc Stack
	var wasmStack ibcporttypes.IBCModule
	wasmStack = wasmmodule.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, app.IBCFeeKeeper)

	var transferStack ibcporttypes.IBCModule = transferIBCModule
	transferStack = transferhook.NewIBCModule(app.TransferhookKeeper, transferStack)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmmoduletypes.ModuleName, wasmStack).
		AddRoute(oracletypes.ModuleName, oracleIBCModule)
	// this line is used by starport scaffolding # ibc/app/router
	app.IBCKeeper.SetRouter(ibcRouter)

	/**** Module Hooks ****/

	// register hooks after all modules have been initialized

	app.StablestakeKeeper = *app.StablestakeKeeper.SetHooks(stablestakekeeper.NewMultiStableStakeHooks(
		app.MasterchefKeeper.StableStakeHooks(),
		app.TierKeeper.StableStakeHooks(),
	))
	stablestakeModule := stablestake.NewAppModule(appCodec, app.StablestakeKeeper, app.AccountKeeper, app.BankKeeper)

	app.EstakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			// insert staking hooks receivers here
			app.SlashingKeeper.Hooks(),
			app.DistrKeeper.Hooks(),
			app.EstakingKeeper.StakingHooks(),
			app.TierKeeper.StakingHooks(),
		),
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	app.AmmKeeper = *app.AmmKeeper.SetHooks(
		ammmoduletypes.NewMultiAmmHooks(
			// insert amm hooks receivers here
			app.PerpetualKeeper.AmmHooks(),
			app.LeveragelpKeeper.AmmHooks(),
			app.MasterchefKeeper.AmmHooks(),
			app.TierKeeper.AmmHooks(),
		),
	)
	ammModule := ammmodule.NewAppModule(appCodec, app.AmmKeeper, app.AccountKeeper, app.BankKeeper)

	app.EpochsKeeper = *app.EpochsKeeper.SetHooks(
		epochsmodulekeeper.NewMultiEpochHooks(
			// insert epoch hooks receivers here
			app.OracleKeeper.Hooks(),
			app.CommitmentKeeper.Hooks(),
			app.BurnerKeeper.Hooks(),
			app.PerpetualKeeper.Hooks(),
		),
	)
	epochsModule := epochsmodule.NewAppModule(appCodec, app.EpochsKeeper)

	app.PerpetualKeeper = *app.PerpetualKeeper.SetHooks(
		perpetualmoduletypes.NewMultiPerpetualHooks(
			// insert perpetual hooks receivers here
			app.AccountedPoolKeeper.PerpetualHooks(),
			app.TierKeeper.PerpetualHooks(),
		),
	)
	perpetualModule := perpetualmodule.NewAppModule(appCodec, app.PerpetualKeeper, app.AccountKeeper, app.BankKeeper)

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
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
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
		epochsModule,
		assetprofileModule,
		oracleModule,
		commitmentModule,
		tokenomicsModule,
		incentiveModule,
		burnerModule,
		ammModule,
		parameterModule,
		perpetualModule,
		accountedPoolModule,
		transferhookModule,
		clockModule,
		stablestakeModule,
		leveragelpModule,
		masterchefModule,
		estakingModule,
		tierModule,
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
		// this line is used by starport scaffolding # stargate/app/initGenesis
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

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
func (app *ElysApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *ElysApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *ElysApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
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
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
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
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService implements the Application.RegisterNodeService method.
func (app *ElysApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(assetprofilemoduletypes.ModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName)
	paramsKeeper.Subspace(commitmentmoduletypes.ModuleName)
	paramsKeeper.Subspace(tokenomicsmoduletypes.ModuleName)
	paramsKeeper.Subspace(burnermoduletypes.ModuleName)
	paramsKeeper.Subspace(perpetualmoduletypes.ModuleName)
	paramsKeeper.Subspace(transferhooktypes.ModuleName)
	paramsKeeper.Subspace(clockmoduletypes.ModuleName)
	paramsKeeper.Subspace(stablestaketypes.ModuleName)
	paramsKeeper.Subspace(leveragelpmoduletypes.ModuleName)
	paramsKeeper.Subspace(masterchefmoduletypes.ModuleName)
	paramsKeeper.Subspace(tiermoduletypes.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace

	return paramsKeeper
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
