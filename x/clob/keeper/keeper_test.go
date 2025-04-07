package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/stretchr/testify/suite"
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
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileAtom)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileUsdc)
}

func (suite *KeeperTestSuite) SetPrice(assets []string, prices []math.LegacyDec) {
	if len(assets) != len(prices) {
		panic("unequal lengths while setting prices during test")
	}
	for i, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
			Asset:       assets[i],
			Price:       price,
			Source:      "test",
			Provider:    "test",
			Timestamp:   uint64(time.Now().Unix()),
			BlockHeight: 1,
		})
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
