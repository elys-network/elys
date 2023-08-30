package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func SetupStableCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.USDC,
		Display: "USDC",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.USDT,
		Display: "USDT",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     sdk.NewDec(100),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdk.NewDec(100),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
}

func TestCalculateRewardsForLPs(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, amm, oracle, bk, ck := app.IncentiveKeeper, app.AmmKeeper, app.OracleKeeper, app.BankKeeper, app.CommitmentKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(100010))

	var committed []sdk.Coins
	var uncommitted []sdk.Coins

	// Prepare uncommitted tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin("ueden", sdk.NewInt(2000)))
	uncommitted = append(uncommitted, uedenToken)

	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenToken)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoins(sdk.NewCoin("ueden", sdk.NewInt(500)))
	lpToken1 := sdk.NewCoins(sdk.NewCoin("lp-elys-usdc", sdk.NewInt(500)))
	lpToken2 := sdk.NewCoins(sdk.NewCoin("lp-ueden-usdc", sdk.NewInt(2000)))
	committed = append(committed, uedenToken)
	committed = append(committed, lpToken1)
	committed = append(committed, lpToken2)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, lpToken1)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], lpToken1)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, lpToken2)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], lpToken2)
	require.NoError(t, err)

	simapp.AddTestCommitment(app, ctx, addr[0], committed, uncommitted)

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.USDC, sdk.NewInt(100000)))

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
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
		Token:  sdk.NewCoin(ptypes.USDC, sdk.NewInt(100)),
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
	require.Equal(t, resp.PoolID, uint64(0))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check balance change on sender
	balances := bk.GetBalance(ctx, addr[0], "amm/pool/0")
	require.Equal(t, balances, sdk.NewCoin("amm/pool/0", sdk.NewInt(0)))

	// check lp token commitment
	commitments, found := ck.GetCommitments(ctx, addr[0].String())
	require.True(t, found)
	require.Len(t, commitments.CommittedTokens, 4)
	require.Equal(t, commitments.CommittedTokens[3].Denom, "amm/pool/0")
	require.Equal(t, commitments.CommittedTokens[3].Amount.String(), "100100000000000000000")

	require.Equal(t, ik.CalculateTVL(ctx), sdk.NewDecWithPrec(1001, 1))

	edenAmountPerEpochLp := sdk.NewInt(1000000)
	totalProxyTVL := ik.CalculateProxyTVL(ctx)

	// Recalculate total committed info
	ik.UpdateTotalCommitmentInfo(ctx)

	gasFeesLPsAmt := sdk.NewDec(1000)
	// Calculate rewards for LPs
	newUncommittedEdenTokensLp, dexRewardsLp := ik.CalculateRewardsForLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLp, gasFeesLPsAmt)

	require.Equal(t, newUncommittedEdenTokensLp, sdk.NewInt(1000000))
	require.Equal(t, dexRewardsLp, sdk.NewInt(1000))
}
