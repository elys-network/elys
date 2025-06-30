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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	simapp "github.com/elys-network/elys/v6/app"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	atypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v6/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v6/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/v6/x/stablestake/types"
	"github.com/stretchr/testify/suite"
)

type TierKeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
	account     sdk.AccAddress
}

const (
	initChain = true
)

func (suite *TierKeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain, suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
	suite.SetStakingParam()
	suite.SetStableStakeParam()
	suite.SetupAssetProfile()
	suite.SetupCoinPrices()
	suite.SetAccount()

	suite.app.StablestakeKeeper.SetPool(suite.ctx, stablestaketypes.Pool{
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		NetAmount:            math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		Id:                   1,
		DepositDenom:         ptypes.BaseCurrency,
	})
}

func (suite *TierKeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(TierKeeperTestSuite))
}

func (suite *TierKeeperTestSuite) SetStakingParam() error {
	return suite.app.StakingKeeper.SetParams(suite.ctx, stakingtypes.Params{
		UnbondingTime:     1209600,
		MaxValidators:     60,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		BondDenom:         "uelys",
		MinCommissionRate: math.LegacyNewDec(0),
	})
}

func (suite *TierKeeperTestSuite) SetStableStakeParam() error {

	params := stablestaketypes.DefaultParams()
	suite.app.StablestakeKeeper.SetParams(suite.ctx, params)
	return nil
}

func (suite *TierKeeperTestSuite) SetAccount() {
	account := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	//10.000 USDC
	usdcAmount := sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 10000000000)}
	// bootstrap balances
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, usdcAmount)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, account, usdcAmount)
	suite.Require().NoError(err)

	//10.000 ATOM
	atomAmount := sdk.Coins{sdk.NewInt64Coin(ptypes.ATOM, 10000000000)}
	// bootstrap balances
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, atomAmount)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, account, atomAmount)
	suite.Require().NoError(err)

	//Elys
	elysCoins := sdk.NewCoins(sdk.NewCoin(ptypes.Elys, math.NewInt(1000000)))

	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, elysCoins)
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, account, elysCoins)
	suite.Require().NoError(err)

	suite.account = account
}

func (suite *TierKeeperTestSuite) SetupAssetProfile() {

	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, atypes.Entry{
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

	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, atypes.Entry{
		BaseDenom:       ptypes.Elys,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})

}

func (suite *TierKeeperTestSuite) SetupCoinPrices() {
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

var oracleProvider = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

type assetPriceInfo struct {
	denom   string
	display string
	price   math.LegacyDec
}

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
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   math.LegacyMustNewDecFromStr("5.0"),
		},
	}
)

func (suite *TierKeeperTestSuite) CreateNewAmmPool(creator sdk.AccAddress, useOracle bool, swapFee, exitFee math.LegacyDec, asset2 string, baseTokenAmount, assetAmount math.Int) ammtypes.Pool {
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

	//Set Perpetual
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:            poolId,
			LeverageMax:          math.LegacyMustNewDecFromStr("10"),
			PoolMaxLeverageRatio: math.LegacyMustNewDecFromStr("0.99"),
		},
	}

	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	pool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("11.5"))
	k := suite.app.PerpetualKeeper
	k.SetPool(suite.ctx, pool)

	return ammPool
}

func (suite *TierKeeperTestSuite) AddAccounts(n int, given []sdk.AccAddress) []sdk.AccAddress {
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

func (suite *TierKeeperTestSuite) GetAccountIssueAmount() math.Int {
	return math.NewInt(10_000_000_000_000)
}
