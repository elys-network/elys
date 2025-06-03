package keeper_test

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/clob/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
	"time"
)

const (
	AvgBlockTime = 5
	BaseDenom    = "uatom"
	QuoteDenom   = "uusdc"
	MarketId     = uint64(1)
)

var (
	assetProfileAtom = assetprofiletypes.Entry{
		BaseDenom:                BaseDenom,
		Decimals:                 6,
		Denom:                    "uatom",
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "ATOM",
		DisplaySymbol:            "ATOM",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                "uatom",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	oracleProfileAtom = oracletypes.AssetInfo{
		Denom:      "uatom",
		Display:    "ATOM",
		BandTicker: "ATOM",
		ElysTicker: "ATOM",
		Decimal:    6,
	}
	assetProfileUsdc = assetprofiletypes.Entry{
		BaseDenom:                QuoteDenom,
		Decimals:                 6,
		Denom:                    QuoteDenom,
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "USDC",
		DisplaySymbol:            "USDC",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                QuoteDenom,
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	oracleProfileUsdc = oracletypes.AssetInfo{
		Denom:      QuoteDenom,
		Display:    "USDC",
		BandTicker: "USDC",
		ElysTicker: "USDC",
		Decimal:    6,
	}
	assetProfileOsmo = assetprofiletypes.Entry{
		BaseDenom:                "uosmo",
		Decimals:                 6,
		Denom:                    "uosmo",
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              "OSMO",
		DisplaySymbol:            "OSMO",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "",
		TransferLimit:            "",
		Permissions:              nil,
		UnitDenom:                "uosmo",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	}
	oracleProfileOsmo = oracletypes.AssetInfo{
		Denom:      "uosmo",
		Display:    "OSMO",
		BandTicker: "OSMO",
		ElysTicker: "OSMO",
		Decimal:    6,
	}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp

	avgBlockTime uint64
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(true, suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(true).WithBlockTime(time.Now())
	suite.app = app
	suite.avgBlockTime = AvgBlockTime

	oracleParams := app.OracleKeeper.GetParams(suite.ctx)
	oracleParams.LifeTimeInBlocks = 10000
	oracleParams.PriceExpiryTime = 84600
	app.OracleKeeper.SetParams(suite.ctx, oracleParams)

	suite.SetAssetProfiles()
	suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{math.LegacyNewDec(10), math.LegacyNewDec(1)})
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) SetAssetProfiles() {
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetProfileAtom)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetProfileUsdc)
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetProfileOsmo)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileAtom)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileUsdc)
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracleProfileOsmo)
}

func (suite *KeeperTestSuite) SetPrice(assets []string, prices []math.LegacyDec) {
	if len(assets) != len(prices) {
		panic("unequal lengths while setting prices during test")
	}
	for i, price := range prices {
		oraclePrice := oracletypes.Price{
			Asset:       assets[i],
			Price:       price,
			Source:      "test",
			Provider:    "test",
			Timestamp:   uint64(suite.ctx.BlockTime().Unix()),
			BlockHeight: uint64(suite.ctx.BlockHeight()),
		}
		suite.app.OracleKeeper.SetPrice(suite.ctx, oraclePrice)
	}
}

func (suite *KeeperTestSuite) IncreaseHeight(height uint64) {
	if height == 0 {
		panic("increment cannot be 0")
	}
	for i := uint64(1); i <= height; i++ {
		//_, err := suite.app.BeginBlocker(suite.ctx)
		//if err != nil {
		//	panic(err)
		//}
		//_, err = suite.app.EndBlocker(suite.ctx)
		//if err != nil {
		//	panic(err)
		//}
		currentHeight := suite.ctx.BlockHeight()
		currentTime := suite.ctx.BlockTime().Unix()
		ctx := suite.ctx.WithBlockHeight(currentHeight + 1)
		ctx = ctx.WithBlockTime(time.Unix(currentTime+int64(suite.avgBlockTime), 0))
		suite.ctx = ctx
	}
}

func (suite *KeeperTestSuite) SetupSubAccounts(total uint64, balance sdk.Coins) []types.SubAccount {
	if total == 0 {
		panic("total subaccounts cannot be 0")
	}

	all := suite.app.ClobKeeper.GetAllSubAccount(suite.ctx)
	var list []types.SubAccount
	for i := uint64(len(all) + 1); i <= uint64(len(all))+total; i++ {
		err := suite.app.BankKeeper.MintCoins(suite.ctx, ammtypes.ModuleName, balance)
		suite.Require().NoError(err)

		subAccountAddress := authtypes.NewModuleAddress("subAccount" + strconv.FormatUint(i, 10))

		subAccount := types.SubAccount{
			Owner:       subAccountAddress.String(),
			Id:          MarketId,
			TradeNounce: 0,
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ammtypes.ModuleName, subAccount.GetTradingAccountAddress(), balance)
		suite.Require().NoError(err)

		suite.app.ClobKeeper.SetSubAccount(suite.ctx, subAccount)
		list = append(list, subAccount)
	}
	return list
}

