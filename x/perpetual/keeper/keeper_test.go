package keeper_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	oraclekeeper "github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
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

var oracleProvider = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

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
			price:   math.LegacyMustNewDecFromStr("5.0"),
		},
	}
)

type PerpetualKeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *PerpetualKeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain, suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(PerpetualKeeperTestSuite))
}

func (suite *PerpetualKeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *PerpetualKeeperTestSuite) ResetAndSetSuite(addr []sdk.AccAddress, useOracle bool, baseTokenAmount, assetAmount math.Int) (ammtypes.Pool, types.Pool) {
	suite.ResetSuite()
	suite.SetupCoinPrices()
	suite.AddAccounts(len(addr), addr)
	poolCreator := addr[0]
	ammPool := suite.CreateNewAmmPool(poolCreator, useOracle, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, baseTokenAmount, assetAmount)
	pool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("10.5"))
	suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	params.BorrowInterestRateMin = math.LegacyMustNewDecFromStr("0.12")
	params.MaximumLongTakeProfitPriceRatio = math.LegacyMustNewDecFromStr("11.000000000000000000")
	params.MinimumLongTakeProfitPriceRatio = math.LegacyMustNewDecFromStr("1.020000000000000000")
	params.MaximumShortTakeProfitPriceRatio = math.LegacyMustNewDecFromStr("0.980000000000000000")
	err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)

	return ammPool, pool
}

func (suite *PerpetualKeeperTestSuite) SetCurrentHeight(h int64) {
	suite.ctx = suite.ctx.WithBlockHeight(h)
}

func (suite *PerpetualKeeperTestSuite) AddBlockTime(d time.Duration) {
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(d))
}

func (suite *PerpetualKeeperTestSuite) SetupCoinPrices() {
	// prices set for USDT and USDC
	provider := oracleProvider

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

func (suite *PerpetualKeeperTestSuite) AddCoinPrices(denoms []string) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range denoms {
		suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
			Denom:   priceMap[v].denom,
			Display: priceMap[v].display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
			Asset:     priceMap[v].display,
			Price:     priceMap[v].price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(suite.ctx.BlockTime().Unix()),
		})
	}
}

func (suite *PerpetualKeeperTestSuite) RemovePrices(ctx sdk.Context, denoms []string) {
	for _, v := range denoms {
		suite.app.OracleKeeper.RemoveAssetInfo(ctx, v)
		suite.app.OracleKeeper.RemovePrice(ctx, priceMap[v].display, "elys", uint64(ctx.BlockTime().Unix()))
	}
}

func (suite *PerpetualKeeperTestSuite) GetAccountIssueAmount() math.Int {
	return math.NewInt(10_000_000_000_000)
}

func (suite *PerpetualKeeperTestSuite) AddAccounts(n int, given []sdk.AccAddress) []sdk.AccAddress {
	issueAmount := suite.GetAccountIssueAmount()
	var addresses []sdk.AccAddress
	if n > len(given) {
		addresses = simapp.AddTestAddrs(suite.app, suite.ctx, n-len(given), issueAmount)
		addresses = append(addresses, given...)
	} else {
		addresses = given
	}
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount),
			sdk.NewCoin(ptypes.Elys, issueAmount),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
		if err != nil {
			panic(err)
		}
	}
	return addresses
}

func (suite *PerpetualKeeperTestSuite) CreateNewAmmPool(creator sdk.AccAddress, useOracle bool, swapFee, exitFee math.LegacyDec, asset2 string, baseTokenAmount, assetAmount math.Int) ammtypes.Pool {
	poolAssets := []ammtypes.PoolAsset{
		{
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, baseTokenAmount),
			Weight:                 math.NewInt(10),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Token:                  sdk.NewCoin(asset2, assetAmount),
			Weight:                 math.NewInt(10),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
	}
	sort.Slice(poolAssets, func(i, j int) bool {
		return strings.Compare(poolAssets[i].Token.Denom, poolAssets[j].Token.Denom) <= 0
	})
	poolParams := ammtypes.PoolParams{
		UseOracle: useOracle,
		SwapFee:   swapFee,
		FeeDenom:  ptypes.BaseCurrency,
	}

	createPoolMsg := &ammtypes.MsgCreatePool{
		Sender:     creator.String(),
		PoolParams: poolParams,
		PoolAssets: poolAssets,
	}

	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, createPoolMsg)
	suite.Require().NoError(err)
	ammPool, _ := suite.app.AmmKeeper.GetPool(suite.ctx, poolId)

	return ammPool
}

func (suite *PerpetualKeeperTestSuite) SetPerpetualPool(poolId uint64) (types.Pool, sdk.AccAddress, ammtypes.Pool) {
	ctx := suite.ctx
	k := suite.app.PerpetualKeeper
	//prices
	suite.SetupCoinPrices()
	//accounts
	accounts := suite.AddAccounts(2, nil)
	poolCreator := accounts[0]

	amount := math.NewInt(100000000000)

	ammPool := suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	suite.Require().NoError(err)

	pool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("11"))
	k.SetPool(ctx, pool)

	return pool, poolCreator, ammPool
}

func (suite *PerpetualKeeperTestSuite) AddLiquidity(ammPool ammtypes.Pool, provider sdk.AccAddress, tokensIn sdk.Coins) {
	numShares, _, err := ammPool.CalcJoinPoolNoSwapShares(tokensIn)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoins(suite.ctx, provider, sdk.MustAccAddressFromBech32(ammPool.GetAddress()), tokensIn)
	suite.Require().NoError(err)
	err = suite.app.AmmKeeper.MintPoolShareToAccount(suite.ctx, ammPool, provider, numShares)
	suite.Require().NoError(err)
	//IncreaseDenomLiquidity
	err = ammPool.IncreaseLiquidity(numShares, tokensIn)
	suite.Require().NoError(err)
	suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)

}
func TestSetGetMTP(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	simapp.SetStakingParam(app, ctx)
	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	for i := 0; i < 2; i++ {
		mtp := types.MTP{
			Address:                       addr[i].String(),
			CollateralAsset:               ptypes.BaseCurrency,
			CustodyAsset:                  "ATOM",
			Collateral:                    math.NewInt(0),
			Liabilities:                   math.NewInt(0),
			BorrowInterestUnpaidLiability: math.NewInt(0),
			BorrowInterestPaidCustody:     math.NewInt(0),
			TakeProfitBorrowFactor:        math.LegacyZeroDec(),
			Custody:                       math.NewInt(0),
			MtpHealth:                     math.LegacyNewDec(0),
			Position:                      types.Position_LONG,
			Id:                            0,
		}
		err := perpetual.SetMTP(ctx, &mtp)
		require.NoError(t, err)
	}

	mtpCount := perpetual.GetMTPCount(ctx)
	require.Equal(t, mtpCount, (uint64)(2))
}

func TestGetAllWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)

	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	// Set whitelisted addresses
	perpetual.WhitelistAddress(ctx, addr[0])
	perpetual.WhitelistAddress(ctx, addr[1])

	// Get all whitelisted addresses
	whitelisted := perpetual.GetAllWhitelistedAddress(ctx)

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

func SetupStableCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
	// prices set for USDT and USDC
	provider := authtypes.NewModuleAddress("provider")
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   "uatom",
		Display: "uatom",
		Decimal: 6,
	})

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     math.LegacyNewDec(23),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     math.LegacyNewDec(5),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     math.LegacyMustNewDecFromStr("5"),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

}
