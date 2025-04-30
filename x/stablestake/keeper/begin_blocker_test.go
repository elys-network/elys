package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/elys-network/elys/x/stablestake/types"
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
			name:           "begin blocker call",
			interestRate:   sdkmath.LegacyMustNewDecFromStr("0.17"),
			redemptionRate: sdkmath.LegacyZeroDec(),
			expectedError:  nil,
		},
		{
			name:           "begin blocker call with more height",
			interestRate:   sdkmath.LegacyMustNewDecFromStr("0.17"),
			redemptionRate: sdkmath.LegacyZeroDec(),
			expectedError:  nil,
			blockHeight:    25768000,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.app.StablestakeKeeper.SetPool(suite.ctx, types.Pool{
				DepositDenom:         "uusdc",
				InterestRate:         tt.interestRate,
				InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.17"),
				InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.12"),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyOneDec(),
				NetAmount:            sdkmath.NewInt(1000000),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				Id:                   1,
			})

			suite.ctx = suite.ctx.WithBlockHeight(tt.blockHeight).WithBlockTime(time.Now())

			suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)

			storedPool, _ := suite.app.StablestakeKeeper.GetPool(suite.ctx, 1)
			require.Equal(suite.T(), tt.interestRate, storedPool.InterestRate)
		})
	}
}
