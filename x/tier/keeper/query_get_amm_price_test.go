package keeper_test

import (
	"cosmossdk.io/math"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/elys-network/elys/v5/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryGetAmmPriceInvalidRequest() {
	_, err := suite.app.TierKeeper.GetAmmPrice(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryGetAmmPriceSuccessful() {

	addr := suite.AddAccounts(2, nil)
	poolCreator := addr[0]

	amount := math.NewInt(500033400000)
	_ = suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount, amount)

	_, err := suite.app.TierKeeper.GetAmmPrice(suite.ctx, &types.QueryGetAmmPriceRequest{
		Denom:   ptypes.ATOM,
		Decimal: 6,
	})

	suite.Require().NoError(err)
}
