package keeper_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/suite"
)

type assetPriceInfo struct {
	denom   string
	display string
	price   sdkmath.LegacyDec
}

const (
	initChain = true
)

var (
	priceMap = map[string]assetPriceInfo{
		"uusdc": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   sdkmath.LegacyOneDec(),
		},
		"uusdt": {
			denom:   "uusdt",
			display: "USDT",
			price:   sdkmath.LegacyOneDec(),
		},
		"USDC": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   sdkmath.LegacyOneDec(),
		},
		"uelys": {
			denom:   ptypes.Elys,
			display: "ELYS",
			price:   sdkmath.LegacyMustNewDecFromStr("3.0"),
		},
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   sdkmath.LegacyMustNewDecFromStr("1.0"),
		},
	}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *KeeperTestSuite) SetupTest() {
	//t.Parallel()
	app := simapp.InitElysTestApp(initChain, suite.Suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetAmmParams() {
	suite.app.AmmKeeper.SetParams(suite.ctx, types.Params{
		PoolCreationFee:             sdkmath.NewInt(10_000_000),
		SlippageTrackDuration:       604800,
		BaseAssets:                  []string{ptypes.BaseCurrency},
		AllowedPoolCreators:         []string{authtypes.NewModuleAddress(govtypes.ModuleName).String()},
		WeightBreakingFeeExponent:   sdkmath.LegacyMustNewDecFromStr("2.5"),
		WeightBreakingFeeMultiplier: sdkmath.LegacyMustNewDecFromStr("0.0005"),
		WeightBreakingFeePortion:    sdkmath.LegacyMustNewDecFromStr("0.5"),
		WeightRecoveryFeePortion:    sdkmath.LegacyMustNewDecFromStr("0.1"),
		ThresholdWeightDifference:   sdkmath.LegacyMustNewDecFromStr("0.3"),
	})
}

func (suite *KeeperTestSuite) SetupAssetProfile() {
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, atypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
		Path:                     "transfer/channel-12",
		IbcChannelId:             "channel-12",
		IbcCounterpartyChannelId: "channel-19",
		DisplayName:              "USDC",
		DisplaySymbol:            "uUSDC",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "uUSDC",
		TransferLimit:            "",
		Permissions:              []string{},
		UnitDenom:                "uusdc",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	})
}

func (suite *KeeperTestSuite) SetupStableCoinPrices() {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusda",
		Display: "USDA",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdkmath.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdkmath.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDA",
		Price:     sdkmath.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
}

func (suite *KeeperTestSuite) SetupCoinPrices() {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range priceMap {
		suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
			Denom:   v.denom,
			Display: v.display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
			Asset:     v.display,
			Price:     v.price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(suite.ctx.BlockTime().Unix()),
		})
	}
}

func SetupMockPools(k *keeper.Keeper, ctx sdk.Context) {
	// Create and set mock pools
	pools := []types.Pool{
		{
			PoolId:  1,
			Address: types.NewPoolAddress(uint64(1)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("denom1", sdkmath.NewInt(1000)), Weight: sdkmath.OneInt()},
				{Token: sdk.NewCoin("denom2", sdkmath.NewInt(1000)), Weight: sdkmath.OneInt()},
			},
			TotalWeight: sdkmath.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(1), types.OneShare),
		},
		{
			PoolId:  2,
			Address: types.NewPoolAddress(uint64(2)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("uusdc", sdkmath.NewInt(1000)), Weight: sdkmath.OneInt()},
				{Token: sdk.NewCoin("denom1", sdkmath.NewInt(1000)), Weight: sdkmath.OneInt()},
			},
			TotalWeight: sdkmath.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(2), types.OneShare),
		},
		{
			PoolId:  3,
			Address: types.NewPoolAddress(uint64(3)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("uusdc", sdkmath.NewInt(1000)), Weight: sdkmath.OneInt()},
				{Token: sdk.NewCoin("denom3", sdkmath.NewInt(1000)), Weight: sdkmath.OneInt()},
			},
			TotalWeight: sdkmath.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(3), types.OneShare),
		},
	}

	for _, pool := range pools {
		k.SetPool(ctx, pool)
	}
}
