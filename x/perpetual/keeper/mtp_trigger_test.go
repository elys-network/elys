package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) resetForMTPTriggerChecksAndUpdates() (types.MTP, types.Pool, ammtypes.Pool, sdk.AccAddress) {
	suite.ResetSuite()
	addr := suite.AddAccounts(1, nil)
	positionCreator := addr[0]
	_, _, ammPool := suite.SetPerpetualPool(1)
	tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)
	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_LONG,
		PoolId:          ammPool.PoolId,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000_000)),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	openPositionMsg2 := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_SHORT,
		PoolId:          ammPool.PoolId,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000_000)),
		TakeProfitPrice: tradingAssetPrice.QuoInt64(4),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	mtpOpenResponse, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
	suite.Require().NoError(err)
	_, err = suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg2)
	suite.Require().NoError(err)
	mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, ammPool.PoolId, positionCreator, mtpOpenResponse.Id)
	suite.Require().NoError(err)
	pool, _ := suite.app.PerpetualKeeper.GetPool(suite.ctx, mtp.Id)
	ammPool, _ = suite.app.PerpetualKeeper.GetAmmPool(suite.ctx, mtp.AmmPoolId)
	return mtp, pool, ammPool, addr[0]
}

func (suite *PerpetualKeeperTestSuite) TestMTPTriggerChecksAndUpdates() {
	mtp, pool, ammPool, _ := suite.resetForMTPTriggerChecksAndUpdates()
	// Define test cases
	testCases := []struct {
		name           string
		setup          func()
		expectedErrMsg string
		repayAmount    math.Int
	}{
		{
			"asset profile not found",
			func() {
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
			},
			"asset info uusdc not found",
			math.NewInt(0),
		},
		{
			"Success: force close when unable to pay funding fee",
			func() {
				mtp, pool, ammPool, _ = suite.resetForMTPTriggerChecksAndUpdates()
				mtp.LastFundingCalcBlock = 1
				mtp.LastFundingCalcTime = 1
				suite.ctx = suite.ctx.WithBlockHeight(1).WithBlockTime(time.Now())
				suite.app.PerpetualKeeper.SetFundingRate(suite.ctx, 1, 1, types.FundingRateBlock{
					FundingRateLong:   math.LegacyNewDec(10000),
					FundingRateShort:  math.LegacyNewDec(10000),
					FundingShareShort: math.LegacyNewDec(10000),
					FundingShareLong:  math.LegacyNewDec(10000),
					BlockHeight:       1,
					BlockTime:         1,
				})
				pool.FundingRate = math.LegacyNewDec(1000_000)
				pool.BorrowInterestRate = math.LegacyNewDec(1)
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
			"",
			math.NewInt(0),
		},
		{
			"paying interest fail",
			func() {
				mtp, pool, ammPool, _ = suite.resetForMTPTriggerChecksAndUpdates()
				mtp.LastInterestCalcBlock = 1
				mtp.LastInterestCalcTime = 1
				mtp.LastFundingCalcBlock = 1
				mtp.LastFundingCalcTime = 1
				suite.ctx = suite.ctx.WithBlockHeight(1).WithBlockTime(time.Now())
				suite.app.PerpetualKeeper.SetFundingRate(suite.ctx, 1, 1, types.FundingRateBlock{
					FundingRateLong:   math.LegacyNewDec(0),
					FundingRateShort:  math.LegacyNewDec(0),
					FundingShareShort: math.LegacyNewDec(0),
					FundingShareLong:  math.LegacyNewDec(0),
					BlockHeight:       1,
					BlockTime:         1,
				})
				suite.app.PerpetualKeeper.SetBorrowRate(suite.ctx, 1, 1, types.InterestBlock{
					InterestRate: math.LegacyNewDec(1000_000),
					BlockHeight:  1,
					BlockTime:    1,
				})
				pool.FundingRate = math.LegacyNewDec(0)
				pool.BorrowInterestRate = math.LegacyNewDec(1000_000)
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
			"",
			math.NewInt(0),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
			_, _, _, _, _, _, _, _, _, _, err := suite.app.PerpetualKeeper.MTPTriggerChecksAndUpdates(suite.ctx, &mtp, &pool, &ammPool)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
