package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
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
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   sdk.ZeroDec(),
	}

	firstPosition, err := suite.app.PerpetualKeeper.Open(ctx, firstOpenPositionMsg, false)
	suite.Require().NoError(err)

	_, err = k.MTP(ctx, &types.MTPRequest{
		Address: firstPositionCreator.String(),
		Id:      firstPosition.Id,
	})

	suite.Require().Nil(err)
}
