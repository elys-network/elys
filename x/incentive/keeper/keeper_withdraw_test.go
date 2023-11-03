package keeper_test

import (
	"math"
	"strings"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

// Test scenario
// We have one validator and delegator.
// In this testing, total delegation amount equals to the genesis delegator delegation amount.
// So the whole Eden amount will be for the genesis delegator
// And we apply the commission rate of the validator and calculate his commission amount
func TestGiveCommissionToValidators(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	delegator := genAccount.String()
	// Calculate delegated amount per delegator
	delegatedAmt := ik.CalculateDelegatedAmount(ctx, delegator)
	newUncommittedEdenTokens := sdk.NewInt(10000)
	dexRewardsByStakers := sdk.NewDec(1000)
	// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
	edenCommissionGiven, dexRewardsCommissionGiven := ik.GiveCommissionToValidators(ctx, delegator, delegatedAmt, newUncommittedEdenTokens, dexRewardsByStakers)

	require.Equal(t, edenCommissionGiven, sdk.NewInt(500))
	require.Equal(t, dexRewardsCommissionGiven, sdk.NewInt(50))
}

func TestProcessWithdrawRewards(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	// Generate 2 random accounts with 10000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	var committed []sdk.Coins
	var uncommitted []sdk.Coins

	// Prepare uncommitted tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000)))
	uedenBToken := sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(2000)))
	lpToken := sdk.NewCoins(sdk.NewCoin("lp-elys-usdc", sdk.NewInt(500)))
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(500)))

	uncommitted = append(uncommitted, uedenToken)
	uncommitted = append(uncommitted, uedenBToken)
	uncommitted = append(uncommitted, lpToken)
	uncommitted = append(uncommitted, usdcToken)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.BaseCurrency, CommitEnabled: false, WithdrawEnabled: true})
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true, WithdrawEnabled: true})
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.EdenB, CommitEnabled: true, WithdrawEnabled: true})
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: "lp-elys-usdc", CommitEnabled: true, WithdrawEnabled: false})

	// Prepare committed tokens
	uedenTokenC := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(1500)))
	uedenBTokenC := sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(500)))
	committed = append(committed, uedenTokenC)
	committed = append(committed, uedenBTokenC)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, addr[0], committed, uncommitted)
	simapp.AddTestCommitment(app, ctx, addr[1], committed, uncommitted)
	_, found := app.CommitmentKeeper.GetCommitments(ctx, addr[0].String())
	require.True(t, found)

	_, found = app.CommitmentKeeper.GetCommitments(ctx, addr[1].String())
	require.True(t, found)

	// Get dex revenue wallet
	dexRewardUSDC := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(5000)))

	// Mint 5000 usdc
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, dexRewardUSDC)
	require.NoError(t, err)

	// Transfer 5000 USDC to dex revenue wallet
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, ctypes.ModuleName, simapp.DexRevenueCollectorName, dexRewardUSDC)
	require.NoError(t, err)

	// Withdraw rewards
	err = ik.ProcessWithdrawRewards(ctx, addr[0].String(), ptypes.Eden)
	require.NoError(t, err)

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], ptypes.Eden)
	require.Equal(t, sdk.Coins{edenCoin}, uedenToken)

	// Withdraw rewards
	err = ik.ProcessWithdrawRewards(ctx, addr[0].String(), ptypes.BaseCurrency)
	require.NoError(t, err)

	usdcCoin := app.BankKeeper.GetBalance(ctx, addr[0], ptypes.BaseCurrency)
	require.Equal(t, sdk.Coins{usdcCoin}, usdcToken)
}

func TestProcessWithdrawValidatorCommission(t *testing.T) {
	app, genAccount, valAddress := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	// Get dex revenue wallet
	dexRewardUSDC := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(5000)))

	// Mint 5000 usdc
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, dexRewardUSDC)
	require.NoError(t, err)

	// Transfer 5000 USDC to dex revenue wallet
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, ctypes.ModuleName, simapp.DexRevenueCollectorName, dexRewardUSDC)
	require.NoError(t, err)

	delegator := genAccount.String()
	// Calculate delegated amount per delegator
	delegatedAmt := ik.CalculateDelegatedAmount(ctx, delegator)
	newUncommittedEdenTokens := sdk.NewInt(10000)
	dexRewardsByStakers := sdk.NewDec(1000)

	// Create an entity in commitment module for validator
	app.CommitmentKeeper.StandardStakingToken(ctx, delegator, valAddress.String(), ptypes.Eden)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.BaseCurrency, CommitEnabled: false, WithdrawEnabled: true})
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true, WithdrawEnabled: true})
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.EdenB, CommitEnabled: true, WithdrawEnabled: true})

	// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
	edenCommissionGiven, dexRewardsCommissionGiven := ik.GiveCommissionToValidators(ctx, delegator, delegatedAmt, newUncommittedEdenTokens, dexRewardsByStakers)

	_, found := app.CommitmentKeeper.GetCommitments(ctx, valAddress.String())
	require.True(t, found)

	require.Equal(t, edenCommissionGiven, sdk.NewInt(500))
	require.Equal(t, dexRewardsCommissionGiven, sdk.NewInt(50))

	delAddr, err := sdk.AccAddressFromBech32(delegator)
	require.NoError(t, err)

	found = false
	// Get all delegations
	delegations := app.StakingKeeper.GetDelegatorDelegations(ctx, delAddr, math.MaxUint16)
	for _, del := range delegations {
		// Get validator address
		valAddr := del.GetValidatorAddr()

		// If it is not requested by the validator creator
		if strings.EqualFold(valAddress.String(), valAddr.String()) {
			found = true
			break
		}
	}

	require.True(t, found)
	err = ik.ProcessWithdrawValidatorCommission(ctx, delegator, valAddress.String(), ptypes.Eden)
	require.NoError(t, err)

	edenCoin := app.BankKeeper.GetBalance(ctx, genAccount, ptypes.Eden)
	require.Equal(t, edenCoin, sdk.NewCoin(ptypes.Eden, edenCommissionGiven))
}
