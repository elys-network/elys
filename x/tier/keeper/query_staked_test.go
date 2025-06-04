package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	stablestakekeeper "github.com/elys-network/elys/v6/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/v6/x/stablestake/types"
	"github.com/elys-network/elys/v6/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryStakedInvalidRequest() {
	k := suite.app.TierKeeper
	_, err := k.Staked(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryStaked() {

	sender := suite.account

	fee := math.LegacyMustNewDecFromStr("0.0002")
	issueAmount := math.NewInt(10_000_000_000_000_000)
	coins := sdk.NewCoins(
		sdk.NewCoin(ptypes.ATOM, issueAmount.MulRaw(100)),
		sdk.NewCoin(ptypes.Elys, issueAmount.MulRaw(100)),
		sdk.NewCoin(ptypes.BaseCurrency, issueAmount.MulRaw(100)),
	)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, coins)
	if err != nil {
		panic(err)
	}
	msgCreatePool := ammtypes.MsgCreatePool{
		Sender: sender.String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:   fee,
			UseOracle: true,
			FeeDenom:  ptypes.BaseCurrency,
		},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewCoin(ptypes.BaseCurrency, issueAmount),
				Weight: math.NewInt(50),
			},
			{
				Token:  sdk.NewCoin(ptypes.ATOM, issueAmount),
				Weight: math.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	if err != nil {
		panic(err)
	}

	msgServer := stablestakekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)

	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.EnabledPools = []uint64{poolId}
	err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)

	//STAKE USDC
	_, err = msgServer.Bond(suite.ctx, &stablestaketypes.MsgBond{
		Creator: sender.String(),
		Amount:  math.NewInt(100000000),
		PoolId:  1,
	})
	suite.Require().NoError(err)

	//TESTING STAKED FUNCTION.
	k := suite.app.TierKeeper
	r, err := k.Staked(suite.ctx, &types.QueryStakedRequest{
		User: sender.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(r.Commitments.TruncateInt(), math.NewInt(100))
}
