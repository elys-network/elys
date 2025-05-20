package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/elys-network/elys/v4/x/oracle/keeper"
	"github.com/elys-network/elys/v4/x/oracle/types"
)

func (suite *KeeperTestSuite) TestPriceFeederMsgServerUpdate() {
	creator := authtypes.NewModuleAddress("A").String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgSetPriceFeeder
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgSetPriceFeeder{
				Feeder:   creator,
				IsActive: false,
			},
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			k, ctx := suite.app.OracleKeeper, suite.ctx
			params := types.DefaultParams()
			suite.app.OracleKeeper.SetParams(ctx, params)
			suite.app.OracleKeeper.SetPriceFeeder(ctx, types.PriceFeeder{
				Feeder:   creator,
				IsActive: true,
			})

			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.SetPriceFeeder(ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				rst, found := k.GetPriceFeeder(ctx,
					sdk.MustAccAddressFromBech32(tc.request.Feeder),
				)
				suite.Require().True(found)
				suite.Require().Equal(tc.request.Feeder, rst.Feeder)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPriceFeederMsgServerDelete() {
	creator := authtypes.NewModuleAddress("A").String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeletePriceFeeder
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeletePriceFeeder{
				Feeder: creator,
			},
		},
	} {
		suite.Run(tc.desc, func() {
			k, ctx := suite.app.OracleKeeper, suite.ctx
			params := types.DefaultParams()
			suite.app.OracleKeeper.SetParams(ctx, params)
			suite.app.OracleKeeper.SetPriceFeeder(ctx, types.PriceFeeder{
				Feeder:   creator,
				IsActive: true,
			})

			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.DeletePriceFeeder(ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				_, found := k.GetPriceFeeder(ctx, sdk.MustAccAddressFromBech32(tc.request.Feeder))
				suite.Require().False(found)
			}
		})
	}
}
