package distribution

import (
	"context"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/exported"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	assetprofilekeeper "github.com/elys-network/elys/v7/x/assetprofile/keeper"
	commitmentkeeper "github.com/elys-network/elys/v7/x/commitment/keeper"
	estakingkeeper "github.com/elys-network/elys/v7/x/estaking/keeper"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
)

var (
	_ module.AppModuleBasic      = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasServices         = AppModule{}
	_ module.HasInvariants       = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
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
	commitmentKeeper   *commitmentkeeper.Keeper
	estakingKeeper     *estakingkeeper.Keeper
	assetprofileKeeper *assetprofilekeeper.Keeper

	feeCollectorName string
}

// NewAppModule creates a new AppModule object using the native x/distribution module
// AppModule constructor.
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, ak distrtypes.AccountKeeper,
	ck *commitmentkeeper.Keeper,
	sk *estakingkeeper.Keeper,
	assetprofileKeeper *assetprofilekeeper.Keeper,
	feeCollectorName string, subspace exported.Subspace,
) AppModule {
	distrAppMod := distr.NewAppModule(cdc, keeper, ak, ck, sk, subspace)
	return AppModule{
		AppModule:          distrAppMod,
		keeper:             keeper,
		accountKeeper:      ak,
		commitmentKeeper:   ck,
		estakingKeeper:     sk,
		assetprofileKeeper: assetprofileKeeper,
		feeCollectorName:   feeCollectorName,
	}
}

// BeginBlock mirror functionality of cosmos-sdk/distribution BeginBlocker
// however it allocates no proposer reward
func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	defer telemetry.ModuleMeasureSince(distrtypes.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	var previousTotalPower int64
	for _, voteInfo := range ctx.VoteInfos() {
		previousTotalPower += voteInfo.Validator.Power
	}

	if ctx.BlockHeight() > 1 {
		am.AllocateEdenUsdcTokens(ctx)
		am.AllocateEdenBTokens(ctx)
	}

	consAddr := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	return am.keeper.SetPreviousProposerConsAddr(ctx, consAddr)
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

// AllocateEdenUsdcTokens handles distribution of the collected fees
// USDC and Eden is distributed for staking Elys and locking Eden and locking EdenB
func (am AppModule) AllocateEdenUsdcTokens(ctx sdk.Context) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the current representatives)
	feeCollector := am.accountKeeper.GetModuleAccount(ctx, ccvconsumertypes.ConsumerRedistributeName)
	feesCollectedInt := am.commitmentKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	usdcDenom, _ := am.assetprofileKeeper.GetUsdcDenom(ctx)
	filteredCoins := FilterDenoms(feesCollectedInt, usdcDenom, ptypes.Eden)
	feesCollected := sdk.NewDecCoinsFromCoins(filteredCoins...)

	// transfer collected fees to the distribution module account
	if filteredCoins.IsAllPositive() {
		err := am.commitmentKeeper.SendCoinsFromModuleToModule(ctx, ccvconsumertypes.ConsumerRedistributeName, distrtypes.ModuleName, filteredCoins)
		if err != nil {
			panic(err)
		}
	}

	// calculate the fraction allocated to representatives by subtracting the community tax.
	// e.g. if community tax is 0.02, representatives fraction will be 0.98 (2% goes to the community pool and the rest to the representatives)
	remaining := feesCollected
	communityTax, err := am.keeper.GetCommunityTax(ctx)
	if err != nil {
		panic(err)
	}
	representativesFraction := math.LegacyOneDec().Sub(communityTax)

	// Note: to prevent negative coin amount issue when invariant's broken,
	// calculation of total bonded tokens manually through iteration
	sumOfValTokens := math.ZeroInt()
	err = am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		sumOfValTokens = sumOfValTokens.Add(validator.GetTokens())
		return false
	})
	if err != nil {
		panic(err)
	}

	// We consider ELYS + EDEN + EDENB bonded tokens here
	totalBondedTokens, err := am.estakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		panic(err)
	}
	if !totalBondedTokens.Equal(sumOfValTokens) {
		ctx.Logger().Error("invariant broken", "sumOfValTokens", sumOfValTokens.String(), "totalBondedTokens", totalBondedTokens.String())
	}

	sumOfValTokensDec := math.LegacyNewDecFromInt(sumOfValTokens)
	// allocate tokens proportionally to representatives voting power
	err = am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		// we get this validator's percentage of the total power by dividing their tokens by the total bonded tokens
		powerFraction := math.LegacyNewDecFromInt(validator.GetTokens()).QuoTruncate(sumOfValTokensDec)
		// we truncate here again, which means that the reward will be slightly lower than it should be
		reward := feesCollected.MulDecTruncate(representativesFraction).MulDecTruncate(powerFraction)
		err = am.keeper.AllocateTokensToValidator(ctx, validator, reward)
		if err != nil {
			panic(err)
		}
		remaining = remaining.Sub(reward)
		return false
	})
	if err != nil {
		panic(err)
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	feePool, err := am.keeper.FeePool.Get(ctx)
	if err != nil {
		panic(err)
	}
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	err = am.keeper.FeePool.Set(ctx, feePool)
	if err != nil {
		panic(err)
	}
}