func (suite *KeeperTestSuite) CreateMarket(baseDenoms ...string) []types.PerpetualMarket {
	all := suite.app.ClobKeeper.GetAllPerpetualMarket(suite.ctx)
	var list []types.PerpetualMarket
	for i, baseDenom := range baseDenoms {
		if baseDenom == "" || baseDenom == QuoteDenom {
			panic("base Denom cannot be uusdc or empty")
		}
		market := types.PerpetualMarket{
			Id:                      uint64(len(all) + i + 1),
			BaseDenom:               baseDenom,
			QuoteDenom:              QuoteDenom,
			InitialMarginRatio:      IMR,
			MaintenanceMarginRatio:  math.LegacyMustNewDecFromStr("0.2"),
			Status:                  1,
			MaxAbsFundingRate:       math.LegacyMustNewDecFromStr("0.05"),
			MaxAbsFundingRateChange: math.LegacyMustNewDecFromStr("0.01"),
			TwapPricesWindow:        15,
			LiquidationFeeShareRate: math.LegacyMustNewDecFromStr("0.01"),
		}
		suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
		list = append(list, market)
	}
	return list
}

func (suite *KeeperTestSuite) GetAccountBalance(addr sdk.AccAddress, denom string) math.Int {
	return suite.app.BankKeeper.GetBalance(suite.ctx, addr, denom).Amount
}

func (suite *KeeperTestSuite) FundAccount(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) WithdrawFromAccount(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, addr, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) SetAccountBalance(addr sdk.AccAddress, coins sdk.Coins) {
	// Implementation depends on test setup - might need direct bank keeper state setting or burn+fund
	currentCoins := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr)
	err := suite.app.BankKeeper.SendCoins(suite.ctx, addr, authtypes.NewModuleAddress("dummy_burn"), currentCoins) // Send all current away
	if err != nil && !errors.Is(err, sdkerrors.ErrInsufficientFunds) {                                             // Ignore error if already empty
		suite.T().Logf("Error burning coins during SetAccountBalance: %v", err) // Log non-critical error
	}
	if !coins.IsZero() {
		suite.FundAccount(addr, coins) // Fund with desired amount
	}
}

func (suite *KeeperTestSuite) BurnAccountBalance(addr sdk.AccAddress, denom string) error {
	ctx := suite.ctx // Get context from suite

	// 1. Get the current balance for the specified denomination
	// Assumes GetBalance returns sdk.Coin{Denom: denom, Amount: math.Int}
	balanceCoin := suite.app.BankKeeper.GetBalance(ctx, addr, denom)
	balanceAmt := balanceCoin.Amount

	// 2. If balance is not positive, there's nothing to burn
	if balanceAmt.IsNil() || !balanceAmt.IsPositive() {
		return nil // Success, nothing to do
	}

	// 3. Define the burn address recipient
	// Common convention. Ensure this module account is initialized in your test app setup.
	// If not, consider using authtypes.FeeCollectorName or another known module address.
	burnAddr := authtypes.NewModuleAddress("burn")

	// 4. Create the sdk.Coins object containing the exact balance amount
	coinsToBurn := sdk.NewCoins(balanceCoin)

	// 5. Attempt to send the coins from the target address to the burn address
	err := suite.app.BankKeeper.SendCoins(ctx, addr, burnAddr, coinsToBurn)
	if err != nil {
		// Wrap the error from the bank keeper for better context
		return fmt.Errorf("failed to send %s from addr %s to burn addr %s: %w",
			coinsToBurn.String(), addr.String(), burnAddr.String(), err)
	}

	// 6. Optionally, verify the balance is now zero (can be useful for debugging tests)
	// finalBalance := suite.bankKeeper.GetBalance(ctx, addr, denom).Amount
	// if !finalBalance.IsZero() {
	// 	return fmt.Errorf("post-burn balance check failed for %s, expected zero, got %s", addr.String(), finalBalance)
	// }

	return nil // Success
}

func (suite *KeeperTestSuite) SetupExchangeTest() (market types.PerpetualMarket, buyerAcc types.SubAccount, sellerAcc types.SubAccount, marketAccAddr sdk.AccAddress) {
	suite.ResetSuite() // Or equivalent setup/teardown

	markets := suite.CreateMarket(BaseDenom)
	market = markets[0]
	marketAccAddr = market.GetAccount()
	// Ensure module account exists
	// suite.accountKeeper.GetModuleAccount(suite.ctx, suite.app.ClobKeeper.MarketModuleName)

	// Set initial funding rate (zero delta assumed for these tests' balance checks)
	initialFundingRate := types.FundingRate{MarketId: MarketId, Rate: math.LegacyMustNewDecFromStr("0.0001"), Block: uint64(suite.ctx.BlockHeight())}
	suite.app.ClobKeeper.SetFundingRate(suite.ctx, initialFundingRate)

	initialBalanceAmt := int64(200_000_000)
	initialMarketBalanceAmt := int64(500_000_000) // Ensure market has funds
	initialBalance := sdk.NewCoin(QuoteDenom, math.NewInt(initialBalanceAmt))
	initialMarketBalanceCoin := sdk.NewCoin(QuoteDenom, math.NewInt(initialMarketBalanceAmt))

	subAccounts := suite.SetupSubAccounts(2, sdk.NewCoins(initialBalance)) // Assumes this creates and funds accounts
	buyerAcc = subAccounts[0]
	sellerAcc = subAccounts[1]

	suite.FundAccount(marketAccAddr, sdk.NewCoins(initialMarketBalanceCoin)) // Fund market

	return market, buyerAcc, sellerAcc, marketAccAddr
}

