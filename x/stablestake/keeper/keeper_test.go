package keeper_test

import (
	"github.com/elys-network/elys/v6/x/amm/types"
	"testing"

	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	atypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/v6/x/stablestake/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v6/app"
	"github.com/stretchr/testify/assert"
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

	poolAddr := types.NewPoolAddress(uint64(1))
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	pool := types.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: types.PoolParams{
			SwapFee:   math.LegacyZeroDec(),
			UseOracle: true,
			FeeDenom:  ptypes.BaseCurrency,
		},
		TotalShares: sdk.NewCoin(types.GetPoolShareDenom(1), math.ZeroInt()),
		PoolAssets: []types.PoolAsset{
			{
				Token:  sdk.NewCoin("uusdc", math.NewInt(1000_000)),
				Weight: math.NewInt(10),
			},
			{
				Token:  sdk.NewCoin("uatom", math.NewInt(1000_000)),
				Weight: math.NewInt(10),
			},
		},
		TotalWeight: math.ZeroInt(),
	}
	suite.app.AmmKeeper.SetPool(suite.ctx, pool)

	leverageLpParams := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	leverageLpParams.EnabledPools = []uint64{1}
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &leverageLpParams)
	suite.Require().NoError(err)

	suite.app.StablestakeKeeper.SetPool(suite.ctx, stablestaketypes.Pool{
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		NetAmount:            math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		Id:                   1,
		DepositDenom:         ptypes.BaseCurrency,
	})

	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusdc",
		Display: "USDC",
		Decimal: 6,
	})
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     math.LegacyOneDec(),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// TestKeeper_Logger tests the Logger function
func TestKeeper_Logger(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.StablestakeKeeper

	logger := app.Logger()

	keeper.Logger(ctx).Info("test")
	logger.Info("test")
}

// TestKeeper_SetHooks_Panic tests the SetHooks function with a nil argument
func TestKeeper_SetHooks_Panic(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	keeper := app.StablestakeKeeper

	assert.Panics(t, func() {
		keeper.SetHooks(nil)
	})
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
