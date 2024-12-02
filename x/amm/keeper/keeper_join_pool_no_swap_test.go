package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestJoinPoolNoSwap() {
	// Define test cases
	testCases := []struct {
		name           string
		setup          func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins)
		expectedErrMsg string
	}{
		{
			"pool does not exist",
			func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins) {
				return sdk.AccAddress([]byte("sender")), 1, sdkmath.NewInt(100), sdk.Coins{}
			},
			"invalid pool id",
		},
		{
			"successful join pool No oracle",
			func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins) {
				suite.SetupCoinPrices()
				addr := suite.AddAccounts(1, nil)
				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr[0], false, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				return addr[0], 1, pool.TotalShares.Amount, sdk.Coins{sdk.NewCoin(types.ATOM, amount.MulRaw(10)), sdk.NewCoin(types.BaseCurrency, amount.MulRaw(10))}
			},
			"",
		},
		{
			"successful join pool with oracle",
			func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins) {
				suite.SetupCoinPrices()
				addr := suite.AddAccounts(1, nil)
				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr[0], true, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				return addr[0], 2, pool.TotalShares.Amount, sdk.Coins{sdk.NewCoin(types.ATOM, amount.MulRaw(10))}
			},
			"",
		},
		{
			"Needed LpLiquidity is more than tokenInMaxs",
			func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins) {
				addr := suite.AddAccounts(1, nil)
				amount := sdkmath.NewInt(100)
				share, _ := sdkmath.NewIntFromString("20000000000000000000000000000000")
				return addr[0], 1, share, sdk.Coins{sdk.NewCoin(types.ATOM, amount), sdk.NewCoin(types.BaseCurrency, amount)}
			},
			"TokenInMaxs is less than the needed LP liquidity to this JoinPoolNoSwap",
		},
		{
			"tokenInMaxs does not contain Needed LpLiquidity coins",
			func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins) {
				addr := suite.AddAccounts(1, nil)
				amount := sdkmath.NewInt(100)
				share, _ := sdkmath.NewIntFromString("20000000000000000000000000000000")
				return addr[0], 1, share, sdk.Coins{sdk.NewCoin("nocoin", amount), sdk.NewCoin(types.BaseCurrency, amount)}
			},
			"TokenInMaxs does not include all the tokens that are part of the target pool",
		},
		{
			"tokenInMaxs does not contain Needed LpLiquidity coins",
			func() (sdk.AccAddress, uint64, sdkmath.Int, sdk.Coins) {
				addr := suite.AddAccounts(1, nil)
				amount := sdkmath.NewInt(100)
				share, _ := sdkmath.NewIntFromString("20000000000000000000000000000000")
				return addr[0], 1, share, sdk.Coins{sdk.NewCoin("nocoin", amount), sdk.NewCoin(types.ATOM, amount), sdk.NewCoin(types.BaseCurrency, amount)}
			},
			"TokenInMaxs includes tokens that are not part of the target pool",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			sender, poolId, shareOutAmount, tokenInMaxs := tc.setup()
			tokenIn, sharesOut, err := suite.app.AmmKeeper.JoinPoolNoSwap(suite.ctx, sender, poolId, shareOutAmount, tokenInMaxs)
			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().True(tokenIn.IsAllLTE(tokenInMaxs))
				suite.Require().True(sharesOut.LTE(shareOutAmount))
			}
		})
	}
}
