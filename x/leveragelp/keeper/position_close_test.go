package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (suite KeeperTestSuite) TestCloseLong() {
	k := suite.app.LeveragelpKeeper

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
		mtp = types.MTP{
			AmmPoolId: 2,
		}
		pool        = types.Pool{}
		repayAmount = math.NewInt(100)
	)

	k.SetMTP(suite.ctx, &mtp)
	k.SetPool(suite.ctx, pool)

	mtpOut, repayAmountOut, err := k.CloseLong(ctx, msg)

	// Expect no error
	suite.Require().Nil(err)
	suite.Require().Equal(repayAmount, repayAmountOut)
	suite.Require().Equal(mtp, *mtpOut)
}
