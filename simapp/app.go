package simapp

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	_ "cosmossdk.io/x/upgrade"
	_ "github.com/cosmos/cosmos-sdk/x/auth"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	_ "github.com/cosmos/cosmos-sdk/x/bank"
	_ "github.com/cosmos/cosmos-sdk/x/consensus"
	_ "github.com/cosmos/cosmos-sdk/x/params"
	_ "github.com/cosmos/cosmos-sdk/x/staking"

	// Cosmos Modules
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	// IBC Modules
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	// Custom Modules
	accountedpoolkeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	assetprofilekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	burnerkeeper "github.com/elys-network/elys/x/burner/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	epochskeeper "github.com/elys-network/elys/x/epochs/keeper"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	tierkeeper "github.com/elys-network/elys/x/tier/keeper"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tradeshieldkeeper "github.com/elys-network/elys/x/tradeshield/keeper"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
)

var DefaultNodeHome string

//go:embed app.yaml
var AppConfigYAML []byte

var (
	_ runtime.AppI            = (*SimApp)(nil)
	_ servertypes.Application = (*SimApp)(nil)
)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// Cosmos Modules
	AccountKeeper   authkeeper.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	ConsensusKeeper consensuskeeper.Keeper
	ParamsKeeper    paramskeeper.Keeper
	StakingKeeper   *stakingkeeper.Keeper
	UpgradeKeeper   *upgradekeeper.Keeper
	// IBC Modules
	CapabilityKeeper     *capabilitykeeper.Keeper
	IBCKeeper            *ibckeeper.Keeper
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	TransferKeeper       transferkeeper.Keeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	// Custom Modules
	AccountedPoolKeeper *accountedpoolkeeper.Keeper
	AMMKeeper           *ammkeeper.Keeper
	AssetProfileKeeper  *assetprofilekeeper.Keeper
	BurnerKeeper        *burnerkeeper.Keeper
	CommitmentKeeper    *commitmentkeeper.Keeper
	EpochsKeeper        *epochskeeper.Keeper
	EstakingKeeper      *estakingkeeper.Keeper
	LeverageLpKeeper    *leveragelpkeeper.Keeper
	MasterchefKeeper    *masterchefkeeper.Keeper
	OracleKeeper        *oraclekeeper.Keeper
	ParameterKeeper     *parameterkeeper.Keeper
	PerpetualKeeper     *perpetualkeeper.Keeper
	StablestakeKeeper   *stablestakekeeper.Keeper
	TierKeeper          *tierkeeper.Keeper
	TokenomicsKeeper    *tokenomicskeeper.Keeper
	TradeshieldKeeper   *tradeshieldkeeper.Keeper
	TransferHookKeeper  *transferhookkeeper.Keeper
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appconfig.LoadYAML(AppConfigYAML),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			},
		),
	)
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*SimApp, error) {
	var (
		app        = &SimApp{}
		appBuilder *runtime.AppBuilder
	)

	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Supply(
				logger,
				appOpts,
			),
		),
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		// Cosmos Modules
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.ConsensusKeeper,
		&app.ParamsKeeper,
		&app.StakingKeeper,
		&app.UpgradeKeeper,
		// Custom Modules
		&app.AccountedPoolKeeper,
		&app.AMMKeeper,
		&app.AssetProfileKeeper,
		&app.BurnerKeeper,
		&app.CommitmentKeeper,
		&app.EpochsKeeper,
		&app.EstakingKeeper,
		&app.LeverageLpKeeper,
		&app.MasterchefKeeper,
		&app.OracleKeeper,
		&app.ParameterKeeper,
		&app.PerpetualKeeper,
		&app.StablestakeKeeper,
		&app.TierKeeper,
		&app.TokenomicsKeeper,
		&app.TradeshieldKeeper,
		&app.TransferHookKeeper,
	); err != nil {
		return nil, err
	}

	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	if err := app.RegisterLegacyModules(); err != nil {
		panic(err)
	}

	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *SimApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *SimApp) SimulationManager() *module.SimulationManager {
	return nil
}

func (app *SimApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	key, _ := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	return key
}

func (app *SimApp) GetMemKey(memKey string) *storetypes.MemoryStoreKey {
	key, _ := app.UnsafeFindStoreKey(memKey).(*storetypes.MemoryStoreKey)
	return key
}

func (app *SimApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *SimApp) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}
