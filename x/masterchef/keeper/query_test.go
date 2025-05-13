package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *MasterchefKeeperTestSuite) SetupApp() {

	// Generate 1 random account with 1000stake balanced
	addr := authtypes.NewModuleAddress(govtypes.ModuleName)

	// Create a pool
	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 110000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	suite.MintMultipleTokenToAddress(addr, coins)

	//app.StakingKeeper.Delegate(ctx, addr[0], math.NewInt(100000), sdk.Unbonded, sdk.NewDec(0.1), math.NewInt(100000))

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100)),
	})

	poolParams := ammtypes.PoolParams{
		SwapFee:   sdkmath.LegacyZeroDec(),
		UseOracle: false,
		FeeDenom:  ptypes.BaseCurrency,
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(addr, poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	poolInfo, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, ammPool.PoolId)
	suite.Require().True(found)

	poolInfo.DexApr = sdkmath.LegacyNewDecWithPrec(1, 2)  // 1%
	poolInfo.EdenApr = sdkmath.LegacyNewDecWithPrec(2, 2) // 2%
	suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, poolInfo)
	estakingParams := suite.app.EstakingKeeper.GetParams(suite.ctx)
	estakingParams.StakeIncentives =
		&estakingtypes.IncentiveInfo{
			EdenAmountPerYear: sdkmath.NewInt(1000000),
			BlocksDistributed: 1000000,
		}
	estakingParams.MaxEdenRewardAprStakers = sdkmath.LegacyNewDecWithPrec(30, 2)
	suite.app.EstakingKeeper.SetParams(suite.ctx, estakingParams)

	mkParams := suite.app.MasterchefKeeper.GetParams(suite.ctx)
	mkParams.LpIncentives = &types.IncentiveInfo{
		EdenAmountPerYear: sdkmath.NewInt(1000000000),
		BlocksDistributed: 1000000,
	}
	suite.app.MasterchefKeeper.SetParams(suite.ctx, mkParams)

	suite.app.StablestakeKeeper.SetPool(suite.ctx, stablestaketypes.Pool{
		InterestRate:         sdkmath.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      sdkmath.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      sdkmath.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     sdkmath.LegacyOneDec(),
		NetAmount:            sdkmath.ZeroInt(),
		MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
		Id:                   stablestaketypes.UsdcPoolId,
		DepositDenom:         ptypes.BaseCurrency,
	})
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
				Apr: sdkmath.LegacyMustNewDecFromStr("0.299999999999999999"),
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
				UsdcAprUsdc:  sdkmath.LegacyZeroDec(),
				EdenAprUsdc:  sdkmath.LegacyZeroDec(),
				UsdcAprEdenb: sdkmath.LegacyZeroDec(),
				EdenAprEdenb: sdkmath.LegacyMustNewDecFromStr("0.299999999999999999"),
				UsdcAprEden:  sdkmath.LegacyZeroDec(),
				EdenAprEden:  sdkmath.LegacyMustNewDecFromStr("0.299999999999999999"),
				EdenbAprEden: sdkmath.LegacyOneDec(),
				UsdcAprElys:  sdkmath.LegacyZeroDec(),
				EdenAprElys:  sdkmath.LegacyMustNewDecFromStr("0.299999999999999999"),
				EdenbAprElys: sdkmath.LegacyOneDec(),
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
					PoolId:        1,
					RewardsUsd:    sdkmath.LegacyNewDec(420),
					RewardCoins:   sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(200)), sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(400))},
					EdenForward:   sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(0)),
					RewardsUsdApr: sdkmath.LegacyMustNewDecFromStr("1531.468531468531468531"),
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
		DexReward: sdkmath.LegacyNewDec(100), EdenReward: sdkmath.LegacyNewDec(200),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()), GasReward: sdkmath.LegacyNewDec(300),
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