// SetPerpetualStateWithEntryFR sets perpetual and owner mapping, applying current funding rate
func (suite *KeeperTestSuite) SetPerpetualStateWithEntryFR(p types.Perpetual, isCross bool) types.Perpetual {
	// Ensure EntryFundingRate matches current rate for test simplicity
	currentFundingRate := suite.app.ClobKeeper.GetFundingRate(suite.ctx, p.MarketId)
	p.EntryFundingRate = currentFundingRate.Rate
	// Assign ID if not set (useful for setup)
	if p.Id == 0 {
		p.Id = suite.app.ClobKeeper.GetAndIncrementPerpetualCounter(suite.ctx, p.MarketId)
	}
	suite.app.ClobKeeper.SetPerpetual(suite.ctx, p)
	subAccountId := p.MarketId
	if isCross {
		subAccountId = types.CrossMarginSubAccountId
	}
	suite.app.ClobKeeper.SetPerpetualOwner(suite.ctx, types.PerpetualOwner{
		Owner: p.Owner, SubAccountId: subAccountId, MarketId: p.MarketId, PerpetualId: p.Id,
	})
	return p // Return potentially updated perpetual (with ID)
}

// GetPerpetualState gets perpetual via owner mapping
func (suite *KeeperTestSuite) GetPerpetualState(ownerAddr sdk.AccAddress, marketId uint64) (types.Perpetual, bool) {
	subAccount := types.SubAccount{Owner: ownerAddr.String(), Id: marketId, TradeNounce: 0}
	ownerMapping, found := suite.app.ClobKeeper.CheckAndGetPerpetualOwner(suite.ctx, subAccount, marketId)
	if !found {
		return types.Perpetual{}, false
	}
	perp, err := suite.app.ClobKeeper.GetPerpetual(suite.ctx, marketId, ownerMapping.PerpetualId)
	if err != nil {
		// Handle specific 'not found' error if GetPerpetual returns one
		if errors.Is(err, types.ErrPerpetualNotFound) { // Assuming such an error exists
			// This case implies inconsistent state (owner mapping exists, perpetual doesn't)
			suite.T().Fatalf("Inconsistent state: PerpetualOwner found for %s/%d, but Perpetual %d not found: %v", ownerAddr.String(), marketId, ownerMapping.PerpetualId, err)
			return types.Perpetual{}, false // Should ideally not happen
		}
		// Fail test for other unexpected errors
		suite.Require().NoError(err, "Failed to get perpetual when owner mapping exists")
	}
	return perp, true
}

// CheckBalanceChange helper for asserting balance changes
func (suite *KeeperTestSuite) CheckBalanceChange(addr sdk.AccAddress, initial math.Int, expectedChange math.Int, msg string) {
	finalBalance := suite.GetAccountBalance(addr, QuoteDenom)
	expectedFinal := initial.Add(expectedChange)
	suite.Require().True(expectedFinal.Equal(finalBalance),
		"%s balance mismatch. Initial %s, Change %s, Expected %s, Got %s",
		msg, initial, expectedChange, expectedFinal, finalBalance)
}

func (suite *KeeperTestSuite) SetTwapPriceDirectly(marketId uint64, price math.LegacyDec) {
	// Implementation depends on how GetCurrentTwapPrice is controlled/mocked
	// Simplest for this example: set the global mock variable
	suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
		MarketId:          marketId,
		Block:             uint64(suite.ctx.BlockHeight()),
		AverageTradePrice: price,
		TotalVolume:       math.LegacyOneDec(),
		CumulativePrice:   math.LegacyOneDec(),
		Timestamp:         1,
	})
	suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
		MarketId:          marketId,
		Block:             uint64(suite.ctx.BlockHeight() + 1),
		AverageTradePrice: price,
		TotalVolume:       math.LegacyOneDec(),
		CumulativePrice:   math.LegacyOneDec().Add(price),
		Timestamp:         2,
	})
}

func (suite *KeeperTestSuite) SetTwapRecordDirectly(val types.TwapPrice) {
	// Implementation depends on how GetCurrentTwapPrice is controlled/mocked
	// Simplest for this example: set the global mock variable
	suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, val)
}
