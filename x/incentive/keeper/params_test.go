package keeper_test

import (
	"strconv"
	"strings"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	testkeeper "github.com/elys-network/elys/testutil/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.IncentiveKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

func TestUpdatePoolMultiplierInfo(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, amm, oracle, bk, ck := app.IncentiveKeeper, app.AmmKeeper, app.OracleKeeper, app.BankKeeper, app.CommitmentKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(100010))

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000)))

	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

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
		ExternalLiquidityRatio:      sdk.OneDec(),
		LpFeePortion:                sdk.ZeroDec(),
		StakingFeePortion:           sdk.ZeroDec(),
		WeightRecoveryFeePortion:    sdk.ZeroDec(),
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

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check balance change on sender
	balances := bk.GetBalance(ctx, addr[0], "amm/pool/1")
	require.Equal(t, balances, sdk.NewCoin("amm/pool/1", sdk.NewInt(0)))

	// check lp token commitment
	commitments, found := ck.GetCommitments(ctx, addr[0].String())
	require.True(t, found)
	require.Len(t, commitments.CommittedTokens, 1)
	require.Equal(t, commitments.CommittedTokens[0].Denom, "amm/pool/1")
	require.Equal(t, commitments.CommittedTokens[0].Amount.String(), "100100000000000000000")

	poolIds := strings.Split("1,2,3", ",")
	multipliers := strings.Split("5,1,1", ",")
	require.Len(t, poolIds, 3)

	poolInfo, found := ik.GetPoolInfo(ctx, resp.PoolID)
	require.True(t, found)
	require.Equal(t, poolInfo.Multiplier, sdk.NewDec(1))

	poolMultipliers := make([]types.PoolMultipliers, 0)
	for i := range poolIds {
		poolId, err := strconv.ParseUint(poolIds[i], 10, 64)
		require.NoError(t, err)

		multiplier, err := sdk.NewDecFromStr(multipliers[i])
		require.NoError(t, err)

		poolMultiplier := types.PoolMultipliers{
			PoolId:     poolId,
			Multiplier: multiplier,
		}

		poolMultipliers = append(poolMultipliers, poolMultiplier)
	}

	ik.UpdatePoolMultipliers(ctx, poolMultipliers)

	poolInfo, found = ik.GetPoolInfo(ctx, resp.PoolID)
	require.True(t, found)

	// After setting up, it should become 5.
	require.Equal(t, poolInfo.Multiplier, sdk.NewDec(5))
}
