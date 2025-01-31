package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type assetPriceInfo struct {
	denom   string
	display string
	price   math.LegacyDec
}

const (
	initChain = true
)

var (
	priceMap = map[string]assetPriceInfo{
		"uusdc": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   math.LegacyOneDec(),
		},
		"uusdt": {
			denom:   "uusdt",
			display: "USDT",
			price:   math.LegacyOneDec(),
		},
		"uelys": {
			denom:   ptypes.Elys,
			display: "ELYS",
			price:   math.LegacyMustNewDecFromStr("3.0"),
		},
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   math.LegacyMustNewDecFromStr("6.0"),
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
	t := suite.T()
	app := simapp.InitElysTestApp(initChain, t)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app

	suite.SetupAssetProfile(suite.ctx)
	suite.SetStakingParam(suite.ctx)
	suite.SetStableStakeParam(suite.ctx)
	suite.SetLeverageParam(suite.ctx)

	suite.app.StablestakeKeeper.SetPool(suite.ctx, stablestaketypes.Pool{
		RedemptionRate:       math.LegacyOneDec(),
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		TotalValue:           math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		PoolId:               1,
		DepositDenom:         ptypes.BaseCurrency,
	})
}

func (suite *KeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) SetCurrentHeight(h int64) {
	suite.ctx = suite.ctx.WithBlockHeight(h)
}

func (suite *KeeperTestSuite) AddBlockTime(d time.Duration) {
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(d)).WithBlockHeight(int64(d.Seconds() / 4))
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

func (suite *KeeperTestSuite) SetPoolThreshold(value math.LegacyDec) {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.PoolOpenThreshold = value
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) SetSafetyFactor(value math.LegacyDec) {
	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.SafetyFactor = value
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	if err != nil {
		panic(err)
	}
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
func (suite *KeeperTestSuite) SetLeverageParam(ctx sdk.Context) error {

	params := &types.DefaultGenesis().Params

	suite.app.LeveragelpKeeper.SetParams(ctx, params)
	return nil
}

func (suite *KeeperTestSuite) SetStakingParam(ctx sdk.Context) error {
	return suite.app.StakingKeeper.SetParams(ctx, stakingtypes.Params{
		UnbondingTime:     1209600,
		MaxValidators:     60,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		BondDenom:         "uelys",
		MinCommissionRate: math.LegacyNewDec(0),
	})
}

func (suite *KeeperTestSuite) SetStableStakeParam(ctx sdk.Context) error {

	params := stablestaketypes.DefaultParams()
	suite.app.StablestakeKeeper.SetParams(ctx, params)
	return nil
}

func (suite *KeeperTestSuite) SetupAssetProfile(ctx sdk.Context) {

	suite.app.AssetprofileKeeper.SetEntry(ctx, atypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "uusdc",
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

func TestGetAllWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	leveragelp := app.LeveragelpKeeper

	simapp.SetStakingParam(app, ctx)
	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

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
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	leveragelp := app.LeveragelpKeeper
	simapp.SetStakingParam(app, ctx)

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

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
			sdk.NewCoin("uusdc", math.NewInt(100)),
			"uusdc",
			ammtypes.Pool{PoolId: 1},
			true,
			"pool 1 not found",
			func() {},
			func() {
			},
		},
		{
			"amm pool not found in transient store ",
			sdk.NewCoin("uusdc", math.NewInt(100).MulRaw(1000_000_000_000)),
			"uusdc",
			ammtypes.Pool{PoolId: 1},
			true,
			"(uusdc) does not exist in the pool",
			func() {
				suite.SetupCoinPrices(suite.ctx)
				addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, math.NewInt(1000000))
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

	leveragelpAmount := math.NewInt(10)
	pool := &types.Pool{
		AmmPoolId:         1,
		LeveragedLpAmount: leveragelpAmount,
	}
	ammPool := ammtypes.Pool{PoolId: 1, Address: ammtypes.NewPoolAddress(uint64(1)).String()}
	totalShares := math.NewInt(100)

	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		expectedValue        math.LegacyDec
	}{
		{
			"amm pool not found",
			func() {},
			math.LegacyZeroDec(),
		},
		{
			"amm pool shares is  0",
			func() {
				app.AmmKeeper.SetPool(ctx, ammPool)
			},
			math.LegacyOneDec(),
		},
		{
			"success",
			func() {
				ammPool.TotalShares = sdk.NewCoin("shares", totalShares)
				app.AmmKeeper.SetPool(ctx, ammPool)
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
