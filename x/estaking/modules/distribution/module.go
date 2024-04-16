package distribution

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/exported"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"

	abci "github.com/cometbft/cometbft/abci/types"
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

	keeper         keeper.Keeper
	accountKeeper  distrtypes.AccountKeeper
	bankKeeper     distrtypes.BankKeeper
	estakingKeeper *estakingkeeper.Keeper

	feeCollectorName string
}

// NewAppModule creates a new AppModule object using the native x/distribution module
// AppModule constructor.
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, ak distrtypes.AccountKeeper,
	bk distrtypes.BankKeeper, sk *estakingkeeper.Keeper, feeCollectorName string, subspace exported.Subspace,
) AppModule {
	distrAppMod := distr.NewAppModule(cdc, keeper, ak, bk, sk, subspace)
	return AppModule{
		AppModule:        distrAppMod,
		keeper:           keeper,
		accountKeeper:    ak,
		bankKeeper:       bk,
		estakingKeeper:   sk,
		feeCollectorName: feeCollectorName,
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

// AllocateTokens handles distribution of the collected fees
func (am AppModule) AllocateTokens(ctx sdk.Context) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the current representatives)
	feeCollector := am.accountKeeper.GetModuleAccount(ctx, am.feeCollectorName)
	feesCollectedInt := am.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	err := am.bankKeeper.SendCoinsFromModuleToModule(ctx, am.feeCollectorName, distrtypes.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	// calculate the fraction allocated to representatives by subtracting the community tax.
	// e.g. if community tax is 0.02, representatives fraction will be 0.98 (2% goes to the community pool and the rest to the representatives)
	remaining := feesCollected
	communityTax := am.keeper.GetCommunityTax(ctx)
	representativesFraction := sdk.OneDec().Sub(communityTax)

	totalBondedTokens := am.estakingKeeper.TotalBondedTokens(ctx)

	// allocate tokens proportionally to representatives voting power
	am.estakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		// we get this validator's percentage of the total power by dividing their tokens by the total bonded tokens
		powerFraction := sdk.NewDecFromInt(validator.GetTokens()).QuoTruncate(sdk.NewDecFromInt(totalBondedTokens))
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
