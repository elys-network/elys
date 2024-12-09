package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (math.LegacyDec, error) {
	masterchefParams := k.GetParams(ctx)
	estakingParams := k.estakingKeeper.GetParams(ctx)

	// If we don't have enough params
	if estakingParams.StakeIncentives == nil || masterchefParams.LpIncentives == nil {
		return math.LegacyZeroDec(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return math.LegacyZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	stkIncentive := estakingParams.StakeIncentives

	totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)

	if query.Denom == ptypes.Eden {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			return k.CalculateStableStakeApr(ctx, &types.QueryStableStakeAprRequest{
				Denom: ptypes.Eden,
			})
		} else {
			// Elys staking, Eden committed, EdenB committed.
			totalStakedSnapshot, err := k.estakingKeeper.TotalBondedTokens(ctx)
			if err != nil {
				return math.LegacyZeroDec(), err
			}

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return math.LegacyZeroDec(), nil
			}

			if stkIncentive == nil || stkIncentive.EdenAmountPerYear.IsNil() {
				return math.LegacyZeroDec(), nil
			}

			// Calculate
			stakersEdenAmount := stkIncentive.EdenAmountPerYear.ToLegacyDec().Quo(math.NewInt(totalBlocksPerYear).ToLegacyDec())

			// Maximum eden APR - 30% by default
			stakersMaxEdenAmount := estakingParams.MaxEdenRewardAprStakers.
				MulInt(totalStakedSnapshot).
				QuoInt64(totalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = math.LegacyMinDec(stakersEdenAmount, stakersMaxEdenAmount)

			// For Eden reward Apr for elys staking
			apr := stakersEdenAmount.
				Mul(math.LegacyNewDec(totalBlocksPerYear)).
				Quo(totalStakedSnapshot.ToLegacyDec())

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
			if err != nil {
				return math.LegacyZeroDec(), err
			}
			apr := params.InterestRate.Mul(res.BorrowRatio)
			return apr, nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			// Get 7 days average rewards
			usdcAmount := k.GetAvgStakerFeesCollected(ctx)
			if usdcAmount.IsZero() {
				return math.LegacyZeroDec(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return math.LegacyZeroDec(), nil
			}

			// Update total committed states
			totalStakedSnapshot, err := k.estakingKeeper.TotalBondedTokens(ctx)
			if err != nil {
				return math.LegacyZeroDec(), err
			}

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return math.LegacyZeroDec(), nil
			}

			// Mutiply by 365 to get yearly rewards
			denom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
			if !found {
				return math.LegacyZeroDec(), assetprofiletypes.ErrAssetProfileNotFound
			}
			entry, found := k.assetProfileKeeper.GetEntry(ctx, denom)
			if !found {
				return math.LegacyZeroDec(), assetprofiletypes.ErrAssetProfileNotFound
			}
			yearlyDexRewardAmount := usdcAmount.Mul(math.LegacyNewDec(365)).Quo(math.LegacyNewDec(int64(entry.Decimals)))

			apr := yearlyDexRewardAmount.
				Quo(edenDenomPrice).
				QuoInt(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := estakingParams.EdenBoostApr
		return apr, nil
	}

	return math.LegacyZeroDec(), nil
}

// Get total dex rewards amount from the specified pool
func (k Keeper) GetDailyRewardsAmountForPool(ctx sdk.Context, poolId uint64) (math.LegacyDec, sdk.Coins) {
	dailyDexRewardsTotal := math.LegacyZeroDec()
	dailyGasRewardsTotal := math.LegacyZeroDec()
	dailyEdenRewardsTotal := math.LegacyZeroDec()
	firstAccum := k.FirstPoolRewardsAccum(ctx, poolId)
	lastAccum := k.LastPoolRewardsAccum(ctx, poolId)
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
