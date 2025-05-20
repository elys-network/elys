package keeper_test

import (
	"sort"
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/v4/app"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	aptypes "github.com/elys-network/elys/v4/x/assetprofile/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v4/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v4/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/v4/x/oracle/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/suite"
)

type assetPriceInfo struct {
	denom   string
	display string
	price   osmomath.BigDec
}

const (
	initChain = true
)

var (
	oracleProvider = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	priceMap = map[string]assetPriceInfo{
		"uusdc": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   osmomath.OneBigDec(),
		},
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   osmomath.MustNewBigDecFromStr("5.0"),
		},
	}

	entry = []aptypes.Entry{
		{
			BaseDenom:   ptypes.BaseCurrency,
			Denom:       ptypes.BaseCurrency,
			Decimals:    6,
			DisplayName: "USDC",
		},
		{
			BaseDenom:   ptypes.ATOM,
			Denom:       ptypes.ATOM,
			Decimals:    6,
			DisplayName: "ATOM",
		},
	}
)

type TradeshieldKeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *TradeshieldKeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain, suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(TradeshieldKeeperTestSuite))
}

func (suite *TradeshieldKeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *TradeshieldKeeperTestSuite) SetupCoinPrices() {
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
			Price:     v.price.Dec(),
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(suite.ctx.BlockTime().Unix()),
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) GetAccountIssueAmount() math.Int {
	return math.NewInt(10_000_000_000_000)
}

func (suite *TradeshieldKeeperTestSuite) SetupAssetProfile() {
	for _, v := range entry {
		suite.app.AssetprofileKeeper.SetEntry(suite.ctx, v)
	}
}

func (suite *TradeshieldKeeperTestSuite) AddAccounts(n int, given []sdk.AccAddress) []sdk.AccAddress {
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

func (suite *TradeshieldKeeperTestSuite) CreateNewAmmPool(creator sdk.AccAddress, useOracle bool, swapFee, exitFee math.LegacyDec, asset2 string, baseTokenAmount, assetAmount math.Int) ammtypes.Pool {
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

func (suite *TradeshieldKeeperTestSuite) SetPerpetualPool(poolId uint64) (types.Pool, sdk.AccAddress, ammtypes.Pool) {
	ctx := suite.ctx
	k := suite.app.PerpetualKeeper
	//prices
	suite.SetupCoinPrices()
	suite.SetupAssetProfile()
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

	pool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("11.5"))
	k.SetPool(ctx, pool)

	return pool, poolCreator, ammPool
}