func (suite *MasterchefKeeperTestSuite) TestExternalIncentiveQuery() {
	tests := []struct {
		desc     string
		request  *types.QueryExternalIncentiveRequest
		response *types.QueryExternalIncentiveResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryExternalIncentiveRequest{
				Id: 1,
			},
			response: &types.QueryExternalIncentiveResponse{
				ExternalIncentive: types.ExternalIncentive{
					Id:             1,
					RewardDenom:    "reward2",
					PoolId:         1,
					FromBlock:      suite.ctx.BlockHeight() - 1,
					ToBlock:        suite.ctx.BlockHeight() + 101,
					AmountPerBlock: sdkmath.OneInt(),
					Apr:            sdkmath.LegacyZeroDec(),
				},
			},
			err: nil,
		},
	}

	suite.SetupApp()

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			suite.app.MasterchefKeeper.SetExternalIncentive(suite.ctx, tc.response.ExternalIncentive)
			response, err := suite.app.MasterchefKeeper.ExternalIncentive(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.String(), response.String())
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestPoolInfo() {
	tests := []struct {
		desc     string
		request  *types.QueryPoolInfoRequest
		response *types.QueryPoolInfoResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryPoolInfoRequest{
				PoolId: 1,
			},
			response: &types.QueryPoolInfoResponse{
				PoolInfo: types.PoolInfo{
					PoolId:               1,
					RewardWallet:         "cosmos1lz2ajk0mvhda7hdzedydeany3f673600pd6euqnjvqv8w5p4az5qmt08tn",
					Multiplier:           sdkmath.LegacyOneDec(),
					EdenApr:              sdkmath.LegacyMustNewDecFromStr("0.02"),
					GasApr:               sdkmath.LegacyZeroDec(),
					DexApr:               sdkmath.LegacyMustNewDecFromStr("0.01"),
					ExternalIncentiveApr: sdkmath.LegacyZeroDec(),
				},
			},
			err: nil,
		},
	}

	suite.SetupApp()

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			response, err := suite.app.MasterchefKeeper.PoolInfo(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.String(), response.String())
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestPoolRewardInfoQuery() {
	suite.SetupApp()
	tests := []struct {
		desc     string
		request  *types.QueryPoolRewardInfoRequest
		response *types.QueryPoolRewardInfoResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryPoolRewardInfoRequest{
				PoolId:      1,
				RewardDenom: "reward",
			},
			response: &types.QueryPoolRewardInfoResponse{
				PoolRewardInfo: types.PoolRewardInfo{
					PoolId:                1,
					RewardDenom:           "reward",
					PoolAccRewardPerShare: sdkmath.LegacyNewDec(0),
					LastUpdatedBlock:      uint64(suite.ctx.BlockHeight()),
				},
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			suite.app.MasterchefKeeper.EndBlocker(suite.ctx)
			suite.app.MasterchefKeeper.SetPoolRewardInfo(suite.ctx, tc.response.PoolRewardInfo)
			response, err := suite.app.MasterchefKeeper.PoolRewardInfo(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.String(), response.String())
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestUserRewardInfoQuery() {
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdkmath.NewInt(100010))
	tests := []struct {
		desc     string
		request  *types.QueryUserRewardInfoRequest
		response *types.QueryUserRewardInfoResponse
		err      error
	}{
		{
			desc: "valid request",
			request: &types.QueryUserRewardInfoRequest{
				PoolId:      1,
				User:        addr[0].String(),
				RewardDenom: "reward",
			},
			response: &types.QueryUserRewardInfoResponse{
				UserRewardInfo: types.UserRewardInfo{
					User:          addr[0].String(),
					PoolId:        1,
					RewardDenom:   "reward",
					RewardDebt:    sdkmath.LegacyOneDec(),
					RewardPending: sdkmath.LegacyOneDec(),
				},
			},
			err: nil,
		},
	}

	suite.SetupApp()

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			suite.app.MasterchefKeeper.SetUserRewardInfo(suite.ctx, tc.response.UserRewardInfo)
			response, err := suite.app.MasterchefKeeper.UserRewardInfo(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.String(), response.String())
			}
		})
	}
}

func (suite *MasterchefKeeperTestSuite) TestTotalPendingRewards() {
	suite.SetupApp()

	// Create test addresses
	addr1 := sdk.AccAddress("test1_______________")
	addr2 := sdk.AccAddress("test2_______________")
	addr3 := sdk.AccAddress("test3_______________")

	// masterchef address
	// elys1nwc45a3fl0dz37m5ulvw8pmpfjnewhgz7t96zn

	// Create test user reward info entries
	userRewardInfos := []types.UserRewardInfo{
		{
			User:          addr1.String(),
			PoolId:        1,
			RewardDenom:   "ueden",
			RewardPending: sdkmath.LegacyNewDec(100),
		},
		{
			User:          addr2.String(),
			PoolId:        2,
			RewardDenom:   "uusdc",
			RewardPending: sdkmath.LegacyNewDec(200),
		},
		{
			User:          addr3.String(),
			PoolId:        1,
			RewardDenom:   "ueden",
			RewardPending: sdkmath.LegacyNewDec(300),
		},
	}

	// Set all user reward info entries
	for _, info := range userRewardInfos {
		suite.app.MasterchefKeeper.SetUserRewardInfo(suite.ctx, info)
	}

	tests := []struct {
		desc     string
		request  *types.QueryTotalPendingRewardsRequest
		response *types.QueryTotalPendingRewardsResponse
		err      error
	}{
		{
			desc: "valid request with default pagination",
			request: &types.QueryTotalPendingRewardsRequest{
				Pagination: &query.PageRequest{
					Limit: 5000,
				},
			},
			response: &types.QueryTotalPendingRewardsResponse{
				TotalPendingRewards: sdk.NewCoins(
					sdk.NewCoin("ueden", sdkmath.NewInt(400)), // 100 + 300
					sdk.NewCoin("uusdc", sdkmath.NewInt(200)),
				),
				Count: 3,
			},
			err: nil,
		},
		{
			desc: "valid request with custom pagination limit",
			request: &types.QueryTotalPendingRewardsRequest{
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			response: &types.QueryTotalPendingRewardsResponse{
				TotalPendingRewards: sdk.NewCoins(
					sdk.NewCoin("ueden", sdkmath.NewInt(100)),
				),
				Count: 1,
			},
			err: nil,
		},
		{
			desc:    "invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	for _, tc := range tests {
		suite.Run(tc.desc, func() {
			response, err := suite.app.MasterchefKeeper.TotalPendingRewards(suite.ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.response.TotalPendingRewards.String(), response.TotalPendingRewards.String())
				suite.Require().Equal(tc.response.Count, response.Count)
			}
		})
	}
}
