package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/gogoproto/proto"
	"github.com/elys-network/elys/v7/app/keepers"
	leveragelpmoduletypes "github.com/elys-network/elys/v7/x/leveragelp/types"
	stablestaketypes "github.com/elys-network/elys/v7/x/stablestake/types"
	"github.com/spf13/cast"

	abci "github.com/cometbft/cometbft/abci/types"
	tmos "github.com/cometbft/cometbft/libs/os"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasm "github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	dbm "github.com/cosmos/cosmos-db"
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
	sigtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/version"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	"github.com/elys-network/elys/v7/app/ante"
	clobabci "github.com/elys-network/elys/v7/x/clob/abci"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	"github.com/elys-network/elys/v7/docs"
)

const (
	AccountAddressPrefix = "elys"
	Name                 = "elys"
)

// this line is used by starport scaffolding

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

	invCheckPeriod uint

	// mm is the module manager
	mm           *module.Manager
	ModuleBasics module.BasicManager

	// sm is the simulation manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

// New returns a reference to an initialized blockchain app
func NewElysApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *ElysApp {

	encodingConfig := MakeEncodingConfig()
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

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

	app := &ElysApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		txConfig:          txConfig,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
	}

	moduleAccountAddresses := app.ModuleAccountAddrs()

	app.AppKeepers = keepers.NewAppKeeper(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		moduleAccountAddresses,
		app.BlockedModuleAccountAddrs(),
		skipUpgradeHeights,
		homePath,
		invCheckPeriod,
		logger,
		appOpts,
		AccountAddressPrefix,
		wasmOpts,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(appModules(app, appCodec, txConfig, skipGenesisInvariants)...)
	app.ModuleBasics = newBasicManagerFromManager(app)

	enabledSignModes := append([]sigtypes.SignMode(nil), authtx.DefaultSignModes...)
	enabledSignModes = append(enabledSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL)

	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfig, err := authtx.NewTxConfigWithOptions(
		appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}
	app.txConfig = txConfig

	// NOTE: upgrade module is required to be prioritized
	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	// Tell the app's module manager how to set the order of BeginBlockers, which are run at the beginning of every block.
	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)

	app.mm.SetOrderEndBlockers(orderEndBlockers()...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)

	// Uncomment if you want to set a custom migration order here.
	allModules := []string{}
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			allModules = append(allModules, moduleName)
		}
	}
	// Ensure "stablestake" appears before "leveragelp"
	stablestakeIndex, leveragelpIndex := -1, -1
	for i, name := range allModules {
		if name == stablestaketypes.ModuleName {
			stablestakeIndex = i
		} else if name == leveragelpmoduletypes.ModuleName {
			leveragelpIndex = i
		}
	}

	// Swap positions if needed
	if stablestakeIndex > leveragelpIndex && leveragelpIndex != -1 {
		allModules[stablestakeIndex], allModules[leveragelpIndex] = allModules[leveragelpIndex], allModules[stablestakeIndex]
	}
	app.mm.SetOrderMigrations(allModules...)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(simulationModules(app, appCodec, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{

			HandlerOptions: sdkante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				SigGasConsumer:  sdkante.DefaultSigVerificationGasConsumer,
				TxFeeChecker:    ante.CheckTxFeeWithValidatorMinGasPrices,
			},

			BankKeeper:            app.BankKeeper,
			ParameterKeeper:       app.ParameterKeeper,
			Cdc:                   appCodec,
			IBCKeeper:             app.IBCKeeper,
			StakingKeeper:         app.StakingKeeper,
			ConsumerKeeper:        app.ConsumerKeeper,
			WasmConfig:            &wasmConfig,
			TXCounterStoreService: runtime.NewKVStoreService(app.AppKeepers.GetKVStoreKey()[wasmTypes.StoreKey]),
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	// Set up CLOB vote extensions and proposal handlers
	clobProposalHandler := clobabci.NewProposalHandler(
		app.Logger(),
		app.ClobKeeper,
		txConfig.TxEncoder(),
		txConfig.TxDecoder(),
	)
	app.SetPrepareProposal(clobProposalHandler.PrepareProposalHandler())
	app.SetProcessProposal(clobProposalHandler.ProcessProposalHandler())

	clobVoteExtensionsHandler := clobabci.NewVoteExtensionHandler(
		app.Logger(),
		app.ClobKeeper,
	)
	app.SetExtendVoteHandler(clobVoteExtensionsHandler.ExtendVoteHandler())
	app.SetVerifyVoteExtensionHandler(clobVoteExtensionsHandler.VerifyVoteExtensionHandler())

	// Uncomment to enable oracle vote extensions as well
	//oracleProposalHandler := oracleabci.NewProposalHandler(
	//	app.Logger(),
	//	app.OracleKeeper,
	//	app.StakingKeeper,
	//)
	//app.SetPrepareProposal(oracleProposalHandler.PrepareProposalHandler())
	//app.SetProcessProposal(oracleProposalHandler.ProcessProposalHandler())
	//
	//oracleVoteExtensionsHandler := oracleabci.NewVoteExtensionHandler(
	//	app.Logger(),
	//	app.OracleKeeper,
	//)
	//app.SetExtendVoteHandler(oracleVoteExtensionsHandler.ExtendVoteHandler())
	//app.SetVerifyVoteExtensionHandler(oracleVoteExtensionsHandler.VerifyVoteExtensionHandler())

	// set ante and post handlers
	app.SetAnteHandler(anteHandler)
	app.setPostHandler()
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.setUpgradeHandler()
	app.setUpgradeStore()

	// At startup, after all modules have been registered, check that all prot
	// annotations are correct.
	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}
	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		// Once we switch to using protoreflect-based antehandlers, we might
		// want to panic here instead of logging a warning.
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
	}

	return app
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

func (app *ElysApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

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
	// Remove the fee-pool from the group of blocked recipient addresses in bank
	// this is required for the consumer chain to be able to send tokens to
	// the provider chain
	delete(modAccAddrs, authtypes.NewModuleAddress(ccvconsumertypes.ConsumerToSendToProviderName).String())
	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ElysApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
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
	return app.GetKVStoreKeys()[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ElysApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.GetTransientStoreKeys()[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *ElysApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.GetMemoryStoreKey()[storeKey]
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
	app.ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

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

func (app *ElysApp) OnTxSucceeded(_ sdk.Context, _, _ string, _ []byte, _ []byte) {
}

func (app *ElysApp) OnTxFailed(_ sdk.Context, _, _ string, _ []byte, _ []byte) {
}

// AutoCliOpts returns the autocli options for the app.
func (app *ElysApp) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

// TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *ElysApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetTxConfig implements the TestingApp interface.
func (app *ElysApp) GetTxConfig() client.TxConfig {
	return app.txConfig
}
