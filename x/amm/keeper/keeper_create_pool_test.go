package keeper_test

import (
	"github.com/elys-network/elys/v5/x/amm/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *AmmKeeperTestSuite) TestCreatePool() {
	// Define test cases
	testCases := []struct {
		name           string
		setup          func() *types.MsgCreatePool
		expectedErrMsg string
	}{
		{
			"asset profile not found",
			func() *types.MsgCreatePool {
				addr := suite.AddAccounts(1, nil)
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
				return &types.MsgCreatePool{
					Sender:     addr[0].String(),
					PoolParams: types.PoolParams{},
					PoolAssets: []types.PoolAsset{},
				}
			},
			"asset profile not found for denom",
		},
		{
			"Balance pool Create Error",
			func() *types.MsgCreatePool {
				suite.ResetSuite()
				addr := suite.AddAccounts(1, nil)
				return &types.MsgCreatePool{
					Sender:     addr[0].String(),
					PoolParams: types.PoolParams{},
					PoolAssets: []types.PoolAsset{},
				}
			},
			"swap_fee is nil",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, msg)
			if tc.expectedErrMsg != "" {
				require.Error(suite.T(), err)
				require.Contains(suite.T(), err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(suite.T(), err)
				require.NotZero(suite.T(), poolId)
			}
		})
	}
}
