package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (elystypes.Dec34, error) {
	masterchefParams := k.GetParams(ctx)
	estakingParams := k.estakingKeeper.GetParams(ctx)

	// If we don't have enough params
	if estakingParams.StakeIncentives == nil || masterchefParams.LpIncentives == nil {
		return elystypes.ZeroDec34(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return elystypes.ZeroDec34(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
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
				return elystypes.ZeroDec34(), err
			}

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return elystypes.ZeroDec34(), nil
			}

			if stkIncentive == nil || stkIncentive.EdenAmountPerYear.IsNil() {
				return elystypes.ZeroDec34(), nil
			}

			// Calculate
			stakersEdenAmount := elystypes.NewDec34FromInt(stkIncentive.EdenAmountPerYear).QuoInt64(totalBlocksPerYear)

			// Maximum eden APR - 30% by default
			stakersMaxEdenAmount := elystypes.NewDec34FromLegacyDec(estakingParams.MaxEdenRewardAprStakers).
				MulInt(totalStakedSnapshot).
				QuoInt64(totalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = elystypes.MinDec34(stakersEdenAmount, stakersMaxEdenAmount)

			// For Eden reward Apr for elys staking
			apr := stakersEdenAmount.
				MulInt64(totalBlocksPerYear).
				QuoInt(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
			if err != nil {
				return elystypes.ZeroDec34(), err
			}
			apr := elystypes.NewDec34FromLegacyDec(params.InterestRate).MulLegacyDec(res.BorrowRatio)
			return apr, nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			// Get 7 days average rewards
			usdcAmount := k.GetAvgStakerFeesCollected(ctx)
			if usdcAmount.IsZero() {
				return elystypes.ZeroDec34(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenDenomPrice, _ := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return elystypes.ZeroDec34(), nil
			}

			// Update total committed states
			totalStakedSnapshot, err := k.estakingKeeper.TotalBondedTokens(ctx)
			if err != nil {
				return elystypes.ZeroDec34(), err
			}

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return elystypes.ZeroDec34(), nil
			}

			// Mutiply by 365 to get yearly rewards
			entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
			if !found {
				return elystypes.ZeroDec34(), assetprofiletypes.ErrAssetProfileNotFound
			}
			yearlyDexRewardAmount := elystypes.NewDec34FromLegacyDec(usdcAmount).MulInt64(365).QuoInt(ammtypes.BaseTokenAmount(entry.Decimals))

			apr := yearlyDexRewardAmount.
				Quo(edenDenomPrice).
				QuoInt(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := elystypes.NewDec34FromLegacyDec(estakingParams.EdenBoostApr)
		return apr, nil
	}

	return elystypes.ZeroDec34(), nil
}

// Get total dex rewards amount from the specified pool
func (k Keeper) GetDailyRewardsAmountForPool(ctx sdk.Context, poolId uint64) (elystypes.Dec34, sdk.Coins) {
	dailyDexRewardsTotal := elystypes.ZeroDec34()
	dailyGasRewardsTotal := elystypes.ZeroDec34()
	dailyEdenRewardsTotal := elystypes.ZeroDec34()
	firstAccum := k.FirstPoolRewardsAccum(ctx, poolId)
	lastAccum := k.LastPoolRewardsAccum(ctx, poolId)
	if lastAccum.Timestamp != 0 {
		if firstAccum.Timestamp == lastAccum.Timestamp {
			dailyDexRewardsTotal = elystypes.NewDec34FromLegacyDec(lastAccum.DexReward)
			dailyGasRewardsTotal = elystypes.NewDec34FromLegacyDec(lastAccum.GasReward)
			dailyEdenRewardsTotal = elystypes.NewDec34FromLegacyDec(lastAccum.EdenReward)
		} else {
			dailyDexRewardsTotal = elystypes.NewDec34FromLegacyDec(lastAccum.DexReward).SubLegacyDec(firstAccum.DexReward)
			dailyGasRewardsTotal = elystypes.NewDec34FromLegacyDec(lastAccum.GasReward).SubLegacyDec(firstAccum.GasReward)
			dailyEdenRewardsTotal = elystypes.NewDec34FromLegacyDec(lastAccum.EdenReward).SubLegacyDec(firstAccum.EdenReward)
		}
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return elystypes.ZeroDec34(), sdk.Coins{}
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, dailyEdenRewardsTotal.ToInt()))
	rewardCoins = rewardCoins.Add(sdk.NewCoin(baseCurrency, dailyDexRewardsTotal.Add(dailyGasRewardsTotal).ToInt()))

	usdcDenomPrice, _ := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	edenDenomPrice, _ := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

	totalRewardsUsd := (usdcDenomPrice.Mul(dailyDexRewardsTotal.Add(dailyGasRewardsTotal))).
		Add(edenDenomPrice.Mul(dailyEdenRewardsTotal))
	return totalRewardsUsd, rewardCoins
}
