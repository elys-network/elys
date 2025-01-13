package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestMsgServerAddPool() {
	for _, tc := range []struct {
		desc        string
		pool        types.Pool
		expPass     bool
		addMultiple bool
	}{
		{
			desc: "successful add pool",
			pool: types.Pool{
				DepositDenom:         "stake",
				RedemptionRate:       sdkmath.LegacyZeroDec(),
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				TotalValue:           sdkmath.ZeroInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				PoolId:               2,
			},
			expPass:     true,
			addMultiple: false,
		},
		{
			desc: "failed add pool for same denom",
			pool: types.Pool{
				DepositDenom:         "stake",
				RedemptionRate:       sdkmath.LegacyZeroDec(),
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				TotalValue:           sdkmath.ZeroInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				PoolId:               2,
			},
			expPass:     false,
			addMultiple: true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			authority := authtypes.NewModuleAddress(govtypes.ModuleName)
			msgServer := keeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
			_, err := msgServer.AddPool(
				suite.ctx,
				&types.MsgAddPool{
					Authority:            authority.String(),
					DepositDenom:         tc.pool.DepositDenom,
					InterestRate:         tc.pool.InterestRate,
					InterestRateMax:      tc.pool.InterestRateMax,
					InterestRateMin:      tc.pool.InterestRateMin,
					InterestRateIncrease: tc.pool.InterestRateIncrease,
					InterestRateDecrease: tc.pool.InterestRateDecrease,
					HealthGainFactor:     tc.pool.HealthGainFactor,
					MaxLeverageRatio:     tc.pool.MaxLeverageRatio,
				})

			if tc.addMultiple {
				_, err = msgServer.AddPool(
					suite.ctx,
					&types.MsgAddPool{
						Authority:            authority.String(),
						DepositDenom:         tc.pool.DepositDenom,
						InterestRate:         tc.pool.InterestRate,
						InterestRateMax:      tc.pool.InterestRateMax,
						InterestRateMin:      tc.pool.InterestRateMin,
						InterestRateIncrease: tc.pool.InterestRateIncrease,
						InterestRateDecrease: tc.pool.InterestRateDecrease,
						HealthGainFactor:     tc.pool.HealthGainFactor,
						MaxLeverageRatio:     tc.pool.MaxLeverageRatio,
					})
			}
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				pools := suite.app.StablestakeKeeper.GetAllPools(suite.ctx)
				suite.Require().Len(pools, 2)
				suite.Require().Equal(tc.pool, pools[1])
			}
		})
	}
}
