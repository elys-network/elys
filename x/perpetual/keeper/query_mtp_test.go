package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/testutil/sample"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestQueryMtp_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg, false)
	suite.Require().NoError(err)

	response, err := k.MTP(ctx, &types.MTPRequest{
		Address: firstPositionCreator.String(),
		Id:      firstPosition.Id,
		PoolId:  firstPool,
	})

	suite.Require().Nil(err)
	suite.Require().Equal(response.Mtp.Mtp.Id, firstPosition.Id)
	suite.Require().Equal(response.Mtp.Mtp.AmmPoolId, firstPool)
}

func (suite *PerpetualKeeperTestSuite) TestQueryMtp_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.MTP(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryMtp_ErrMTPDoesNotExist() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.MTP(ctx, &types.MTPRequest{
		Address: sample.AccAddress(),
		Id:      2,
	})

	suite.Require().ErrorIs(err, types.ErrMTPDoesNotExist)
}

func (suite *PerpetualKeeperTestSuite) TestQueryMtp_ErrAccAddressFromBech32() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.MTP(ctx, &types.MTPRequest{
		Address: "error",
		Id:      2,
	})

	suite.Require().ErrorContains(err, "invalid bech32 string length 5")
}

func (suite *PerpetualKeeperTestSuite) TestQueryMtp_BaseCurrencyNotFound() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg, false)
	suite.Require().NoError(err)

	suite.app.AssetprofileKeeper.RemoveEntry(ctx, ptypes.BaseCurrency)

	_, err = k.MTP(ctx, &types.MTPRequest{
		Address: firstPositionCreator.String(),
		Id:      firstPosition.Id,
		PoolId:  firstPool,
	})

	suite.Require().ErrorContains(err, "base currency not found")
}

func (suite *PerpetualKeeperTestSuite) TestQueryMtp_ErrFillMTPData() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg, false)
	suite.Require().NoError(err)

	suite.app.AmmKeeper.RemovePool(ctx, firstPool)

	_, err = k.MTP(ctx, &types.MTPRequest{
		Address: firstPositionCreator.String(),
		Id:      firstPosition.Id,
	})

	suite.Require().ErrorContains(err, "not found")
}