// AllocateEdenBTokens handles distribution of the collected fees
// EdenB is distributed for staking Elys and locking Eden, not for locking EdenB
func (am AppModule) AllocateEdenBTokens(ctx sdk.Context) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the current representatives)
	feeCollector := am.accountKeeper.GetModuleAccount(ctx, ccvconsumertypes.ConsumerRedistributeName)
	feesCollectedInt := am.commitmentKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	filteredCoins := FilterDenoms(feesCollectedInt, ptypes.EdenB)
	feesCollected := sdk.NewDecCoinsFromCoins(filteredCoins...)

	// transfer collected fees to the distribution module account
	if filteredCoins.IsAllPositive() {
		err := am.commitmentKeeper.SendCoinsFromModuleToModule(ctx, ccvconsumertypes.ConsumerRedistributeName, distrtypes.ModuleName, filteredCoins)
		if err != nil {
			panic(err)
		}
	}

	// calculate the fraction allocated to representatives by subtracting the community tax.
	// e.g. if community tax is 0.02, representatives fraction will be 0.98 (2% goes to the community pool and the rest to the representatives)
	remaining := feesCollected
	communityTax, err := am.keeper.GetCommunityTax(ctx)
	if err != nil {
		panic(err)
	}
	representativesFraction := math.LegacyOneDec().Sub(communityTax)

	edenBValidator := am.estakingKeeper.GetEdenBValidator(ctx).GetOperator()

	// Note: to prevent negative coin amount issue when invariant's broken,
	// calculation of total bonded tokens manually through iteration
	sumOfValTokens := math.ZeroInt()
	err = am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		if validator.GetOperator() != edenBValidator {
			sumOfValTokens = sumOfValTokens.Add(validator.GetTokens())
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// We consider ELYS and EDEN bonded tokens, EdenB bonded tokens are excluded here
	totalBondedElysEdenTokens, err := am.estakingKeeper.TotalBondedElysEdenTokens(ctx)
	if err != nil {
		panic(err)
	}
	if !totalBondedElysEdenTokens.Equal(sumOfValTokens) {
		ctx.Logger().Error("invariant broken", "sumOfValTokens", sumOfValTokens.String(), "totalBondedElysEdenTokens", totalBondedElysEdenTokens.String())
	}

	sumOfValTokensDec := math.LegacyNewDecFromInt(sumOfValTokens)
	// allocate tokens proportionally to representatives voting power
	err = am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		if validator.GetOperator() == edenBValidator {
			return false
		}
		// we get this validator's percentage of the total power by dividing their tokens by the total bonded tokens
		powerFraction := math.LegacyNewDecFromInt(validator.GetTokens()).QuoTruncate(sumOfValTokensDec)
		// we truncate here again, which means that the reward will be slightly lower than it should be
		reward := feesCollected.MulDecTruncate(representativesFraction).MulDecTruncate(powerFraction)
		err = am.keeper.AllocateTokensToValidator(ctx, validator, reward)
		if err != nil {
			panic(err)
		}
		remaining = remaining.Sub(reward)
		return false
	})
	if err != nil {
		panic(err)
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	feePool, err := am.keeper.FeePool.Get(ctx)
	if err != nil {
		panic(err)
	}
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	err = am.keeper.FeePool.Set(ctx, feePool)
	if err != nil {
		panic(err)
	}
}
