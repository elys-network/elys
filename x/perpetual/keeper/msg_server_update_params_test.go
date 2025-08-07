package keeper_test

import (
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/testutil/sample"
	leveragelpmodulekeeper "github.com/elys-network/elys/v7/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v7/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/keeper"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestMsgServerUpdateParams_ErrUnauthorised() {
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.UpdateParams(suite.ctx, &types.MsgUpdateParams{
		Authority: sample.AccAddress(),
		Params:    &types.Params{},
	})
	suite.Require().ErrorIs(err, govtypes.ErrInvalidSigner)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerUpdateParams_ErrSetParams() {
	params := types.DefaultParams()
	params.LeverageMax = math.LegacyNewDec(-12)
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.UpdateParams(suite.ctx, &types.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    &params,
	})
	suite.Require().ErrorContains(err, "LeverageMax is negative")
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerUpdateParams_Successful() {
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)

	params := types.DefaultGenesis().Params
	params.LeverageMax = math.LegacyNewDec(int64(200000))

	_, err := msg.UpdateParams(suite.ctx, &types.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    &params,
	})
	suite.Require().Nil(err)
	gotParams := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	suite.Equal(params.LeverageMax, gotParams.LeverageMax)
}

func (suite *PerpetualKeeperTestSuite) TestMsgServerUpdateParams_Pools() {

	suite.SetupCoinPrices()
	addr := suite.AddAccounts(1, nil)
	amount := math.NewInt(1000)
	poolCreator := addr[0]
	ammPool := suite.CreateNewAmmPool(poolCreator, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:       ammPool.PoolId,
			LeverageMax:     math.LegacyMustNewDecFromStr("10"),
			AdlTriggerRatio: math.LegacyNewDec(1),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	params := types.DefaultGenesis().Params
	params.LeverageMax = math.LegacyNewDec(int64(200000))

	_, err = msg.UpdateParams(suite.ctx, &types.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    &params,
	})
	suite.Require().Nil(err)
	gotParams := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	suite.Equal(params.LeverageMax, gotParams.LeverageMax)
}
