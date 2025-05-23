package keeper_test

import (
	"cosmossdk.io/math"
	typessdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/elys-network/elys/v5/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryStakedPoolSInvalidRequest() {
	k := suite.app.TierKeeper
	_, err := k.StakedPool(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryStakedPoolSuccessful() {
	addr := suite.AddAccounts(1, nil)[0]
	poolCreator := addr

	amount := math.NewInt(500033400000)
	ammPool := suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount, amount)
	tokens := []typessdk.Coin{
		{
			Denom:  ptypes.ATOM,
			Amount: math.NewInt(10000000),
		},
		{
			Denom:  ptypes.BaseCurrency,
			Amount: math.NewInt(540000000),
		},
	}
	estimation, err := suite.app.AmmKeeper.JoinPoolEstimation(suite.ctx, &ammtypes.QueryJoinPoolEstimationRequest{
		PoolId:    ammPool.PoolId,
		AmountsIn: tokens,
	})
	suite.Require().NoError(err)

	_, _, err = suite.app.AmmKeeper.JoinPoolNoSwap(
		suite.ctx,
		suite.account,
		ammPool.PoolId,
		estimation.ShareAmountOut.Amount.ToLegacyDec().TruncateInt(),
		tokens)
	suite.Require().NoError(err)

	stakedPoolResponse, err := suite.app.TierKeeper.StakedPool(suite.ctx, &types.QueryStakedPoolRequest{
		User: suite.account.String(),
	})

	suite.Require().NoError(err)
	suite.Require().Equal("59.999999999999999999", stakedPoolResponse.Total.String())

}
