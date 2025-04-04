package keeper_test

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestUpdatePool() {
	sender := authtypes.NewModuleAddress(govtypes.ModuleName)
	tests := []struct {
		name      string
		authority string
		pool      types.Pool
		expected  error
	}{
		{
			name:      "valid authority",
			authority: sender.String(),
			pool: types.Pool{
				DepositDenom:         "stake",
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				TotalValue:           sdkmath.OneInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				Id:                   1,
			},
			expected: nil,
		},
		{
			name:      "invalid authority",
			authority: "invalid_authority",
			pool: types.Pool{
				DepositDenom:         "stake",
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				TotalValue:           sdkmath.OneInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
				Id:                   1,
			},
			expected: errors.New("invalid authority"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			msgServer := keeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)

			msg := &types.MsgUpdatePool{
				Authority:            tt.authority,
				InterestRateMax:      tt.pool.InterestRateMax,
				InterestRateMin:      tt.pool.InterestRateMin,
				InterestRateIncrease: tt.pool.InterestRateIncrease,
				InterestRateDecrease: tt.pool.InterestRateDecrease,
				HealthGainFactor:     tt.pool.HealthGainFactor,
				MaxLeverageRatio:     tt.pool.MaxLeverageRatio,
				PoolId:               1,
			}

			_, err := msgServer.UpdatePool(suite.ctx, msg)
			if tt.expected != nil {
				require.ErrorContains(suite.T(), err, tt.expected.Error())
			} else {
				require.NoError(suite.T(), err)
				storedPool, _ := suite.app.StablestakeKeeper.GetPool(suite.ctx, 1)
				require.Equal(suite.T(), tt.pool.MaxLeverageRatio, storedPool.MaxLeverageRatio)
			}
		})
	}
}
