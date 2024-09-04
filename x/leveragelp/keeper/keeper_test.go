package keeper_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
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

func (suite *KeeperTestSuite) SetupCoinPrices(ctx sdk.Context) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range priceMap {
		suite.app.OracleKeeper.SetAssetInfo(ctx, oracletypes.AssetInfo{
			Denom:   v.denom,
			Display: v.display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(ctx, oracletypes.Price{
			Asset:     v.display,
			Price:     v.price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(ctx.BlockTime().Unix()),
		})
	}
}

func (suite *KeeperTestSuite) AddCoinPrices(ctx sdk.Context, denoms []string) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range denoms {
		suite.app.OracleKeeper.SetAssetInfo(ctx, oracletypes.AssetInfo{
			Denom:   priceMap[v].denom,
			Display: priceMap[v].display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(ctx, oracletypes.Price{
			Asset:     priceMap[v].display,
			Price:     priceMap[v].price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(ctx.BlockTime().Unix()),
		})
	}
}

func (suite *KeeperTestSuite) RemovePrices(ctx sdk.Context, denoms []string) {
	for _, v := range denoms {
		suite.app.OracleKeeper.RemoveAssetInfo(ctx, v)
		suite.app.OracleKeeper.RemovePrice(ctx, priceMap[v].display, "elys", uint64(ctx.BlockTime().Unix()))
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

func TestGetWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	leveragelp := app.LeveragelpKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	// Set whitelisted addresses
	leveragelp.WhitelistAddress(ctx, addr[0])
	leveragelp.WhitelistAddress(ctx, addr[1])

	// Get all whitelisted addresses
	whitelisted, _, _ := leveragelp.GetWhitelistedAddress(ctx, nil)

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

func (suite *KeeperTestSuite) TestEstimateSwapGivenOut() {
	app := suite.app
	ctx := suite.ctx

	leveragelp := app.LeveragelpKeeper


	testCases := []struct {
		name                 string
		tokenOutAmount       sdk.Coin
		tokenInDenom         string
		ammPool              ammtypes.Pool
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func()
	}{
		{
			"pool not found",
			sdk.NewCoin("uusdc", sdk.NewInt(100)),
			"uusdc",
			ammtypes.Pool{PoolId: 1},
			true,
			"pool 1 not found",
			func() {},
			func() {
			},
		},
		{
			"pool not enabled",
			sdk.NewCoin("uusdc", sdk.NewInt(100)),
			"uusdc",
			ammtypes.Pool{PoolId: 1},
			true,
			"pool 1: leveragelp disabled pool",
			func() {
				pool := types.NewPool(1)
				pool.Enabled = false
				leveragelp.SetPool(ctx, pool)
			},
			func() {
			},
		},
		{
			"amm pool not created",
			sdk.NewCoin("uusdc", sdk.NewInt(100).MulRaw(1000_000_000_000)),
			"uusdc",
			ammtypes.Pool{PoolId: 1},
			true,
			"invalid pool",
			func() {
				pool := types.NewPool(1)
				pool.Enabled = true
				leveragelp.SetPool(ctx, pool)
			},
			func() {
			},
		},
		{
			"amm pool not found in transient store ",
			sdk.NewCoin("uusdc", sdk.NewInt(100).MulRaw(1000_000_000_000)),
			"uusdc",
			ammtypes.Pool{PoolId: 1},
			true,
			"(uusdc) does not exist in the pool",
			func() {
				suite.SetupCoinPrices(suite.ctx)
				addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
				asset1 := ptypes.ATOM
				asset2 := ptypes.BaseCurrency
				initializeForClose(suite, addresses, asset1, asset2)
			},
			func() {
			},
		},
	}

	for _, tc := range testCases {
		tc.prerequisiteFunction()
		_, err := leveragelp.EstimateSwapGivenOut(ctx, tc.tokenOutAmount, tc.tokenInDenom, tc.ammPool)
		if tc.expectErr {
			suite.Require().Error(err)
			suite.Require().Contains(err.Error(), tc.expectErrMsg)
		} else {
			suite.Require().NoError(err)
		}
	}
}

func (suite *KeeperTestSuite) TestCalculatePoolHealth() {
	app := suite.app
	ctx := suite.ctx

	leveragelp := app.LeveragelpKeeper

	leveragelpAmount := sdk.NewInt(10)
	pool := &types.Pool{
		AmmPoolId:         1,
		LeveragedLpAmount: leveragelpAmount,
	}
	ammPool := ammtypes.Pool{PoolId: 1}
	totalShares := sdk.NewInt(100)

	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		expectedValue        sdk.Dec
	}{
		{
			"amm pool not found",
			func() {},
			sdk.ZeroDec(),
		},
		{
			"amm pool shares is  0",
			func() {
				_ = app.AmmKeeper.SetPool(ctx, ammPool)
			},
			sdk.OneDec(),
		},
		{
			"success",
			func() {
				ammPool.TotalShares = sdk.NewCoin("shares", totalShares)
				_ = app.AmmKeeper.SetPool(ctx, ammPool)
			},
			(totalShares.Sub(leveragelpAmount)).ToLegacyDec().QuoInt(totalShares),
		},
	}

	for _, tc := range testCases {
		tc.prerequisiteFunction()
		out := leveragelp.CalculatePoolHealth(ctx, pool)
		suite.Require().Equal(tc.expectedValue, out)
	}
}
