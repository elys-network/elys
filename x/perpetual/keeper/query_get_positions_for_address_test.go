package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/testutil/sample"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestQueryGetPositionsForAddress_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetPositionsForAddress(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetPositionsForAddress_ErrAccAddressFromBech32() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetPositionsForAddress(ctx, &types.PositionsForAddressRequest{
		Address: "error",
	})

	suite.Require().ErrorContains(err, "invalid bech32 string length 5")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetPositionsForAddress_ErrPageSize() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetPositionsForAddress(ctx, &types.PositionsForAddressRequest{
		Address: sample.AccAddress(),
		Pagination: &query.PageRequest{
			Limit: 12000,
		},
	})

	suite.Require().ErrorContains(err, "page size greater than max")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetPositionsForAddress_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	firstPool := uint64(1)
	suite.SetPerpetualPool(firstPool)

	amount := math.NewInt(400)
	addr := suite.AddAccounts(2, nil)
	firstPositionCreator := addr[0]
	secondPositionCreator := addr[1]

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	secondOpenPositionMsg := &types.MsgOpen{
		Creator:         secondPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          firstPool,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	_, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)
	_, err = suite.app.PerpetualKeeper.Open(ctx, secondOpenPositionMsg)
	suite.Require().NoError(err)

	response, err := k.GetPositionsForAddress(ctx, &types.PositionsForAddressRequest{
		Address: firstPositionCreator.String(),
		Pagination: &query.PageRequest{
			Limit: 10,
		},
	})

	suite.Require().Nil(err)
	suite.Require().Len(response.Mtps, 1)
}
