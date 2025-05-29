package keeper_test

import (
	"cosmossdk.io/math"
	aptypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestGetAssetPriceAndAssetUsdcDenomRatio() {
	suite.SetupCoinPrices()
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, aptypes.Entry{
		BaseDenom:       ptypes.BaseCurrency,
		Denom:           ptypes.BaseCurrency,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "uusdc",
		Price:     osmomath.MustNewDecFromStr("0.98"),
		Source:    "elys",
		Provider:  oracleProvider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})

	tradingAssetPrice, tradingAssetPriceDenomRatio, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)
	suite.Require().Equal(math.LegacyMustNewDecFromStr("5"), tradingAssetPrice)
	suite.Require().Equal(osmomath.MustNewBigDecFromStr("5.102040816326530612244897959183673469"), tradingAssetPriceDenomRatio)

	tradingAssetPrice, tradingAssetPriceDenomRatio, err = suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, "wei")
	suite.Require().NoError(err)
	suite.Require().Equal(math.LegacyMustNewDecFromStr("1500"), tradingAssetPrice)
	suite.Require().Equal(osmomath.MustNewBigDecFromStr("0.000000001530612244897959183673469387"), tradingAssetPriceDenomRatio)
}

func (suite *PerpetualKeeperTestSuite) TestConvertPriceToAssetUsdcDenomRatio() {
	suite.SetupCoinPrices()
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, aptypes.Entry{
		BaseDenom:       ptypes.BaseCurrency,
		Denom:           ptypes.BaseCurrency,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "uusdc",
		Price:     osmomath.MustNewDecFromStr("0.98"),
		Source:    "elys",
		Provider:  oracleProvider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})

	tradingAssetPriceDenomRatio, err := suite.app.PerpetualKeeper.ConvertPriceToAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM, math.LegacyMustNewDecFromStr("5"))
	suite.Require().NoError(err)
	suite.Require().Equal(osmomath.MustNewBigDecFromStr("5.102040816326530612244897959183673469"), tradingAssetPriceDenomRatio)

	tradingAssetPriceDenomRatio, err = suite.app.PerpetualKeeper.ConvertPriceToAssetUsdcDenomRatio(suite.ctx, "wei", math.LegacyMustNewDecFromStr("1500"))
	suite.Require().NoError(err)
	suite.Require().Equal(osmomath.MustNewBigDecFromStr("0.000000001530612244897959183673469387"), tradingAssetPriceDenomRatio)

	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:      "afet",
		Display:    "FET",
		BandTicker: "FET",
		ElysTicker: "FET",
		Decimal:    18,
	})

	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "uusdc",
		Price:     osmomath.MustNewDecFromStr("1"),
		Source:    "elys",
		Provider:  oracleProvider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})

	tradingAssetPriceDenomRatio, err = suite.app.PerpetualKeeper.ConvertPriceToAssetUsdcDenomRatio(suite.ctx, "afet", math.LegacyMustNewDecFromStr("0.5"))
	suite.Require().NoError(err)
	suite.Require().Equal(osmomath.MustNewBigDecFromStr("0.0000000000005"), tradingAssetPriceDenomRatio)
}
