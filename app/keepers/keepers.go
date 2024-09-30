package keepers

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmmodulekeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	wasmbindingsclient "github.com/elys-network/elys/wasmbindings/client"
	accountedpoolmodulekeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	accountedpoolmoduletypes "github.com/elys-network/elys/x/accountedpool/types"
	ammmodulekeeper "github.com/elys-network/elys/x/amm/keeper"
	ammmoduletypes "github.com/elys-network/elys/x/amm/types"
	assetprofilemodulekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	assetprofilemoduletypes "github.com/elys-network/elys/x/assetprofile/types"
	burnermodulekeeper "github.com/elys-network/elys/x/burner/keeper"
	burnermoduletypes "github.com/elys-network/elys/x/burner/types"
	clockmodulekeeper "github.com/elys-network/elys/x/clock/keeper"
	clockmoduletypes "github.com/elys-network/elys/x/clock/types"
	commitmentmodulekeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmentmoduletypes "github.com/elys-network/elys/x/commitment/types"
	epochsmodulekeeper "github.com/elys-network/elys/x/epochs/keeper"
	epochsmoduletypes "github.com/elys-network/elys/x/epochs/types"
	estakingmodulekeeper "github.com/elys-network/elys/x/estaking/keeper"
	estakingmoduletypes "github.com/elys-network/elys/x/estaking/types"
	incentivemodulekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivemoduletypes "github.com/elys-network/elys/x/incentive/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	masterchefmodulekeeper "github.com/elys-network/elys/x/masterchef/keeper"
	masterchefmoduletypes "github.com/elys-network/elys/x/masterchef/types"
	oraclemodule "github.com/elys-network/elys/x/oracle"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parametermodulekeeper "github.com/elys-network/elys/x/parameter/keeper"
	parametermoduletypes "github.com/elys-network/elys/x/parameter/types"
	perpetualmodulekeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualmoduletypes "github.com/elys-network/elys/x/perpetual/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	tiermodulekeeper "github.com/elys-network/elys/x/tier/keeper"
	tiermoduletypes "github.com/elys-network/elys/x/tier/types"
	tokenomicsmodulekeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicsmoduletypes "github.com/elys-network/elys/x/tokenomics/types"
	tradeshieldmodulekeeper "github.com/elys-network/elys/x/tradeshield/keeper"
	tradeshieldmoduletypes "github.com/elys-network/elys/x/tradeshield/types"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
	transferhooktypes "github.com/elys-network/elys/x/transferhook/types"
	"os"
	// this line is used by starport scaffolding # stargate/app/moduleImport
)

type AppKeepers struct {
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        *govkeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper

	// IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCKeeper             *ibckeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	ICAControllerKeeper   icacontrollerkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	WasmKeeper            wasmmodulekeeper.Keeper
	GroupKeeper           groupkeeper.Keeper

	ICAModule      ica.AppModule
	IBCFeeKeeper   ibcfeekeeper.Keeper
	IBCHooksKeeper *ibchookskeeper.Keeper
	TransferModule transfer.AppModule

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCFeeKeeper        capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper          capabilitykeeper.ScopedKeeper

	EpochsKeeper        *epochsmodulekeeper.Keeper
	AssetprofileKeeper  assetprofilemodulekeeper.Keeper
	ScopedOracleKeeper  capabilitykeeper.ScopedKeeper
	OracleKeeper        oraclekeeper.Keeper
	CommitmentKeeper    *commitmentmodulekeeper.Keeper
	TokenomicsKeeper    tokenomicsmodulekeeper.Keeper
	IncentiveKeeper     incentivemodulekeeper.Keeper
	BurnerKeeper        burnermodulekeeper.Keeper
	AmmKeeper           *ammmodulekeeper.Keeper
	ParameterKeeper     parametermodulekeeper.Keeper
	PerpetualKeeper     perpetualmodulekeeper.Keeper
	TransferhookKeeper  transferhookkeeper.Keeper
	ContractKeeper      *wasmmodulekeeper.PermissionedKeeper
	ClockKeeper         clockmodulekeeper.Keeper
	AccountedPoolKeeper accountedpoolmodulekeeper.Keeper
	StablestakeKeeper   *stablestakekeeper.Keeper
	LeveragelpKeeper    *leveragelpmodulekeeper.Keeper
	MasterchefKeeper    masterchefmodulekeeper.Keeper
	EstakingKeeper      *estakingmodulekeeper.Keeper
	TierKeeper          tiermodulekeeper.Keeper
	TradeshieldKeeper   tradeshieldmodulekeeper.Keeper

	Ics20WasmHooks   *ibchooks.WasmHooks
	HooksICS4Wrapper ibchooks.ICS4Middleware
}

func NewAppKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	modAccAddrs map[string]bool,
	blockedAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	logger log.Logger,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	enabledProposals []wasmtypes.ProposalType,
	AccountAddressPrefix string,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	/*
		configure state listening capabilities using AppOptions
		we are doing nothing with the returned streamingServices and waitGroup in this case
	*/
	// load state streaming if enabled

	if err := bApp.RegisterStreamingServices(appOpts, appKeepers.keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(appKeepers.ConsensusParamsKeeper.ParamsStore)

	appKeepers.ParameterKeeper = *parametermodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[parametermoduletypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	appKeepers.ScopedICAControllerKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	//appKeepers.ScopedICSproviderkeeper = appKeepers.CapabilityKeeper.ScopeToModule(providertypes.ModuleName)
	appKeepers.ScopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[crisistypes.StoreKey]),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.AccountKeeper.AddressCodec(),
	)

	// Add normal keepers
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[banktypes.StoreKey]),
		appKeepers.AccountKeeper,
		blockedAddress,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]),
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]),
		appKeepers.AccountKeeper,
	)

	appKeepers.AssetprofileKeeper = *assetprofilemodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[assetprofilemoduletypes.StoreKey]),
		&appKeepers.TransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	commitmentKeeper := commitmentmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[commitmentmoduletypes.StoreKey]),

		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.AssetprofileKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.TokenomicsKeeper = *tokenomicsmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[tokenomicsmoduletypes.StoreKey]),
		appKeepers.CommitmentKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.EstakingKeeper = estakingmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[estakingmoduletypes.StoreKey]),
		appKeepers.ParameterKeeper,
		appKeepers.StakingKeeper,
		appKeepers.CommitmentKeeper,
		&appKeepers.DistrKeeper,
		appKeepers.AssetprofileKeeper,
		appKeepers.TokenomicsKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.EstakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			// insert staking hooks receivers here
			appKeepers.SlashingKeeper.Hooks(),
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.EstakingKeeper.StakingHooks(),
			appKeepers.TierKeeper.StakingHooks(),
		),
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]),
		appKeepers.EstakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	groupConfig := group.DefaultConfig()
	appKeepers.GroupKeeper = groupkeeper.NewKeeper(
		appKeepers.keys[group.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
		groupConfig,
	)

	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(appKeepers.keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBC Fee Module keeper
	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper,
	)
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCFeeKeeper, // ISC4 Wrapper: PFM Router middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	//transferIBCModule := transfer.NewIBCModule(appKeepers.TransferKeeper)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.ScopedICAHostKeeper,
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// required since ibc-go v7.5.0
	appKeepers.ICAHostKeeper.WithQueryRouter(bApp.GRPCQueryRouter())

	appKeepers.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedICAControllerKeeper,
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]),
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
		appKeepers.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	appKeepers.ScopedOracleKeeper = appKeepers.CapabilityKeeper.ScopeToModule(oracletypes.ModuleName)
	appKeepers.OracleKeeper = *oraclekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[oracletypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedOracleKeeper,
	)

	oracleIBCModule := oraclemodule.NewIBCModule(appKeepers.OracleKeeper)

	appKeepers.EpochsKeeper = epochsmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[epochsmoduletypes.StoreKey]),
	)

	appKeepers.EpochsKeeper = appKeepers.EpochsKeeper.SetHooks(
		epochsmodulekeeper.NewMultiEpochHooks(
			// insert epoch hooks receivers here
			appKeepers.OracleKeeper.Hooks(),
			appKeepers.CommitmentKeeper.Hooks(),
			appKeepers.BurnerKeeper.Hooks(),
			appKeepers.PerpetualKeeper.Hooks(),
		),
	)

	appKeepers.AccountedPoolKeeper = *accountedpoolmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[accountedpoolmoduletypes.StoreKey]),
		appKeepers.BankKeeper,
	)

	appKeepers.AmmKeeper = ammmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[ammmoduletypes.StoreKey]),
		appKeepers.tkeys[ammmoduletypes.TStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		&appKeepers.ParameterKeeper,
		appKeepers.BankKeeper,
		appKeepers.AccountKeeper,
		appKeepers.OracleKeeper,
		appKeepers.CommitmentKeeper,
		appKeepers.AssetprofileKeeper,
		appKeepers.AccountedPoolKeeper,
	)

	appKeepers.AmmKeeper.SetHooks(
		ammmoduletypes.NewMultiAmmHooks(
			// insert amm hooks receivers here
			appKeepers.PerpetualKeeper.AmmHooks(),
			appKeepers.LeveragelpKeeper.AmmHooks(),
			appKeepers.MasterchefKeeper.AmmHooks(),
			appKeepers.TierKeeper.AmmHooks(),
		),
	)

	appKeepers.StablestakeKeeper = stablestakekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[stablestaketypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.BankKeeper,
		appKeepers.CommitmentKeeper,
		appKeepers.AssetprofileKeeper,
	)
	appKeepers.StablestakeKeeper.SetHooks(stablestakekeeper.NewMultiStableStakeHooks(
		appKeepers.MasterchefKeeper.StableStakeHooks(),
		appKeepers.TierKeeper.StableStakeHooks(),
	))

	appKeepers.CommitmentKeeper.SetHooks(
		commitmentmodulekeeper.NewMultiCommitmentHooks(
			appKeepers.EstakingKeeper.CommitmentHooks(),
		),
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.CommitmentKeeper,
		appKeepers.EstakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.MasterchefKeeper = *masterchefmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[masterchefmoduletypes.StoreKey]),
		appKeepers.ParameterKeeper,
		appKeepers.CommitmentKeeper,
		appKeepers.AmmKeeper,
		appKeepers.OracleKeeper,
		appKeepers.AssetprofileKeeper,
		appKeepers.AccountedPoolKeeper,
		appKeepers.StablestakeKeeper,
		appKeepers.TokenomicsKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.IncentiveKeeper = *incentivemodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[incentivemoduletypes.StoreKey]),
		appKeepers.ParameterKeeper,
		commitmentKeeper,
		appKeepers.StakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.AmmKeeper,
		appKeepers.OracleKeeper,
		appKeepers.AssetprofileKeeper,
		appKeepers.AccountedPoolKeeper,
		appKeepers.StablestakeKeeper,
		appKeepers.TokenomicsKeeper,
		&appKeepers.MasterchefKeeper,
		appKeepers.EstakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.BurnerKeeper = *burnermodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[burnermoduletypes.StoreKey]),
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	wasmDir := homePath
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	bankKeeper := appKeepers.BankKeeper.(bankkeeper.BaseKeeper)
	wasmOpts = append(
		wasmbindingsclient.RegisterCustomPlugins(
			&appKeepers.AccountedPoolKeeper,
			appKeepers.AmmKeeper,
			&appKeepers.AssetprofileKeeper,
			&appKeepers.AccountKeeper,
			&bankKeeper,
			&appKeepers.BurnerKeeper,
			&appKeepers.ClockKeeper,
			appKeepers.CommitmentKeeper,
			appKeepers.EpochsKeeper,
			&appKeepers.IncentiveKeeper,
			appKeepers.LeveragelpKeeper,
			&appKeepers.PerpetualKeeper,
			&appKeepers.OracleKeeper,
			&appKeepers.ParameterKeeper,
			appKeepers.StablestakeKeeper,
			appKeepers.StakingKeeper,
			&appKeepers.TokenomicsKeeper,
			&appKeepers.TransferhookKeeper,
			&appKeepers.MasterchefKeeper,
			appKeepers.EstakingKeeper,
			&appKeepers.TierKeeper,
		),
		wasmOpts...,
	)

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[wasmtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmkeeper.BuiltInCapabilities(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	appKeepers.TransferhookKeeper = *transferhookkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[transferhooktypes.StoreKey]),
		*appKeepers.AmmKeeper)
	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		appKeepers.keys[ibchookstypes.StoreKey],
	)
	appKeepers.IBCHooksKeeper = &hooksKeeper

	wasmHooks := ibchooks.NewWasmHooks(appKeepers.IBCHooksKeeper, &appKeepers.WasmKeeper, AccountAddressPrefix) // The contract keeper needs to be set later
	appKeepers.Ics20WasmHooks = &wasmHooks

	// provider depends on gov, so gov must be registered first
	govConfig := govtypes.DefaultConfig()
	// set the MaxMetadataLen for proposals to the same value as it was pre-sdk v0.47.x
	govConfig.MaxMetadataLen = 10200
	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.DistrKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		//AddRoute(upgradetypes.RouterKey, upgradetypes.NewSoftwareUpgradeProposal(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))

	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	appKeepers.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.Ics20WasmHooks,
	)

	// set the contract keeper for the Ics20WasmHooks
	appKeepers.ContractKeeper = wasmmodulekeeper.NewDefaultPermissionKeeper(appKeepers.WasmKeeper)
	appKeepers.Ics20WasmHooks.ContractKeeper = &appKeepers.WasmKeeper

	appKeepers.PerpetualKeeper = *perpetualmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[perpetualmoduletypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.AmmKeeper,
		appKeepers.BankKeeper,
		appKeepers.OracleKeeper,
		appKeepers.AssetprofileKeeper,
		&appKeepers.ParameterKeeper,
	)

	appKeepers.ClockKeeper = *clockmodulekeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[clockmoduletypes.StoreKey]),
		appCodec,
		*appKeepers.ContractKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.LeveragelpKeeper = leveragelpmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[leveragelpmoduletypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.AmmKeeper,
		appKeepers.BankKeeper,
		appKeepers.OracleKeeper,
		appKeepers.StablestakeKeeper,
		appKeepers.CommitmentKeeper,
		appKeepers.AssetprofileKeeper,
		appKeepers.MasterchefKeeper,
	)

	appKeepers.LeveragelpKeeper.SetHooks(leveragelpmoduletypes.NewMultiLeverageLpHooks(
		appKeepers.TierKeeper.LeverageLpHooks(),
	))

	appKeepers.TierKeeper = *tiermodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[tiermoduletypes.StoreKey]),
		appKeepers.BankKeeper,
		appKeepers.OracleKeeper,
		appKeepers.AssetprofileKeeper,
		appKeepers.AmmKeeper,
		appKeepers.EstakingKeeper,
		appKeepers.MasterchefKeeper,
		appKeepers.CommitmentKeeper,
		appKeepers.StakingKeeper,
		appKeepers.PerpetualKeeper,
		appKeepers.LeveragelpKeeper,
		appKeepers.StablestakeKeeper,
	)

	appKeepers.TradeshieldKeeper = *tradeshieldmodulekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[tradeshieldmoduletypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	/**** IBC Routing ****/

	// Sealing prevents other modules from creating scoped sub-keepers
	appKeepers.CapabilityKeeper.Seal()

	appKeepers.ICAModule = ica.NewAppModule(&appKeepers.ICAControllerKeeper, &appKeepers.ICAHostKeeper)
	appKeepers.TransferModule = transfer.NewAppModule(appKeepers.TransferKeeper)

	// Create Transfer Stack (from bottom to top of stack)
	// - core IBC
	// - ibcfee
	// - ratelimit
	// - pfm
	// - provider
	// - transfer
	//
	// This is how transfer stack will work in the end:
	// * RecvPacket -> IBC core -> Fee -> RateLimit -> PFM -> Provider -> Transfer (AddRoute)
	// * SendPacket -> Transfer -> Provider -> PFM -> RateLimit -> Fee -> IBC core (ICS4Wrapper)

	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, appKeepers.IBCFeeKeeper)

	// Create ICAHost Stack
	var icaHostStack porttypes.IBCModule = icahost.NewIBCModule(appKeepers.ICAHostKeeper)

	// Create Interchain Accounts Controller Stack
	var icaControllerStack porttypes.IBCModule = icacontroller.NewIBCMiddleware(nil, appKeepers.ICAControllerKeeper)

	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, appKeepers.IBCFeeKeeper)

	// Create IBC Router & seal
	ibcRouter := porttypes.NewRouter().
		AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasmStack).
		AddRoute(oracletypes.ModuleName, oracleIBCModule)

	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	return appKeepers

}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, ok := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	if !ok {
		panic("couldn't load subspace for module: " + moduleName)
	}
	return subspace
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)

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
