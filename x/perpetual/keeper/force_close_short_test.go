package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestForceCloseShort_PoolNotFound() {

	ctx := suite.ctx
	k := suite.app.PerpetualKeeper
	mtp := &types.MTP{
		AmmPoolId: 8000,
	}

	pool := &types.Pool{}

	_, err := k.ForceCloseShort(ctx, mtp, pool, false, ptypes.BaseCurrency)

	suite.Require().ErrorIs(err, types.ErrPoolDoesNotExist)
}

func (suite *PerpetualKeeperTestSuite) TestForceCloseShort_Successful() {

	ctx := suite.ctx
	k := suite.app.PerpetualKeeper
	//prices

	suite.SetupCoinPrices()
	//accounts
	accounts := suite.AddAccounts(2, nil)
	poolCreator := accounts[0]
	positionCreator := accounts[1]

	amount := math.NewInt(1000)

	ammPool := suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   ammPool.PoolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	suite.Require().NoError(err)

	pool := types.NewPool(ammPool)
	k.SetPool(ctx, pool)

	openPositionMsg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_SHORT,
		PoolId:          ammPool.PoolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("0.95"),
		StopLossPrice:   math.LegacyZeroDec(),
	}

	position, err := k.Open(ctx, openPositionMsg, false)

	suite.Require().Nil(err)

	mtp, err := k.GetMTP(ctx, positionCreator, position.Id)

	suite.Require().Nil(err)

	_, err = k.ForceCloseShort(ctx, &mtp, &pool, false, ptypes.BaseCurrency)

	suite.Require().Nil(err)
}
