package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestUpdateParams() {
	sender := authtypes.NewModuleAddress(govtypes.ModuleName)
	tests := []struct {
		name      string
		authority string
		params    types.Params
		expected  error
	}{
		{
			name:      "valid authority",
			authority: sender.String(),
			params: types.Params{
				DepositDenom:         "stake",
				RedemptionRate:       sdkmath.LegacyNewDec(1),
				EpochLength:          100,
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				TotalValue:           sdkmath.OneInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
			},
			expected: nil,
		},
		{
			name:      "invalid authority",
			authority: "invalid_authority",
			params: types.Params{
				DepositDenom:         "stake",
				RedemptionRate:       sdkmath.LegacyNewDec(1),
				EpochLength:          100,
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.1"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyMustNewDecFromStr("0.01"),
				TotalValue:           sdkmath.OneInt(),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.1"),
			},
			expected: errors.New("invalid authority"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			msgServer := keeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)

			msg := &types.MsgUpdateParams{
				Authority: tt.authority,
				Params:    &tt.params,
			}

			_, err := msgServer.UpdateParams(suite.ctx, msg)
			if tt.expected != nil {
				require.ErrorContains(suite.T(), err, tt.expected.Error())
			} else {
				require.NoError(suite.T(), err)
				storedParams := suite.app.StablestakeKeeper.GetParams(suite.ctx)
				require.Equal(suite.T(), tt.params.TotalValue, storedParams.TotalValue)
			}
		})
	}
}
