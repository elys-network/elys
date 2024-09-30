package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/incentive/types"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeService        store.KVStoreService
		parameterKeeper     types.ParameterKeeper
		commitmentKeeper    types.CommitmentKeeper
		stk                 types.StakingKeeper
		authKeeper          types.AccountKeeper
		bankKeeper          types.BankKeeper
		amm                 types.AmmKeeper
		oracleKeeper        types.OracleKeeper
		assetProfileKeeper  types.AssetProfileKeeper
		accountedPoolKeeper types.AccountedPoolKeeper
		stableKeeper        types.StableStakeKeeper
		tokenomicsKeeper    types.TokenomicsKeeper
		masterchef          *masterchefkeeper.Keeper
		estaking            *estakingkeeper.Keeper

		authority string // gov module addresss
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	parameterKeeper types.ParameterKeeper,
	ck types.CommitmentKeeper,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	amm types.AmmKeeper,
	ok types.OracleKeeper,
	ap types.AssetProfileKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
	stableKeeper types.StableStakeKeeper,
	tokenomicsKeeper types.TokenomicsKeeper,
	masterchef *masterchefkeeper.Keeper,
	estaking *estakingkeeper.Keeper,
	feeCollectorName string,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		parameterKeeper:     parameterKeeper,
		commitmentKeeper:    ck,
		stk:                 sk,
		authKeeper:          ak,
		bankKeeper:          bk,
		amm:                 amm,
		oracleKeeper:        ok,
		assetProfileKeeper:  ap,
		accountedPoolKeeper: accountedPoolKeeper,
		stableKeeper:        stableKeeper,
		tokenomicsKeeper:    tokenomicsKeeper,
		masterchef:          masterchef,
		estaking:            estaking,
		authority:           authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Caculate total TVL
func (k Keeper) CalculateTVL(ctx sdk.Context) math.LegacyDec {
	TVL := math.LegacyZeroDec()

	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}
		TVL = TVL.Add(tvl)
		return false
	})

	return TVL
}

// Get total dex rewards amount from the specified pool
func (k Keeper) GetDailyRewardsAmountForPool(ctx sdk.Context, poolId uint64) (math.LegacyDec, sdk.Coins) {
	dailyDexRewardsTotal := math.LegacyZeroDec()
	dailyGasRewardsTotal := math.LegacyZeroDec()
	dailyEdenRewardsTotal := math.LegacyZeroDec()
	firstAccum := k.masterchef.FirstPoolRewardsAccum(ctx, poolId)
	lastAccum := k.masterchef.LastPoolRewardsAccum(ctx, poolId)
	if lastAccum.Timestamp != 0 {
		if firstAccum.Timestamp == lastAccum.Timestamp {
			dailyDexRewardsTotal = lastAccum.DexReward
			dailyGasRewardsTotal = lastAccum.GasReward
			dailyEdenRewardsTotal = lastAccum.EdenReward
		} else {
			dailyDexRewardsTotal = lastAccum.DexReward.Sub(firstAccum.DexReward)
			dailyGasRewardsTotal = lastAccum.GasReward.Sub(firstAccum.GasReward)
			dailyEdenRewardsTotal = lastAccum.EdenReward.Sub(firstAccum.EdenReward)
		}
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return math.LegacyZeroDec(), sdk.Coins{}
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, dailyEdenRewardsTotal.RoundInt()))
	rewardCoins = rewardCoins.Add(sdk.NewCoin(baseCurrency, dailyDexRewardsTotal.Add(dailyGasRewardsTotal).RoundInt()))

	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

	totalRewardsUsd := usdcDenomPrice.Mul(dailyDexRewardsTotal.Add(dailyGasRewardsTotal)).
		Add(edenDenomPrice.Mul(dailyEdenRewardsTotal))
	return totalRewardsUsd, rewardCoins
}
