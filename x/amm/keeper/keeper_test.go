package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/suite"
)

type assetPriceInfo struct {
	denom   string
	display string
	price   sdk.Dec
}

const (
	initChain = true
)

var (
	priceMap = map[string]assetPriceInfo{
		"uusdc": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   sdk.OneDec(),
		},
		"uusdt": {
			denom:   "uusdt",
			display: "USDT",
			price:   sdk.OneDec(),
		},
		"uelys": {
			denom:   ptypes.Elys,
			display: "ELYS",
			price:   sdk.MustNewDecFromStr("3.0"),
		},
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   sdk.MustNewDecFromStr("1.0"),
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
	app := simapp.InitElysTestApp(initChain)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain, tmproto.Header{})
	suite.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
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
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDA",
		Price:     sdk.NewDec(1),
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
				{Token: sdk.NewCoin("denom1", sdk.NewInt(1000)), Weight: sdk.OneInt()},
				{Token: sdk.NewCoin("denom2", sdk.NewInt(1000)), Weight: sdk.OneInt()},
			},
			TotalWeight: sdk.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(1), types.OneShare),
		},
		{
			PoolId:  2,
			Address: types.NewPoolAddress(uint64(2)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("uusdc", sdk.NewInt(1000)), Weight: sdk.OneInt()},
				{Token: sdk.NewCoin("denom1", sdk.NewInt(1000)), Weight: sdk.OneInt()},
			},
			TotalWeight: sdk.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(2), types.OneShare),
		},
		{
			PoolId:  3,
			Address: types.NewPoolAddress(uint64(3)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("uusdc", sdk.NewInt(1000)), Weight: sdk.OneInt()},
				{Token: sdk.NewCoin("denom3", sdk.NewInt(1000)), Weight: sdk.OneInt()},
			},
			TotalWeight: sdk.NewInt(2),
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
