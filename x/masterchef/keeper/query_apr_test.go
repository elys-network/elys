package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	"github.com/elys-network/elys/x/masterchef/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SetupApp(t *testing.T) (keeper.Keeper, sdk.Context) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle, estaking := app.MasterchefKeeper, app.AmmKeeper, app.OracleKeeper, app.EstakingKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(100010))

	// Create a pool
	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], coins)
	require.NoError(t, err)

	//app.StakingKeeper.Delegate(ctx, addr[0], sdk.NewInt(100000), sdk.Unbonded, sdk.NewDec(0.1), sdk.NewInt(100000))

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdk.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdk.NewInt(1000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdk.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100)),
	})

	poolParams := &ammtypes.PoolParams{
		SwapFee:                     sdk.ZeroDec(),
		ExitFee:                     sdk.ZeroDec(),
		UseOracle:                   false,
		WeightBreakingFeeMultiplier: sdk.ZeroDec(),
		WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
		ExternalLiquidityRatio:      sdk.OneDec(),
		WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   sdk.ZeroDec(),
		FeeDenom:                    "",
	}

	// Create a Elys+USDC pool
	msgServer := ammkeeper.NewMsgServerImpl(amm)
	resp, err := msgServer.CreatePool(
		sdk.WrapSDKContext(ctx),
		&ammtypes.MsgCreatePool{
			Sender:     addr[0].String(),
			PoolParams: poolParams,
			PoolAssets: poolAssets,
		})

	require.NoError(t, err)
	require.Equal(t, resp.PoolID, uint64(1))

	poolInfo, found := mk.GetPoolInfo(ctx, resp.PoolID)
	require.True(t, found)

	poolInfo.DexApr = sdk.NewDecWithPrec(1, 2)  // 1%
	poolInfo.EdenApr = sdk.NewDecWithPrec(2, 2) // 2%
	mk.SetPoolInfo(ctx, poolInfo)
	estakingParams := estaking.GetParams(ctx)
	estakingParams.StakeIncentives =
		&estakingtypes.IncentiveInfo{
			EdenAmountPerYear: sdk.NewInt(1000000),
			BlocksDistributed: 1000000,
		}
	estakingParams.MaxEdenRewardAprStakers = sdk.NewDecWithPrec(30, 2)
	estaking.SetParams(ctx, estakingParams)

	mkParams := mk.GetParams(ctx)
	mkParams.LpIncentives = &types.IncentiveInfo{
		EdenAmountPerYear: sdk.NewInt(1000000000),
		BlocksDistributed: 1000000,
	}
	mk.SetParams(ctx, mkParams)

	return mk, ctx
}

func TestApr(t *testing.T) {
	tests := []struct {
		desc     string
		request  *types.QueryAprRequest
		response *types.QueryAprResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryAprRequest{
				WithdrawType: 0,
				Denom:        "ueden",
			},
			response: &types.QueryAprResponse{
				Apr: sdk.MustNewDecFromStr("0.299999999999999995"),
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	mk, ctx := SetupApp(t)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := mk.Apr(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestAprs(t *testing.T) {
	tests := []struct {
		desc     string
		request  *types.QueryAprsRequest
		response *types.QueryAprsResponse
		err      error
	}{
		{
			desc:    "valid request",
			request: &types.QueryAprsRequest{},
			response: &types.QueryAprsResponse{
				UsdcAprUsdc:  sdk.ZeroDec(),
				EdenAprUsdc:  sdk.ZeroDec(),
				UsdcAprEdenb: sdk.ZeroDec(),
				EdenAprEdenb: sdk.MustNewDecFromStr("0.299999999999999995"),
				UsdcAprEden:  sdk.ZeroDec(),
				EdenAprEden:  sdk.MustNewDecFromStr("0.299999999999999995"),
				EdenbAprEden: sdk.OneDec(),
				UsdcAprElys:  sdk.ZeroDec(),
				EdenAprElys:  sdk.MustNewDecFromStr("0.299999999999999995"),
				EdenbAprElys: sdk.OneDec(),
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	mk, ctx := SetupApp(t)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := mk.Aprs(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response.String(), response.String())
			}
		})
	}
}

func TestPoolRewards(t *testing.T) {
	tests := []struct {
		desc     string
		request  *types.QueryPoolRewardsRequest
		response *types.QueryPoolRewardsResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryPoolRewardsRequest{
				PoolIds: []uint64{1},
			},
			response: &types.QueryPoolRewardsResponse{
				Pools: []types.PoolRewards{{
					PoolId:      1,
					RewardsUsd:  sdk.NewDec(420),
					RewardCoins: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdk.NewInt(200)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(400))},
					EdenForward: sdk.NewCoin(ptypes.Eden, sdk.NewInt(0)),
				},
				},
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	mk, ctx := SetupApp(t)

	ctx.BlockTime()
	mk.SetPoolRewardsAccum(ctx, types.PoolRewardsAccum{
		PoolId: 1, BlockHeight: ctx.BlockHeight(),
		DexReward: sdk.NewDec(100), EdenReward: sdk.NewDec(200),
		Timestamp: uint64(ctx.BlockTime().Unix()), GasReward: sdk.NewDec(300),
	})

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := mk.PoolRewards(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response.String(), response.String())
			}
		})
	}
}
