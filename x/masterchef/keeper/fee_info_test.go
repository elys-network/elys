package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/elys-network/elys/v5/x/masterchef/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *MasterchefKeeperTestSuite) TestFeeInfo() {

	dateString := suite.ctx.BlockTime().Format("2006-01-02")
	dayAfterDateString := suite.ctx.BlockTime().AddDate(0, 0, 1).Format("2006-01-02")
	twoDayAfterDateString := suite.ctx.BlockTime().AddDate(0, 0, 2).Format("2006-01-02")
	addFeeInfoTestCases := []struct {
		name        string
		lp          osmomath.BigDec
		expLp       sdkmath.Int
		stakers     osmomath.BigDec
		expStakers  sdkmath.Int
		protocol    osmomath.BigDec
		expProtocol sdkmath.Int
		gas         bool
		date        string
	}{
		{
			name:        "AddFeeInfo for gas fees",
			lp:          osmomath.NewBigDec(100),
			expLp:       sdkmath.NewInt(100),
			stakers:     osmomath.NewBigDec(50),
			expStakers:  sdkmath.NewInt(50),
			protocol:    osmomath.NewBigDec(25),
			expProtocol: sdkmath.NewInt(25),
			gas:         true,
			date:        dateString,
		},
		{
			name:        "AddFeeInfo for dex fees",
			lp:          osmomath.NewBigDec(200),
			expLp:       sdkmath.NewInt(200),
			stakers:     osmomath.NewBigDec(100),
			expStakers:  sdkmath.NewInt(100),
			protocol:    osmomath.NewBigDec(50),
			expProtocol: sdkmath.NewInt(50),
			gas:         false,
			date:        dateString,
		},
	}

	for _, tc := range addFeeInfoTestCases {
		suite.Run(tc.name, func() {
			suite.app.MasterchefKeeper.AddFeeInfo(suite.ctx, tc.lp, tc.stakers, tc.protocol, tc.gas)
			info := suite.app.MasterchefKeeper.GetFeeInfo(suite.ctx, tc.date)
			if tc.gas {
				suite.Require().Equal(tc.expLp, info.GasLp)
				suite.Require().Equal(tc.expStakers, info.GasStakers)
				suite.Require().Equal(tc.expProtocol, info.GasProtocol)
			} else {
				suite.Require().Equal(tc.expLp, info.DexLp)
				suite.Require().Equal(tc.expStakers, info.DexStakers)
				suite.Require().Equal(tc.expProtocol, info.DexProtocol)
			}
		})
	}

	addEdenInfoTestCases := []struct {
		name    string
		eden    osmomath.BigDec
		expEden sdkmath.Int
		date    string
	}{
		{
			name:    "Add eden Info",
			eden:    osmomath.NewBigDec(100),
			expEden: sdkmath.NewInt(100),
			date:    dateString,
		},
	}

	for _, tc := range addEdenInfoTestCases {
		suite.Run(tc.name, func() {
			suite.app.MasterchefKeeper.AddEdenInfo(suite.ctx, tc.eden)
			info := suite.app.MasterchefKeeper.GetFeeInfo(suite.ctx, tc.date)
			suite.Require().Equal(tc.expEden, info.EdenLp)
		})
	}

	feeInfo := types.FeeInfo{
		GasLp:        sdkmath.NewInt(300),
		GasStakers:   sdkmath.NewInt(150),
		GasProtocol:  sdkmath.NewInt(75),
		DexLp:        sdkmath.NewInt(400),
		DexStakers:   sdkmath.NewInt(200),
		DexProtocol:  sdkmath.NewInt(100),
		PerpLp:       sdkmath.NewInt(500),
		PerpStakers:  sdkmath.NewInt(250),
		PerpProtocol: sdkmath.NewInt(125),
		EdenLp:       sdkmath.NewInt(2000),
	}
	zeroInfo := types.FeeInfo{
		GasLp:        sdkmath.NewInt(0),
		GasStakers:   sdkmath.NewInt(0),
		GasProtocol:  sdkmath.NewInt(0),
		DexLp:        sdkmath.NewInt(0),
		DexStakers:   sdkmath.NewInt(0),
		DexProtocol:  sdkmath.NewInt(0),
		PerpLp:       sdkmath.NewInt(0),
		PerpStakers:  sdkmath.NewInt(0),
		PerpProtocol: sdkmath.NewInt(0),
		EdenLp:       sdkmath.NewInt(0),
	}

	setFeeInfoTestCases := []struct {
		name                 string
		info                 types.FeeInfo
		expInfo              types.FeeInfo
		setDate              string
		getonExistentDate    string
		getOnnonExistentDate string
	}{
		{
			name:                 "Set and get Fee Info on data",
			info:                 feeInfo,
			expInfo:              feeInfo,
			setDate:              dayAfterDateString,
			getonExistentDate:    dayAfterDateString,
			getOnnonExistentDate: twoDayAfterDateString,
		},
	}

	for _, tc := range setFeeInfoTestCases {
		suite.Run(tc.name, func() {
			suite.app.MasterchefKeeper.SetFeeInfo(suite.ctx, tc.info, tc.setDate)

			infoOnExistentDate := suite.app.MasterchefKeeper.GetFeeInfo(suite.ctx, tc.getonExistentDate)
			suite.Require().Equal(tc.expInfo, infoOnExistentDate)

			infoOnnonExistentDate := suite.app.MasterchefKeeper.GetFeeInfo(suite.ctx, tc.getOnnonExistentDate)
			suite.Require().Equal(zeroInfo, infoOnnonExistentDate)
		})
	}

	//Test RemoveFeeInfo
	suite.app.MasterchefKeeper.RemoveFeeInfo(suite.ctx, dateString)
	removedInfo := suite.app.MasterchefKeeper.GetFeeInfo(suite.ctx, dateString)
	suite.Require().Equal(zeroInfo, removedInfo)

	// Test GetAllFeeInfos
	allInfos := suite.app.MasterchefKeeper.GetAllFeeInfos(suite.ctx)
	suite.Require().Len(allInfos, 1)
	suite.Require().Equal(feeInfo, allInfos[0])
}
