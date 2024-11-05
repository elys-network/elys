package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestBeginBlocker() {
	tests := []struct {
		name           string
		blockHeight    int64
		epochLength    int64
		epochPosition  int64
		interestRate   sdk.Dec
		redemptionRate sdk.Dec
		expectedError  error
	}{
		{
			name:           "epoch has passed",
			epochLength:    1,
			epochPosition:  0,
			interestRate:   sdk.MustNewDecFromStr("0.17"),
			redemptionRate: sdk.ZeroDec(),
			expectedError:  nil,
		},
		{
			name:           "epoch has not passed",
			epochLength:    2,
			epochPosition:  1,
			interestRate:   sdk.NewDec(5),
			redemptionRate: sdk.NewDec(10),
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			params := suite.app.StablestakeKeeper.GetParams(suite.ctx)
			params.InterestRate = tt.interestRate
			params.RedemptionRate = tt.redemptionRate
			params.EpochLength = tt.epochLength
			params.TotalValue = sdk.NewInt(1000000)
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
