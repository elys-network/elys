package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	simapp "github.com/elys-network/elys/v5/app"
	commitmentkeeper "github.com/elys-network/elys/v5/x/commitment/keeper"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestClaimRewardProgram(t *testing.T) {
	// Setup test environment
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(false)
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create test address
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))[0]
	addrStr := addr.String()

	// Test cases
	tests := []struct {
		name          string
		setup         func()
		msg           *types.MsgClaimRewardProgram
		expectedError error
		checkResponse func(*types.MsgClaimRewardProgramResponse)
	}{
		{
			name: "successful claim",
			setup: func() {
				// Set reward program
				rewardProgram := types.RewardProgram{
					Address: addrStr,
					Amount:  sdkmath.NewInt(1000),
					Claimed: false,
				}
				keeper.SetRewardProgram(ctx, rewardProgram)

				// Set params
				params := types.Params{
					EnableClaim:                   true,
					StartRewardProgramClaimHeight: 1,
					EndRewardProgramClaimHeight:   100,
				}
				keeper.SetParams(ctx, params)

				// Set block height
				ctx = ctx.WithBlockHeight(50)
			},
			msg: &types.MsgClaimRewardProgram{
				ClaimAddress: addrStr,
			},
			expectedError: nil,
			checkResponse: func(resp *types.MsgClaimRewardProgramResponse) {
				require.Equal(t, sdkmath.NewInt(1000), resp.EdenAmount)
			},
		},
		{
			name: "reward program not found",
			setup: func() {
				// Set empty reward program
				rewardProgram := types.RewardProgram{
					Address: addrStr,
					Amount:  sdkmath.ZeroInt(),
					Claimed: false,
				}
				keeper.SetRewardProgram(ctx, rewardProgram)
			},
			msg: &types.MsgClaimRewardProgram{
				ClaimAddress: addrStr,
			},
			expectedError: types.ErrRewardProgramNotFound,
		},
		{
			name: "claim not enabled",
			setup: func() {
				// Set reward program
				rewardProgram := types.RewardProgram{
					Address: addrStr,
					Amount:  sdkmath.NewInt(1000),
					Claimed: false,
				}
				keeper.SetRewardProgram(ctx, rewardProgram)

				// Set params with claim disabled
				params := types.Params{
					EnableClaim:                   false,
					StartRewardProgramClaimHeight: 1,
					EndRewardProgramClaimHeight:   100,
				}
				keeper.SetParams(ctx, params)
			},
			msg: &types.MsgClaimRewardProgram{
				ClaimAddress: addrStr,
			},
			expectedError: types.ErrClaimNotEnabled,
		},
		{
			name: "already claimed",
			setup: func() {
				// Set reward program as claimed
				rewardProgram := types.RewardProgram{
					Address: addrStr,
					Amount:  sdkmath.NewInt(1000),
					Claimed: true,
				}
				// Set params with claim disabled
				params := types.Params{
					EnableClaim:                   true,
					StartRewardProgramClaimHeight: 1,
					EndRewardProgramClaimHeight:   100,
				}
				keeper.SetParams(ctx, params)
				keeper.SetRewardProgram(ctx, rewardProgram)
			},
			msg: &types.MsgClaimRewardProgram{
				ClaimAddress: addrStr,
			},
			expectedError: types.ErrRewardProgramAlreadyClaimed,
		},
		{
			name: "claim not started",
			setup: func() {
				// Set reward program
				rewardProgram := types.RewardProgram{
					Address: addrStr,
					Amount:  sdkmath.NewInt(1000),
					Claimed: false,
				}
				keeper.SetRewardProgram(ctx, rewardProgram)

				// Set params with future start height
				params := types.Params{
					EnableClaim:                   true,
					StartRewardProgramClaimHeight: 100,
					EndRewardProgramClaimHeight:   200,
				}
				keeper.SetParams(ctx, params)

				// Set current block height before start
				ctx = ctx.WithBlockHeight(50)
			},
			msg: &types.MsgClaimRewardProgram{
				ClaimAddress: addrStr,
			},
			expectedError: types.ErrRewardProgramNotStarted,
		},
		{
			name: "claim ended",
			setup: func() {
				// Set reward program
				rewardProgram := types.RewardProgram{
					Address: addrStr,
					Amount:  sdkmath.NewInt(1000),
					Claimed: false,
				}
				keeper.SetRewardProgram(ctx, rewardProgram)

				// Set params with past end height
				params := types.Params{
					EnableClaim:                   true,
					StartRewardProgramClaimHeight: 1,
					EndRewardProgramClaimHeight:   50,
				}
				keeper.SetParams(ctx, params)

				// Set current block height after end
				ctx = ctx.WithBlockHeight(100)
			},
			msg: &types.MsgClaimRewardProgram{
				ClaimAddress: addrStr,
			},
			expectedError: types.ErrRewardProgramEnded,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test case
			tc.setup()

			// Execute claim
			resp, err := msgServer.ClaimRewardProgram(ctx, tc.msg)

			// Check error
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				return
			}

			// Check success case
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Check response
			if tc.checkResponse != nil {
				tc.checkResponse(resp)
			}

			// Verify reward program is marked as claimed
			rewardProgram := keeper.GetRewardProgram(ctx, addr)
			require.True(t, rewardProgram.Claimed)

			// Verify total claimed amount
			total := keeper.GetTotalRewardProgramClaimed(ctx)
			require.Equal(t, sdkmath.NewInt(1000), total.TotalEdenClaimed)

			// Verify balance
			balance := app.CommitmentKeeper.GetAllBalances(ctx, addr)
			require.Equal(t, sdkmath.NewInt(1000), balance.AmountOf(ptypes.Eden))
		})
	}
}
