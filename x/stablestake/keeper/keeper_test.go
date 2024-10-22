package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	simapp "github.com/elys-network/elys/app"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/suite"
)

const (
	initChain = true
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain, suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
	suite.SetStakingParam()
	suite.SetStableStakeParam()
	suite.SetupAssetProfile()
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetStakingParam() error {
	return suite.app.StakingKeeper.SetParams(suite.ctx, stakingtypes.Params{
		UnbondingTime:     1209600,
		MaxValidators:     60,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		BondDenom:         "uelys",
		MinCommissionRate: math.LegacyNewDec(0),
	})
}

func (suite *KeeperTestSuite) SetStableStakeParam() error {

	params := stablestaketypes.DefaultParams()
	suite.app.StablestakeKeeper.SetParams(suite.ctx, params)
	return nil
}

func (suite *KeeperTestSuite) SetupAssetProfile() {

	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, atypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "uusdc",
		Path:                     "transfer/channel-12",
		IbcChannelId:             "channel-12",
		IbcCounterpartyChannelId: "channel-19",
		DisplayName:              "USDC",
		DisplaySymbol:            "uUSDC",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "uUSDC",
		TransferLimit:            "",
		Permissions:              []string{},
		UnitDenom:                "uusdc",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	})
}
