package keeper_test

import (
	"errors"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestUpdateParams() {
	sender := authtypes.NewModuleAddress(govtypes.ModuleName)
	tests := []struct {
		name      string
		authority string
		params    types.Params
		expected  error
	}{
		{
			name:      "invalid authority",
			authority: "invalid_authority",
			params:    types.Params{},
			expected:  errors.New("invalid authority"),
		},
		{
			name:      "valid authority",
			authority: sender.String(),
			params: types.Params{
				MarketOrderEnabled:   true,
				StakeEnabled:         true,
				ProcessOrdersEnabled: true,
				SwapEnabled:          true,
				PerpetualEnabled:     true,
				RewardEnabled:        true,
				LeverageEnabled:      true,
				LimitProcessOrder:    100,
				RewardPercentage:     math.LegacyMustNewDecFromStr("0.1"),
				MarginError:          math.LegacyMustNewDecFromStr("0.05"),
				MinimumDeposit:       math.NewInt(100000),
				Tolerance:            math.LegacyMustNewDecFromStr("0.05"),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			msgServer := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)

			msg := &types.MsgUpdateParams{
				Authority: tt.authority,
				Params:    &tt.params,
			}

			_, err := msgServer.UpdateParams(suite.ctx, msg)
			if tt.expected != nil {
				suite.Require().ErrorContains(err, tt.expected.Error())
			} else {
				suite.Require().NoError(err)
				storedParams := suite.app.TradeshieldKeeper.GetParams(suite.ctx)
				suite.Require().Equal(tt.params, storedParams)
			}
		})
	}
}
