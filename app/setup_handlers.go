package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	accountedpooltypes "github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	burnertypes "github.com/elys-network/elys/x/burner/types"
	clocktypes "github.com/elys-network/elys/x/clock/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	liquidityprovidertypes "github.com/elys-network/elys/x/liquidityprovider/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parametertypes "github.com/elys-network/elys/x/parameter/types"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
)

func SetupHandlers(app *ElysApp) {
	setUpgradeHandler(app)

	loadUpgradeStore(app)
}

func setUpgradeHandler(app *ElysApp) {
	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		subspace := subspace

		app.Logger().Info("Setting up upgrade handler for " + subspace.Name())

		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable() //nolint:staticcheck
		case minttypes.ModuleName:
			keyTable = minttypes.ParamKeyTable() //nolint:staticcheck
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable() //nolint:staticcheck
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
		case ammtypes.ModuleName:
			keyTable = ammtypes.ParamKeyTable() //nolint:staticcheck
		case assetprofiletypes.ModuleName:
			keyTable = assetprofiletypes.ParamKeyTable() //nolint:staticcheck
		case burnertypes.ModuleName:
			keyTable = burnertypes.ParamKeyTable() //nolint:staticcheck
		case commitmenttypes.ModuleName:
			keyTable = commitmenttypes.ParamKeyTable() //nolint:staticcheck
		case incentivetypes.ModuleName:
			keyTable = incentivetypes.ParamKeyTable() //nolint:staticcheck
		case liquidityprovidertypes.ModuleName:
			keyTable = liquidityprovidertypes.ParamKeyTable() //nolint:staticcheck
		case margintypes.ModuleName:
			keyTable = margintypes.ParamKeyTable() //nolint:staticcheck
		case oracletypes.ModuleName:
			keyTable = oracletypes.ParamKeyTable() //nolint:staticcheck
		case parametertypes.ModuleName:
			keyTable = parametertypes.ParamKeyTable() //nolint:staticcheck
		case tokenomicstypes.ModuleName:
			keyTable = tokenomicstypes.ParamKeyTable() //nolint:staticcheck
		case accountedpooltypes.ModuleName:
			keyTable = accountedpooltypes.ParamKeyTable() //nolint:staticcheck
		case clocktypes.ModuleName:
			keyTable = clocktypes.ParamKeyTable() //nolint:staticcheck
		}

		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())

	app.UpgradeKeeper.SetUpgradeHandler(
		version.Version,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			app.Logger().Info("Running upgrade handler for " + version.Version)

			// Migrate Tendermint consensus parameters from x/params module to a
			// dedicated x/consensus module.
			baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

			return app.mm.RunMigrations(ctx, app.configurator, vm)
		},
	)
}

func loadUpgradeStore(app *ElysApp) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if shouldLoadUpgradeStore(app, upgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				clocktypes.ModuleName,
			},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func shouldLoadUpgradeStore(app *ElysApp, upgradeInfo upgradetypes.Plan) bool {
	return upgradeInfo.Name == version.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height)
}
