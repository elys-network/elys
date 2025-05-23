package keeper_test

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestExitPool() {
	testCases := []struct {
		name           string
		setup          func() (sdk.AccAddress, uint64, math.Int, sdk.Coins, string, bool)
		expectedErrMsg string
	}{
		{
			"pool does not exist",
			func() (sdk.AccAddress, uint64, math.Int, sdk.Coins, string, bool) {
				addr := suite.AddAccounts(1, nil)
				return addr[0], 1, math.NewInt(100), sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))), "uatom", false
			},
			"invalid pool id",
		},
		{
			"exiting more shares than available",
			func() (sdk.AccAddress, uint64, math.Int, sdk.Coins, string, bool) {
				suite.SetupCoinPrices()
				addr := suite.AddAccounts(1, nil)
				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr[0], true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				return addr[0], 1, pool.TotalShares.Amount.Add(sdkmath.NewInt(10)), sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))), "uatom", false
			},
			"Trying to exit >= the number of shares contained in the pool",
		},
		{
			"exiting negative shares",
			func() (sdk.AccAddress, uint64, math.Int, sdk.Coins, string, bool) {
				addr := suite.AddAccounts(1, nil)
				return addr[0], 1, sdkmath.NewInt(1).Neg(), sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))), "uatom", false
			},
			"Trying to exit a negative amount of shares",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			exiter, poolId, inShares, minTokensOut, tokenOutDenom, isLiq := tc.setup()
			_, _, _, _, _, err := suite.app.AmmKeeper.ExitPool(suite.ctx, exiter, poolId, inShares, minTokensOut, tokenOutDenom, isLiq, true)
			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
