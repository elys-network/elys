package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/testutil/sample"
	"github.com/elys-network/elys/v7/x/stablestake/keeper"
	"github.com/elys-network/elys/v7/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestMsgServerAddPool() {
	for _, tc := range []struct {
		desc        string
		pool        types.Pool
		expPass     bool
		addMultiple bool
		sender      string
	}{
		{
			desc: "successful add pool",
			pool: types.Pool{
				DepositDenom:         "stake",
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				NetAmount:            sdkmath.ZeroInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				MaxWithdrawRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				Id:                   2,
			},
			expPass:     true,
			sender:      "",
			addMultiple: false,
		},
		{
			desc: "fails to add pool for same denom",
			pool: types.Pool{
				DepositDenom:         "stake",
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				NetAmount:            sdkmath.ZeroInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				MaxWithdrawRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				Id:                   2,
			},
			sender:      "",
			expPass:     false,
			addMultiple: true,
		},
		{
			desc: "failed to add pool with invalid creatoe",
			pool: types.Pool{
				DepositDenom:         "stake",
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				NetAmount:            sdkmath.ZeroInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				MaxWithdrawRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				Id:                   2,
			},
			sender:      "invalid_address",
			expPass:     false,
			addMultiple: false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			authority := authtypes.NewModuleAddress(govtypes.ModuleName)
			if tc.sender != "" {
				authority = sdk.MustAccAddressFromBech32(sample.AccAddress())
			}
			msgServer := keeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
			_, err := msgServer.AddPool(
				suite.ctx,
				&types.MsgAddPool{
					Sender:               authority.String(),
					DepositDenom:         tc.pool.DepositDenom,
					InterestRate:         tc.pool.InterestRate,
					InterestRateMax:      tc.pool.InterestRateMax,
					InterestRateMin:      tc.pool.InterestRateMin,
					InterestRateIncrease: tc.pool.InterestRateIncrease,
					InterestRateDecrease: tc.pool.InterestRateDecrease,
					HealthGainFactor:     tc.pool.HealthGainFactor,
					MaxLeverageRatio:     tc.pool.MaxLeverageRatio,
					MaxWithdrawRatio:     tc.pool.MaxWithdrawRatio,
				})

			if tc.addMultiple {
				_, err = msgServer.AddPool(
					suite.ctx,
					&types.MsgAddPool{
						Sender:               authority.String(),
						DepositDenom:         tc.pool.DepositDenom,
						InterestRate:         tc.pool.InterestRate,
						InterestRateMax:      tc.pool.InterestRateMax,
						InterestRateMin:      tc.pool.InterestRateMin,
						InterestRateIncrease: tc.pool.InterestRateIncrease,
						InterestRateDecrease: tc.pool.InterestRateDecrease,
						HealthGainFactor:     tc.pool.HealthGainFactor,
						MaxLeverageRatio:     tc.pool.MaxLeverageRatio,
						MaxWithdrawRatio:     tc.pool.MaxWithdrawRatio,
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
