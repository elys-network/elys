package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	simapp "github.com/elys-network/elys/v6/app"
	aptypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ckeeper "github.com/elys-network/elys/v6/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/v6/x/commitment/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
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

func (suite *EstakingKeeperTestSuite) SetupTest() {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount(suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
	suite.genAccount = genAccount
	suite.valAddr = valAddr
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
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, aptypes.Entry{
		BaseDenom:       ptypes.Eden,
		Denom:           ptypes.Eden,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, aptypes.Entry{
		BaseDenom:       ptypes.BaseCurrency,
		Denom:           ptypes.BaseCurrency,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, aptypes.Entry{
		BaseDenom:       ptypes.EdenB,
		Denom:           ptypes.EdenB,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
}

// Prepare unclaimed tokens
func (suite *EstakingKeeperTestSuite) PrepareUnclaimedTokens() sdk.Coins {
	unclaimed := sdk.Coins{}
	unclaimed = unclaimed.Add(sdk.NewCoin(ptypes.Eden, math.NewInt(2000)))
	unclaimed = unclaimed.Add(sdk.NewCoin(ptypes.EdenB, math.NewInt(20000)))

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
	committed = committed.Add(sdk.NewCoin(ptypes.Eden, math.NewInt(10000)))
	committed = committed.Add(sdk.NewCoin(ptypes.EdenB, math.NewInt(5000)))

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

func (suite *EstakingKeeperTestSuite) GetAccountIssueAmount() math.Int {
	return math.NewInt(10_000_000_000_000)
}

func (suite *EstakingKeeperTestSuite) AddAccounts(n int, given []sdk.AccAddress) []sdk.AccAddress {
	issueAmount := suite.GetAccountIssueAmount()
	var addresses []sdk.AccAddress
	if n > len(given) {
		addresses = simapp.AddTestAddrs(suite.app, suite.ctx, n-len(given), issueAmount)
		addresses = append(addresses, given...)
	} else {
		addresses = given
	}
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount),
			sdk.NewCoin(ptypes.Elys, issueAmount),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
		if err != nil {
			panic(err)
		}
	}
	return addresses
}

func (suite *EstakingKeeperTestSuite) TestEstakingExtendedFunctions() {
	// create validator with 50% commission
	validators, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	suite.Require().Nil(err)
	suite.Require().True(len(validators) > 0)
	valAddr := validators[0].GetOperator()

	operatorAddr, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().Nil(err)

	delegations, err := suite.app.StakingKeeper.GetValidatorDelegations(suite.ctx, operatorAddr)
	suite.Require().Nil(err)

	suite.Require().True(len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	totalBonded, err := suite.app.EstakingKeeper.TotalBondedTokens(suite.ctx)
	suite.Require().Nil(err)
	suite.Require().Equal(totalBonded.String(), "1000000")

	// set commitments
	commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, addr)
	commitments.AddClaimed(sdk.NewInt64Coin(ptypes.Eden, 1000_000))
	suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitments)

	suite.SetAssetProfile()

	commitmentMsgServer := ckeeper.NewMsgServerImpl(*suite.app.CommitmentKeeper)
	_, err = commitmentMsgServer.CommitClaimedRewards(suite.ctx, &ctypes.MsgCommitClaimedRewards{
		Creator: addr.String(),
		Denom:   ptypes.Eden,
		Amount:  math.NewInt(1000_000),
	})
	suite.Require().Nil(err)

	totalBonded, err = suite.app.EstakingKeeper.TotalBondedTokens(suite.ctx)
	suite.Require().Nil(err)
	suite.Require().Equal(totalBonded.String(), "2000000")

	edenVal := suite.app.EstakingKeeper.GetEdenValidator(suite.ctx)
	suite.Require().Equal(edenVal.GetMoniker(), "EdenValidator")

	edenBVal := suite.app.EstakingKeeper.GetEdenBValidator(suite.ctx)
	suite.Require().Equal(edenBVal.GetMoniker(), "EdenBValidator")

	operatorValAddr, err := sdk.ValAddressFromBech32(edenVal.GetOperator())
	suite.Require().Nil(err)

	operatorBValAddr, err := sdk.ValAddressFromBech32(edenBVal.GetOperator())
	suite.Require().Nil(err)

	validator, err := suite.app.EstakingKeeper.Validator(suite.ctx, operatorValAddr)
	suite.Require().Nil(err)
	suite.Require().Equal(validator, edenVal)
	validator, err = suite.app.EstakingKeeper.Validator(suite.ctx, operatorBValAddr)
	suite.Require().Nil(err)
	suite.Require().Equal(validator, edenBVal)

	edenDel, err := suite.app.EstakingKeeper.Delegation(suite.ctx, addr, operatorValAddr)
	suite.Require().Nil(err)
	suite.Require().Equal(edenDel.GetShares(), math.LegacyNewDec(1000_000))

	edenBDel, err := suite.app.EstakingKeeper.Delegation(suite.ctx, addr, operatorBValAddr)
	suite.Require().Nil(err)
	suite.Require().Nil(edenBDel)

	numDelegations := int64(0)
	suite.app.EstakingKeeper.IterateDelegations(suite.ctx, addr, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
		numDelegations++
		return false
	})
	suite.Require().Equal(numDelegations, int64(2))

	numBondedValidators := int64(0)
	suite.app.EstakingKeeper.IterateBondedValidatorsByPower(suite.ctx, func(index int64, delegation stakingtypes.ValidatorI) (stop bool) {
		numBondedValidators++
		return false
	})
	suite.Require().Equal(numBondedValidators, int64(2))

	// test IterateValidators
	numValidators := int64(0)
	suite.app.EstakingKeeper.IterateValidators(suite.ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		numValidators++
		return false
	})
	suite.Require().Equal(numValidators, int64(3))

	// test Slash
	edenValConsAddr, err := edenVal.GetConsAddr()
	suite.Require().Nil(err)
	suite.app.EstakingKeeper.Slash(suite.ctx, edenValConsAddr, 0, 0, math.LegacyNewDecWithPrec(5, 1))

	// test SlashWithInfractionReason
	suite.app.EstakingKeeper.SlashWithInfractionReason(suite.ctx, edenValConsAddr, 0, 0, math.LegacyNewDecWithPrec(5, 1), stakingtypes.Infraction_INFRACTION_UNSPECIFIED)

	// test WithdrawEdenBReward
	err = suite.app.EstakingKeeper.WithdrawEdenBReward(suite.ctx, addr)
	suite.Require().Error(err)

	// test WithdrawEdenReward
	err = suite.app.EstakingKeeper.WithdrawEdenReward(suite.ctx, addr)
	suite.Require().Nil(err)

	// test DelegationRewards
	rewards, err := suite.app.EstakingKeeper.DelegationRewards(suite.ctx, edenDel.GetDelegatorAddr(), valAddr)
	suite.Require().Nil(err)
	suite.Require().Nil(rewards)

}
