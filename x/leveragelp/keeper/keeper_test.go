package keeper_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
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
			price:   sdk.MustNewDecFromStr("6.0"),
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

func (suite *KeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) SetCurrentHeight(h int64) {
	suite.ctx = suite.ctx.WithBlockHeight(h)
}

func (suite *KeeperTestSuite) AddBlockTime(d time.Duration) {
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(d))
}

func (suite *KeeperTestSuite) EnableWhiteListing() {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.WhitelistingEnabled = true
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) DisableWhiteListing() {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.WhitelistingEnabled = false
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) SetMaxOpenPositions(value int64) {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.MaxOpenPositions = value
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) SetPoolThreshold(value sdk.Dec) {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.PoolOpenThreshold = value
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) SetSafetyFactor(value sdk.Dec) {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.SafetyFactor = value
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) EnablePool(poolId uint64) {
	pool, found := suite.app.LeveragelpKeeper.GetPool(suite.ctx, poolId)
	if !found {
		panic("pool not found")
	}
	pool.Enabled = true
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
}

func (suite *KeeperTestSuite) DisablePool(poolId uint64) {
	pool, found := suite.app.LeveragelpKeeper.GetPool(suite.ctx, poolId)
	if !found {
		panic("pool not found")
	}
	pool.Enabled = false
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func SetupCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range priceMap {
		oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
			Denom:   v.denom,
			Display: v.display,
			Decimal: 6,
		})
		oracle.SetPrice(ctx, oracletypes.Price{
			Asset:     v.display,
			Price:     v.price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(ctx.BlockTime().Unix()),
		})
	}
}

func AddCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper, denoms []string) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range denoms {
		oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
			Denom:   priceMap[v].denom,
			Display: priceMap[v].display,
			Decimal: 6,
		})
		oracle.SetPrice(ctx, oracletypes.Price{
			Asset:     priceMap[v].display,
			Price:     priceMap[v].price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(ctx.BlockTime().Unix()),
		})
	}
}

func RemovePrices(ctx sdk.Context, oracle oraclekeeper.Keeper, denoms []string) {
	for _, v := range denoms {
		oracle.RemoveAssetInfo(ctx, v)
		oracle.RemovePrice(ctx, priceMap[v].display, "elys", uint64(ctx.BlockTime().Unix()))
	}
}

func TestGetAllWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	leveragelp := app.LeveragelpKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	// Set whitelisted addresses
	leveragelp.WhitelistAddress(ctx, addr[0])
	leveragelp.WhitelistAddress(ctx, addr[1])

	// Get all whitelisted addresses
	whitelisted := leveragelp.GetAllWhitelistedAddress(ctx)

	// length should be 2
	require.Equal(t, len(whitelisted), 2)

	// If addr[0] is whitelisted
	require.Contains(t,
		whitelisted,
		addr[0],
	)

	// If addr[1] is whitelisted
	require.Contains(t,
		whitelisted,
		addr[1],
	)
}
