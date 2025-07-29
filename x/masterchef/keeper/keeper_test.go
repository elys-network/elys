package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	simapp "github.com/elys-network/elys/v7/app"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/masterchef/keeper"
	"github.com/elys-network/elys/v7/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

const (
	initChain = true
)

type MasterchefKeeperTestSuite struct {
	suite.Suite

	legacyAmino      *codec.LegacyAmino
	ctx              sdk.Context
	app              *simapp.ElysApp
	msgServer        types.MsgServer
	queryClient      types.QueryClient
	genesisAccount   sdk.AccAddress
	genValidatorAddr sdk.ValAddress
}

func (suite *MasterchefKeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain, suite.T())
	suite.app = app
	suite.SetEssentials()
}

func (suite *MasterchefKeeperTestSuite) SetupTestWithGenesisAcc() {
	app, genAcc, genVal := simapp.InitElysTestAppWithGenAccount(suite.T())
	suite.app = app
	suite.genesisAccount = genAcc
	suite.genValidatorAddr = genVal
	suite.SetEssentials()
}

func (suite *MasterchefKeeperTestSuite) SetEssentials() {
	suite.legacyAmino = suite.app.LegacyAmino()
	suite.ctx = suite.app.BaseApp.NewContext(initChain)
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.MasterchefKeeper)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.MasterchefKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	// Setup Masterchef and staking module params
	simapp.SetMasterChefParams(suite.app, suite.ctx)
	err := simapp.SetStakingParam(suite.app, suite.ctx)
	suite.Require().NoError(err)
	simapp.SetupAssetProfile(suite.app, suite.ctx)

	// Setup coin prices
	suite.SetupStableCoinPrices()
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(MasterchefKeeperTestSuite))
}

func (suite *MasterchefKeeperTestSuite) ResetSuite(withGenesisAcc bool) {
	if withGenesisAcc {
		suite.SetupTestWithGenesisAcc()
	} else {
		suite.SetupTest()
	}
}

func (suite *MasterchefKeeperTestSuite) AddPoolInfo() {
	suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, types.PoolInfo{
		PoolId:               2,
		RewardWallet:         "elys1d96rzrky937s3s397g5xh5qvcwgkeqysmh8sg2kn359fhfvzeyrsnalu2u",
		Multiplier:           sdkmath.LegacyMustNewDecFromStr("1.00"),
		GasApr:               sdkmath.LegacyMustNewDecFromStr("0.00"),
		EdenApr:              sdkmath.LegacyMustNewDecFromStr("0.50"),
		DexApr:               sdkmath.LegacyMustNewDecFromStr("0.00"),
		ExternalIncentiveApr: sdkmath.LegacyMustNewDecFromStr("0.00"),
		EnableEdenRewards:    false,
	})
}

func (suite *MasterchefKeeperTestSuite) SetupStableCoinPrices() {
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
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})

	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdkmath.LegacyNewDec(1000000),
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdkmath.LegacyNewDec(1000000),
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     sdkmath.LegacyNewDec(100),
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdkmath.LegacyNewDec(100),
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
}

func (suite *MasterchefKeeperTestSuite) CreateNewAmmPool(creator sdk.AccAddress, poolAssets []ammtypes.PoolAsset, poolParams ammtypes.PoolParams) ammtypes.Pool {

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

func (suite *MasterchefKeeperTestSuite) MintTokenToAddress(address sdk.AccAddress, amount sdkmath.Int, denom string) {

	token := sdk.NewCoins(sdk.NewCoin(denom, amount))

	err := suite.app.BankKeeper.MintCoins(suite.ctx, ammtypes.ModuleName, token)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ammtypes.ModuleName, address, token)
	suite.Require().NoError(err)
}

func (suite *MasterchefKeeperTestSuite) MintMultipleTokenToAddress(address sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, ammtypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ammtypes.ModuleName, address, coins)
	suite.Require().NoError(err)
}
