package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *MasterchefKeeperTestSuite) SetupApp() {

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 1, math.NewInt(100010))

	// Create a pool
	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	suite.MintMultipleTokenToAddress(addr[0], coins)

	//app.StakingKeeper.Delegate(ctx, addr[0], math.NewInt(100000), sdk.Unbonded, sdk.NewDec(0.1), math.NewInt(100000))

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: math.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, math.NewInt(1000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: math.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100)),
	})

	poolParams := &ammtypes.PoolParams{
		SwapFee:                     math.LegacyZeroDec(),
		ExitFee:                     math.LegacyZeroDec(),
		UseOracle:                   false,
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
		FeeDenom:                    "",
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(addr[0], poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	poolInfo, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, ammPool.PoolId)
	suite.Require().True(found)

	poolInfo.DexApr = math.LegacyNewDecWithPrec(1, 2)  // 1%
	poolInfo.EdenApr = math.LegacyNewDecWithPrec(2, 2) // 2%
	suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, poolInfo)
	estakingParams := suite.app.EstakingKeeper.GetParams(suite.ctx)
	estakingParams.StakeIncentives =
		&estakingtypes.IncentiveInfo{
			EdenAmountPerYear: math.NewInt(1000000),
			BlocksDistributed: 1000000,
		}
	estakingParams.MaxEdenRewardAprStakers = math.LegacyNewDecWithPrec(30, 2)
	suite.app.EstakingKeeper.SetParams(suite.ctx, estakingParams)

	mkParams := suite.app.MasterchefKeeper.GetParams(suite.ctx)
	mkParams.LpIncentives = &types.IncentiveInfo{
		EdenAmountPerYear: math.NewInt(1000000000),
		BlocksDistributed: 1000000,
	}
	suite.app.MasterchefKeeper.SetParams(suite.ctx, mkParams)
}

func (suite *MasterchefKeeperTestSuite) TestApr() {
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
				Apr: math.LegacyMustNewDecFromStr("0.299999999999999995"),
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	suite.SetupApp()

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			response, err := suite.app.MasterchefKeeper.Apr(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response, response)
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestAprs() {
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
				UsdcAprUsdc:  math.LegacyZeroDec(),
				EdenAprUsdc:  math.LegacyZeroDec(),
				UsdcAprEdenb: math.LegacyZeroDec(),
				EdenAprEdenb: math.LegacyMustNewDecFromStr("0.299999999999999995"),
				UsdcAprEden:  math.LegacyZeroDec(),
				EdenAprEden:  math.LegacyMustNewDecFromStr("0.299999999999999995"),
				EdenbAprEden: math.LegacyOneDec(),
				UsdcAprElys:  math.LegacyZeroDec(),
				EdenAprElys:  math.LegacyMustNewDecFromStr("0.299999999999999995"),
				EdenbAprElys: math.LegacyOneDec(),
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	suite.SetupApp()

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			response, err := suite.app.MasterchefKeeper.Aprs(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.String(), response.String())
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestPoolRewards() {
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
					RewardsUsd:  math.LegacyNewDec(420),
					RewardCoins: sdk.Coins{sdk.NewCoin(ptypes.Eden, math.NewInt(200)), sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(400))},
					EdenForward: sdk.NewCoin(ptypes.Eden, math.NewInt(0)),
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

	suite.SetupApp()

	suite.ctx.BlockTime()
	suite.app.MasterchefKeeper.SetPoolRewardsAccum(suite.ctx, types.PoolRewardsAccum{
		PoolId: 1, BlockHeight: suite.ctx.BlockHeight(),
		DexReward: math.LegacyNewDec(100), EdenReward: math.LegacyNewDec(200),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()), GasReward: math.LegacyNewDec(300),
	})

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			response, err := suite.app.MasterchefKeeper.PoolRewards(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.String(), response.String())
			}
		})
	}
}
