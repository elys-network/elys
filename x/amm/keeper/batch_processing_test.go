package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestBatchProcessing() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"GetAllSwapExactAmountInRequests",
			func() {
				suite.ResetSuite()

				lastSwapIndex := suite.app.AmmKeeper.GetLastSwapRequestIndex(suite.ctx)
				suite.app.AmmKeeper.SetSwapExactAmountInRequests(suite.ctx, &types.MsgSwapExactAmountIn{
					TokenIn: sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100)),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.ATOM,
						},
					},
					Sender: authtypes.NewModuleAddress("random").String(),
				}, lastSwapIndex+1)
				suite.app.AmmKeeper.SetLastSwapRequestIndex(suite.ctx, lastSwapIndex+1)
			},
			func() {
				list := suite.app.AmmKeeper.GetAllSwapExactAmountInRequests(suite.ctx)
				suite.Require().Equal(1, len(list))
			},
		},
		{
			"GetAllSwapExactAmountOutRequests",
			func() {
				suite.ResetSuite()

				lastSwapIndex := suite.app.AmmKeeper.GetLastSwapRequestIndex(suite.ctx)
				suite.app.AmmKeeper.SetSwapExactAmountOutRequests(suite.ctx, &types.MsgSwapExactAmountOut{
					TokenOut: sdk.NewCoin(ptypes.ATOM, math.NewInt(100)),
					Routes: []types.SwapAmountOutRoute{
						{
							PoolId:       1,
							TokenInDenom: ptypes.BaseCurrency,
						},
					},
					Sender: authtypes.NewModuleAddress("random").String(),
				}, lastSwapIndex+1)
				suite.app.AmmKeeper.SetLastSwapRequestIndex(suite.ctx, lastSwapIndex+1)
			},
			func() {
				list := suite.app.AmmKeeper.GetAllSwapExactAmountOutRequests(suite.ctx)
				suite.Require().Equal(1, len(list))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
