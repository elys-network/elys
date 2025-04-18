package keeper_test

import (
	"cosmossdk.io/math"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/clob/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
	"time"
)

var (
	assetProfileAtom = assetprofiletypes.Entry{
		BaseDenom:                "uatom",
		Decimals:                 6,
		Denom:                    "uatom",
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "ATOM",
		DisplaySymbol:            "ATOM",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                "uatom",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	oracleProfileAtom = oracletypes.AssetInfo{
		Denom:      "uatom",
		Display:    "ATOM",
		BandTicker: "ATOM",
		ElysTicker: "ATOM",
		Decimal:    6,
	}
	assetProfileUsdc = assetprofiletypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "uusdc",
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "USDC",
		DisplaySymbol:            "USDC",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                "uusdc",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	oracleProfileUsdc = oracletypes.AssetInfo{
		Denom:      "uusdc",
		Display:    "USDC",
		BandTicker: "USDC",
		ElysTicker: "USDC",
		Decimal:    6,
	}
	assetProfileOsmo = assetprofiletypes.Entry{
		BaseDenom:                "uosmo",
		Decimals:                 6,
		Denom:                    "uosmo",
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "OSMO",
		DisplaySymbol:            "OSMO",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                "uosmo",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	oracleProfileOsmo = oracletypes.AssetInfo{
		Denom:      "uosmo",
		Display:    "OSMO",
		BandTicker: "OSMO",
		ElysTicker: "OSMO",
		Decimal:    6,
	}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp

	avgBlockTime uint64
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(true, suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(true).WithBlockTime(time.Now())
	suite.app = app
	suite.avgBlockTime = 5

	oracleParams := app.OracleKeeper.GetParams(suite.ctx)
	oracleParams.LifeTimeInBlocks = 10000
	oracleParams.PriceExpiryTime = 84600
	app.OracleKeeper.SetParams(suite.ctx, oracleParams)

	suite.SetAssetProfiles()
	suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{math.LegacyNewDec(10), math.LegacyNewDec(1)})
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) SetAssetProfiles() {
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetProfileAtom)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetProfileUsdc)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetProfileOsmo)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileAtom)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileUsdc)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileOsmo)
}

func (suite *KeeperTestSuite) SetPrice(assets []string, prices []math.LegacyDec) {
	if len(assets) != len(prices) {
		panic("unequal lengths while setting prices during test")
	}
	for i, price := range prices {
		oraclePrice := oracletypes.Price{
			Asset:       assets[i],
			Price:       price,
			Source:      "test",
			Provider:    "test",
			Timestamp:   uint64(suite.ctx.BlockTime().Unix()),
			BlockHeight: uint64(suite.ctx.BlockHeight()),
		}
		suite.app.OracleKeeper.SetPrice(suite.ctx, oraclePrice)
	}
}

func (suite *KeeperTestSuite) IncreaseHeight(height uint64) {
	if height == 0 {
		panic("increment cannot be 0")
	}
	for i := uint64(1); i <= height; i++ {
		//_, err := suite.app.BeginBlocker(suite.ctx)
		//if err != nil {
		//	panic(err)
		//}
		//_, err = suite.app.EndBlocker(suite.ctx)
		//if err != nil {
		//	panic(err)
		//}
		currentHeight := suite.ctx.BlockHeight()
		currentTime := suite.ctx.BlockTime().Unix()
		ctx := suite.ctx.WithBlockHeight(currentHeight + 1)
		ctx = ctx.WithBlockTime(time.Unix(currentTime+int64(suite.avgBlockTime), 0))
		suite.ctx = ctx
	}
}

func (suite *KeeperTestSuite) SetupSubAccounts(total uint64, balance sdk.Coins) []types.SubAccount {
	if total == 0 {
		panic("total subaccounts cannot be 0")
	}

	all := suite.app.ClobKeeper.GetAllSubAccount(suite.ctx)
	var list []types.SubAccount
	for i := uint64(len(all) + 1); i <= uint64(len(all))+total; i++ {
		err := suite.app.BankKeeper.MintCoins(suite.ctx, ammtypes.ModuleName, balance)
		suite.Require().NoError(err)

		subAccountAddress := authtypes.NewModuleAddress("subAccount" + strconv.FormatUint(i, 10))

		subAccount := types.SubAccount{
			Owner:            subAccountAddress.String(),
			MarketId:         1,
			AvailableBalance: balance,
			TotalBalance:     balance,
			TradeNounce:      0,
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ammtypes.ModuleName, subAccount.GetTradingAccountAddress(), balance)
		suite.Require().NoError(err)

		suite.app.ClobKeeper.SetSubAccount(suite.ctx, subAccount)
		list = append(list, subAccount)
	}
	return list
}

func (suite *KeeperTestSuite) CreateMarket(baseDenoms ...string) []types.PerpetualMarket {
	all := suite.app.ClobKeeper.GetAllPerpetualMarket(suite.ctx)
	var list []types.PerpetualMarket
	for _, baseDenom := range baseDenoms {
		if baseDenom == "" || baseDenom == "uusdc" {
			panic("base Denom cannot be uusdc or empty")
		}
		market := types.PerpetualMarket{
			Id:                     uint64(len(all) + 1),
			BaseDenom:              baseDenom,
			QuoteDenom:             "uusdc",
			InitialMarginRatio:     math.LegacyMustNewDecFromStr("0.1"),
			MaintenanceMarginRatio: math.LegacyMustNewDecFromStr("0.2"),
			Status:                 1,
			MaxFundingRate:         math.LegacyMustNewDecFromStr("0.05"),
			MaxFundingRateChange:   math.LegacyMustNewDecFromStr("0.01"),
			MaxTwapPricesTime:      15,
		}
		suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
		list = append(list, market)
	}
	return list
}

func (suite *KeeperTestSuite) GetAccountBalance(addr sdk.AccAddress, denom string) math.Int {
	return suite.app.BankKeeper.GetBalance(suite.ctx, addr, denom).Amount
}

func (suite *KeeperTestSuite) FundAccount(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) WithdrawFromAccount(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, addr, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) SetAccountBalance(addr sdk.AccAddress, coins sdk.Coins) {
	// Implementation depends on test setup - might need direct bank keeper state setting or burn+fund
	currentCoins := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr)
	err := suite.app.BankKeeper.SendCoins(suite.ctx, addr, authtypes.NewModuleAddress("dummy_burn"), currentCoins) // Send all current away
	if err != nil && !errors.Is(err, sdkerrors.ErrInsufficientFunds) {                                             // Ignore error if already empty
		suite.T().Logf("Error burning coins during SetAccountBalance: %v", err) // Log non-critical error
	}
	if !coins.IsZero() {
		suite.FundAccount(addr, coins) // Fund with desired amount
	}
}
