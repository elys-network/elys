package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/utils"
	assetprofiletypes "github.com/elys-network/elys/v4/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/v4/x/commitment/types"
	"github.com/elys-network/elys/v4/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	stabletypes "github.com/elys-network/elys/v4/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (osmomath.BigDec, error) {
	masterchefParams := k.GetParams(ctx)
	estakingParams := k.estakingKeeper.GetParams(ctx)

	// If we don't have enough params
	if estakingParams.StakeIncentives == nil || masterchefParams.LpIncentives == nil {
		return osmomath.ZeroBigDec(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return osmomath.ZeroBigDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
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
				return osmomath.ZeroBigDec(), err
			}

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return osmomath.ZeroBigDec(), nil
			}

			if stkIncentive == nil || stkIncentive.EdenAmountPerYear.IsNil() {
				return osmomath.ZeroBigDec(), nil
			}

			// Calculate
			stakersEdenAmount := osmomath.BigDecFromSDKInt(stkIncentive.EdenAmountPerYear).QuoInt64(totalBlocksPerYear)

			// Maximum eden APR - 30% by default
			stakersMaxEdenAmount := estakingParams.GetBigDecMaxEdenRewardAprStakers().
				Mul(osmomath.BigDecFromSDKInt(totalStakedSnapshot)).
				QuoInt64(totalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = osmomath.MinBigDec(stakersEdenAmount, stakersMaxEdenAmount)

			// For Eden reward Apr for elys staking
			apr := stakersEdenAmount.
				MulInt64(totalBlocksPerYear).
				Quo(osmomath.BigDecFromSDKInt(totalStakedSnapshot))

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			borrowPool, found := k.stableKeeper.GetPoolByDenom(ctx, baseCurrency)
			if !found {
				return osmomath.ZeroBigDec(), fmt.Errorf("pool not found for denom %s", baseCurrency)
			}
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{
				PoolId: stabletypes.UsdcPoolId,
			})
			if err != nil {
				return osmomath.ZeroBigDec(), err
			}
			apr := borrowPool.GetBigDecInterestRate().MulDec(res.BorrowRatio)
			return apr, nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			// Get x days average rewards
			usdcAmount := k.GetAvgStakerFeesCollected(ctx, int(query.Days))
			if usdcAmount.IsZero() {
				return osmomath.ZeroBigDec(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be available in amm pool and has the same value as Elys
			edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return osmomath.ZeroBigDec(), nil
			}

			// Update total committed states
			totalStakedSnapshot, err := k.estakingKeeper.TotalBondedTokens(ctx)
			if err != nil {
				return osmomath.ZeroBigDec(), err
			}

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return osmomath.ZeroBigDec(), nil
			}

			// Multiply by 365 to get yearly rewards
			entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
			if !found {
				return osmomath.ZeroBigDec(), assetprofiletypes.ErrAssetProfileNotFound
			}
			yearlyDexRewardAmount := usdcAmount.MulInt64(365).QuoInt64(utils.Pow10Int64(entry.Decimals))

			apr := yearlyDexRewardAmount.
				Quo(edenDenomPrice).
				Quo(osmomath.BigDecFromSDKInt(totalStakedSnapshot))

			return apr, nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := estakingParams.GetBigDecEdenBoostApr()
		return apr, nil
	}

	return osmomath.ZeroBigDec(), nil
}

// Get total dex rewards amount from the specified pool
func (k Keeper) GetDailyRewardsAmountForPool(ctx sdk.Context, poolId uint64) (osmomath.BigDec, sdk.Coins) {
	dailyDexRewardsTotal := osmomath.ZeroBigDec()
	dailyGasRewardsTotal := osmomath.ZeroBigDec()
	dailyEdenRewardsTotal := osmomath.ZeroBigDec()
	firstAccum := k.FirstPoolRewardsAccum(ctx, poolId)
	lastAccum := k.LastPoolRewardsAccum(ctx, poolId)
	if lastAccum.Timestamp != 0 {
		if firstAccum.Timestamp == lastAccum.Timestamp {
			dailyDexRewardsTotal = lastAccum.GetBigDecDexReward()
			dailyGasRewardsTotal = lastAccum.GetBigDecGasReward()
			dailyEdenRewardsTotal = lastAccum.GetBigDecEdenReward()
		} else {
			dailyDexRewardsTotal = lastAccum.GetBigDecDexReward().Sub(firstAccum.GetBigDecDexReward())
			dailyGasRewardsTotal = lastAccum.GetBigDecGasReward().Sub(firstAccum.GetBigDecGasReward())
			dailyEdenRewardsTotal = lastAccum.GetBigDecEdenReward().Sub(firstAccum.GetBigDecEdenReward())
		}
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return osmomath.ZeroBigDec(), sdk.Coins{}
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, dailyEdenRewardsTotal.Dec().RoundInt()))
	rewardCoins = rewardCoins.Add(sdk.NewCoin(baseCurrency, dailyDexRewardsTotal.Add(dailyGasRewardsTotal).Dec().RoundInt()))

	usdcDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, baseCurrency)
	edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

	totalRewardsUsd := usdcDenomPrice.Mul(dailyDexRewardsTotal.Add(dailyGasRewardsTotal)).
		Add(edenDenomPrice.Mul(dailyEdenRewardsTotal))
	return totalRewardsUsd, rewardCoins
}
