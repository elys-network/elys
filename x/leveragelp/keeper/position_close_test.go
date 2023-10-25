package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (suite KeeperTestSuite) TestCloseLong() {
	k := suite.app.LeveragelpKeeper

	var (
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
		mtp = types.MTP{
			Address:           msg.Creator,
			Collateral:        sdk.NewInt64Coin("uusdc", 1000),
			Liabilities:       sdk.NewInt(4000),
			InterestPaid:      sdk.ZeroInt(),
			Leverage:          sdk.NewDec(5),
			LeveragedLpAmount: sdk.ZeroInt(),
			MtpHealth:         sdk.OneDec(),
			Id:                1,
			AmmPoolId:         1,
		}
		pool = types.Pool{
			AmmPoolId: 1,
		}
		repayAmount = math.NewInt(0)
	)

	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId: 1,
	})
	k.SetPool(suite.ctx, pool)
	k.SetMTP(suite.ctx, &mtp)
	k.SetMTP(suite.ctx, &types.MTP{})
	mtpOut, repayAmountOut, err := k.CloseLong(suite.ctx, msg)

	// Expect no error
	suite.Require().Nil(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
	suite.Require().Equal(mtp, *mtpOut)
}

func (suite KeeperTestSuite) TestForceCloseLong() {
	k := suite.app.LeveragelpKeeper
	mtp := types.MTP{
		Address:           "creator",
		Collateral:        sdk.NewInt64Coin("uusdc", 1000),
		Liabilities:       sdk.NewInt(4000),
		InterestPaid:      sdk.ZeroInt(),
		Leverage:          sdk.NewDec(5),
		LeveragedLpAmount: sdk.ZeroInt(),
		MtpHealth:         sdk.OneDec(),
		Id:                1,
		AmmPoolId:         1,
	}
	pool := types.Pool{
		AmmPoolId: 1,
	}
	repayAmount := math.NewInt(0)
	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId: 1,
	})
	k.SetPool(suite.ctx, pool)
	k.SetMTP(suite.ctx, &mtp)
	k.SetMTP(suite.ctx, &types.MTP{})
	repayAmountOut, err := k.ForceCloseLong(suite.ctx, mtp, pool)

	// Expect no error
	suite.Require().Nil(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
}
