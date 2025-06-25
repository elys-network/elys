package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	tiertypes "github.com/elys-network/elys/v6/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryPerpetualRequest() {
	k := suite.app.TierKeeper
	_, err := k.Perpetual(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryPerpetual() {
	addr := suite.AddAccounts(2, nil)
	poolCreator := addr[0]

	amount := math.NewInt(500033400000)
	ammPool := suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount, amount)

	firstPositionCreator := addr[1]
	amountToOpen := math.NewInt(400)

	firstOpenPositionMsg := &types.MsgOpen{
		Creator:         firstPositionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_SHORT,
		PoolId:          ammPool.PoolId,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amountToOpen),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	_, err := suite.app.PerpetualKeeper.Open(suite.ctx, firstOpenPositionMsg)
	suite.Require().NoError(err)

	_, err = suite.app.TierKeeper.Perpetual(suite.ctx, &tiertypes.QueryPerpetualRequest{
		User: firstOpenPositionMsg.Creator,
	})

	suite.Require().NoError(err)

}
