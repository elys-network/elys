package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
		storeKey            storetypes.StoreKey
		memKey              storetypes.StoreKey
		cmk                 types.CommitmentKeeper
		stk                 types.StakingKeeper
		tci                 *types.TotalCommitmentInfo
		authKeeper          types.AccountKeeper
		bankKeeper          types.BankKeeper
		amm                 types.AmmKeeper
		oracleKeeper        types.OracleKeeper
		assetProfileKeeper  types.AssetProfileKeeper
		accountedPoolKeeper types.AccountedPoolKeeper
		epochsKeeper        types.EpochsKeeper
		stableKeeper        types.StableStakeKeeper
		tokenomicsKeeper    types.TokenomicsKeeper
		masterchef          *masterchefkeeper.Keeper
		estaking            *estakingkeeper.Keeper

		authority string // gov module addresss
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ck types.CommitmentKeeper,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	amm types.AmmKeeper,
	ok types.OracleKeeper,
	ap types.AssetProfileKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
	epochsKeeper types.EpochsKeeper,
	stableKeeper types.StableStakeKeeper,
	tokenomicsKeeper types.TokenomicsKeeper,
	masterchef *masterchefkeeper.Keeper,
	estaking *estakingkeeper.Keeper,
	feeCollectorName string,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		cmk:                 ck,
		stk:                 sk,
		tci:                 &types.TotalCommitmentInfo{},
		authKeeper:          ak,
		bankKeeper:          bk,
		amm:                 amm,
		oracleKeeper:        ok,
		assetProfileKeeper:  ap,
		accountedPoolKeeper: accountedPoolKeeper,
		epochsKeeper:        epochsKeeper,
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
func (k Keeper) CalculateTVL(ctx sdk.Context) sdk.Dec {
	TVL := sdk.ZeroDec()

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
func (k Keeper) GetDailyRewardsAmountForPool(ctx sdk.Context, poolId uint64) (sdk.Dec, sdk.Coins) {
	poolInfo, found := k.masterchef.GetPool(ctx, poolId)
	if !found {
		return sdk.ZeroDec(), sdk.Coins{}
	}

	// Fetch incentive params
	params := k.masterchef.GetParams(ctx)
	if params.LpIncentives == nil {
		return sdk.ZeroDec(), sdk.Coins{}
	}

	// Dex reward Apr per pool =  total accumulated usdc rewards for 7 day * 52/ tvl of pool
	dailyDexRewardsTotal := poolInfo.DexRewardAmountGiven.
		QuoInt(poolInfo.NumBlocks)

	// Eden reward Apr per pool = (total LM Eden reward allocated per day*((tvl of pool * multiplier)/total proxy TVL) ) * 365 / TVL of pool
	dailyEdenRewardsTotal := poolInfo.EdenRewardAmountGiven.
		Quo(poolInfo.NumBlocks)

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroDec(), sdk.Coins{}
	}
	baseCurrency := entry.Denom

	rewardCoins := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, dailyEdenRewardsTotal))
	rewardCoins = rewardCoins.Add(sdk.NewCoin(baseCurrency, math.Int(dailyDexRewardsTotal)))

	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

	totalRewardsUsd := usdcDenomPrice.Mul(dailyDexRewardsTotal).Add(edenDenomPrice.MulInt(dailyEdenRewardsTotal))
	return totalRewardsUsd, rewardCoins
}
