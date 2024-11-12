package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"time"

	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestBeginBlocker() {
	tests := []struct {
		name           string
		blockHeight    int64
		epochLength    int64
		epochPosition  int64
		interestRate   sdkmath.LegacyDec
		redemptionRate sdkmath.LegacyDec
		expectedError  error
	}{
		{
			name:           "epoch has passed",
			epochLength:    1,
			epochPosition:  0,
			interestRate:   sdkmath.LegacyMustNewDecFromStr("0.17"),
			redemptionRate: sdkmath.LegacyZeroDec(),
			expectedError:  nil,
		},
		{
			name:           "epoch has not passed",
			epochLength:    2,
			epochPosition:  1,
			interestRate:   sdkmath.LegacyNewDec(5),
			redemptionRate: sdkmath.LegacyNewDec(10),
			expectedError:  nil,
		},
		{
			name:           "delete old data",
			epochLength:    2,
			epochPosition:  1,
			interestRate:   sdkmath.LegacyNewDec(5),
			redemptionRate: sdkmath.LegacyNewDec(10),
			expectedError:  nil,
			blockHeight:    95768100,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			params := suite.app.StablestakeKeeper.GetParams(suite.ctx)
			params.InterestRate = tt.interestRate
			params.RedemptionRate = tt.redemptionRate
			params.EpochLength = tt.epochLength
			params.TotalValue = sdkmath.NewInt(1000000)
			suite.app.StablestakeKeeper.SetParams(suite.ctx, params)

			suite.ctx = suite.ctx.WithBlockHeight(tt.blockHeight).WithBlockTime(time.Now())

			suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)

			storedParams := suite.app.StablestakeKeeper.GetParams(suite.ctx)
			if tt.epochPosition == 0 {
				require.Equal(suite.T(), tt.interestRate, storedParams.InterestRate)
				require.Equal(suite.T(), tt.redemptionRate, storedParams.RedemptionRate)
			} else {
				require.NotEqual(suite.T(), tt.interestRate, storedParams.InterestRate)
				require.NotEqual(suite.T(), tt.redemptionRate, storedParams.RedemptionRate)
			}
		})
	}
}
