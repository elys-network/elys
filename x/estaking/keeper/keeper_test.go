package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	simapp "github.com/elys-network/elys/app"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	ckeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	initChain = true
)

type EstakingKeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
	genAccount  sdk.AccAddress
	valAddr     sdk.ValAddress
}

func (k *EstakingKeeperTestSuite) SetupTest() {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount()

	k.legacyAmino = app.LegacyAmino()
	k.ctx = app.BaseApp.NewContext(initChain, tmproto.Header{})
	k.app = app
	k.genAccount = genAccount
	k.valAddr = valAddr
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(EstakingKeeperTestSuite))
}

func (suite *EstakingKeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

// Add testing commitments
func (suite *EstakingKeeperTestSuite) AddTestCommitment(committed sdk.Coins) {
	commitment := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, suite.genAccount)

	for _, c := range committed {
		commitment.AddCommittedTokens(c.Denom, c.Amount, uint64(suite.ctx.BlockTime().Unix()))
	}

	suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitment)
}

// Add testing claimed
func (suite *EstakingKeeperTestSuite) AddTestClaimed(committed sdk.Coins) {
	commitment := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, suite.genAccount)
	commitment.Claimed = commitment.Claimed.Add(committed...)
	suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitment)
}

// Set asset profile
func (suite *EstakingKeeperTestSuite) SetAssetProfile() {
	// Set assetprofile entry for denom
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, aptypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true, WithdrawEnabled: true})
}

// Prepare unclaimed tokens
func (suite *EstakingKeeperTestSuite) PrepareUnclaimedTokens() sdk.Coins {
	unclaimed := sdk.Coins{}
	unclaimed = unclaimed.Add(sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000)))
	unclaimed = unclaimed.Add(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(20000)))

	// Mint coins
	err := suite.app.BankKeeper.MintCoins(suite.ctx, ctypes.ModuleName, unclaimed)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ctypes.ModuleName, suite.genAccount, unclaimed)
	suite.Require().NoError(err)

	return unclaimed
}

// Prepare committed tokens
func (suite *EstakingKeeperTestSuite) PrepareCommittedTokens() sdk.Coins {
	committed := sdk.Coins{}
	committed = committed.Add(sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000)))
	committed = committed.Add(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(5000)))

	// Mint coins
	err := suite.app.BankKeeper.MintCoins(suite.ctx, ctypes.ModuleName, committed)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ctypes.ModuleName, suite.genAccount, committed)
	suite.Require().NoError(err)

	commitment := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, suite.genAccount)
	commitment.Claimed = commitment.Claimed.Add(committed...)
	suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitment)

	return committed
}

func TestEstakingExtendedFunctions(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	stakingKeeper := app.StakingKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	valAddr := validators[0].GetOperator()
	delegations := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	totalBonded := estakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, totalBonded.String(), "1000000")

	// set commitments
	commitments := app.CommitmentKeeper.GetCommitments(ctx, addr)
	commitments.AddClaimed(sdk.NewInt64Coin(ptypes.Eden, 1000_000))
	app.CommitmentKeeper.SetCommitments(ctx, commitments)
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{
		BaseDenom:       ptypes.Eden,
		Denom:           ptypes.Eden,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	commitmentMsgServer := ckeeper.NewMsgServerImpl(app.CommitmentKeeper)
	_, err := commitmentMsgServer.CommitClaimedRewards(sdk.WrapSDKContext(ctx), &ctypes.MsgCommitClaimedRewards{
		Creator: addr.String(),
		Denom:   ptypes.Eden,
		Amount:  sdk.NewInt(1000_000),
	})
	require.Nil(t, err)

	totalBonded = estakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, totalBonded.String(), "2000000")

	edenVal := estakingKeeper.GetEdenValidator(ctx)
	require.Equal(t, edenVal.GetMoniker(), "EdenValidator")

	edenBVal := estakingKeeper.GetEdenBValidator(ctx)
	require.Equal(t, edenBVal.GetMoniker(), "EdenBValidator")

	require.Equal(t, estakingKeeper.Validator(ctx, edenVal.GetOperator()), edenVal)
	require.Equal(t, estakingKeeper.Validator(ctx, edenBVal.GetOperator()), edenBVal)

	edenDel := estakingKeeper.Delegation(ctx, addr, edenVal.GetOperator())
	require.Equal(t, edenDel.GetShares(), sdk.NewDec(1000_000))

	edenBDel := estakingKeeper.Delegation(ctx, addr, edenBVal.GetOperator())
	require.Nil(t, edenBDel)

	numDelegations := int64(0)
	estakingKeeper.IterateDelegations(ctx, addr, func(index int64, delegation stypes.DelegationI) (stop bool) {
		numDelegations++
		return false
	})
	require.Equal(t, numDelegations, int64(2))

	numBondedValidators := int64(0)
	estakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, delegation stypes.ValidatorI) (stop bool) {
		numBondedValidators++
		return false
	})
	require.Equal(t, numBondedValidators, int64(2))
}
