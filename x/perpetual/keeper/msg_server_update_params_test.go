package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/testutil/sample"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
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
	msg := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
	_, err := msg.UpdateParams(suite.ctx, &types.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params: &types.Params{
			LeverageMax: math.LegacyNewDec(int64(-200)),
		},
	})
	suite.Require().ErrorContains(err, "leverage max must be positive")
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
	amount := sdk.NewInt(1000)
	poolId := uint64(1)
	poolCreator := addr[0]
	_ = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
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
