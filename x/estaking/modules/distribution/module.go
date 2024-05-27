package distribution

import (
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/exported"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	assetprofilekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/distribution AppModuleBasic.
type AppModuleBasic struct {
	distr.AppModuleBasic
}

// AppModule embeds the Cosmos SDK's x/distribution AppModule
type AppModule struct {
	// embed the Cosmos SDK's x/distribution AppModule
	distr.AppModule

	keeper             keeper.Keeper
	accountKeeper      distrtypes.AccountKeeper
	bankKeeper         distrtypes.BankKeeper
	estakingKeeper     *estakingkeeper.Keeper
	assetprofileKeeper *assetprofilekeeper.Keeper

	feeCollectorName string
}

// NewAppModule creates a new AppModule object using the native x/distribution module
// AppModule constructor.
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, ak distrtypes.AccountKeeper,
	bk distrtypes.BankKeeper,
	sk *estakingkeeper.Keeper,
	assetprofileKeeper *assetprofilekeeper.Keeper,
	feeCollectorName string, subspace exported.Subspace,
) AppModule {
	distrAppMod := distr.NewAppModule(cdc, keeper, ak, bk, sk, subspace)
	return AppModule{
		AppModule:          distrAppMod,
		keeper:             keeper,
		accountKeeper:      ak,
		bankKeeper:         bk,
		estakingKeeper:     sk,
		assetprofileKeeper: assetprofileKeeper,
		feeCollectorName:   feeCollectorName,
	}
}

// BeginBlocker mirror functionality of cosmos-sdk/distribution BeginBlocker
// however it allocates no proposer reward
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	defer telemetry.ModuleMeasureSince(distrtypes.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if ctx.BlockHeight() > 1 {
		am.AllocateTokens(ctx)
	}
}

// RegisterInvariants registers the distribution module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

func FilterDenoms(coins sdk.Coins, denoms ...string) sdk.Coins {
	filtered := sdk.Coins{}
	for _, denom := range denoms {
		filtered = filtered.Add(sdk.NewCoin(denom, coins.AmountOf(denom)))
	}
	return filtered
}

// AllocateTokens handles distribution of the collected fees
func (am AppModule) AllocateTokens(ctx sdk.Context) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the current representatives)
	feeCollector := am.accountKeeper.GetModuleAccount(ctx, am.feeCollectorName)
	feesCollectedInt := am.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	usdcDenom, _ := am.assetprofileKeeper.GetUsdcDenom(ctx)
	filteredCoins := FilterDenoms(feesCollectedInt, usdcDenom, ptypes.Eden, ptypes.EdenB)
	feesCollected := sdk.NewDecCoinsFromCoins(filteredCoins...)

	// transfer collected fees to the distribution module account
	if filteredCoins.IsAllPositive() {
		err := am.bankKeeper.SendCoinsFromModuleToModule(ctx, am.feeCollectorName, distrtypes.ModuleName, filteredCoins)
		if err != nil {
			panic(err)
		}
	}

	// calculate the fraction allocated to representatives by subtracting the community tax.
	// e.g. if community tax is 0.02, representatives fraction will be 0.98 (2% goes to the community pool and the rest to the representatives)
	remaining := feesCollected
	communityTax := am.keeper.GetCommunityTax(ctx)
	representativesFraction := sdk.OneDec().Sub(communityTax)

	// Note: to prevent negative coin amount issue when invariant's broken,
	// calculation of total bonded tokens manually through iteration
	sumOfValTokens := math.ZeroInt()
	am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		sumOfValTokens = sumOfValTokens.Add(validator.GetTokens())
		return false
	})

	totalBondedTokens := am.estakingKeeper.TotalBondedTokens(ctx)
	if !totalBondedTokens.Equal(sumOfValTokens) {
		ctx.Logger().Error("invariant broken", "sumOfValTokens", sumOfValTokens.String(), "totalBondedTokens", totalBondedTokens.String())
	}

	sumOfValTokensDec := sdk.NewDecFromInt(sumOfValTokens)
	// allocate tokens proportionally to representatives voting power
	am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		// we get this validator's percentage of the total power by dividing their tokens by the total bonded tokens
		powerFraction := sdk.NewDecFromInt(validator.GetTokens()).QuoTruncate(sumOfValTokensDec)
		// we truncate here again, which means that the reward will be slightly lower than it should be
		reward := feesCollected.MulDecTruncate(representativesFraction).MulDecTruncate(powerFraction)
		am.keeper.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
		return false
	})

	// temporary workaround to keep CanWithdrawInvariant happy
	feePool := am.keeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	am.keeper.SetFeePool(ctx, feePool)
}
