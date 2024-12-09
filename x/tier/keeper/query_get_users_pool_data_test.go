package keeper_test

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	ltypes "github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryGetUsersPoolDataInvalidRequest() {
	_, err := suite.app.TierKeeper.GetUsersPoolData(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryGetUsersPoolDataSuccessful() {

	addr := suite.AddAccounts(2, nil)
	poolCreator := addr[0]

	amount := math.NewInt(500033400000)
	ammPool := suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount, amount)

	msgBond := stabletypes.MsgBond{
		Creator: addr[0].String(),
		Amount:  math.NewInt(10_000_000_000),
	}

	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
	_, err := stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	suite.Require().NoError(err)

	collateralAmount := sdkmath.NewInt(10_000_000)

	_, err = suite.app.LeveragelpKeeper.Open(suite.ctx, &ltypes.MsgOpen{
		Creator:          addr[1].String(),
		CollateralAsset:  ptypes.BaseCurrency,
		CollateralAmount: collateralAmount,
		AmmPoolId:        ammPool.PoolId,
		Leverage:         sdkmath.LegacyMustNewDecFromStr("2.0"),
		StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("50.0"),
	})
	suite.Require().NoError(err)
	_, err = suite.app.TierKeeper.GetUsersPoolData(suite.ctx, &types.QueryGetUsersPoolDataRequest{})
	suite.Require().NoError(err)
}
