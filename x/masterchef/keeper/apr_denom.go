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

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (sdk.Dec, error) {
	masterchefParams := k.GetParams(ctx)
	estakingParams := k.estakingKeeper.GetParams(ctx)

	// If we don't have enough params
	if estakingParams.StakeIncentives == nil || masterchefParams.LpIncentives == nil {
		return sdk.ZeroDec(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return sdk.ZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	stkIncentive := estakingParams.StakeIncentives

	totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear

	if query.Denom == ptypes.Eden {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			return k.CalculateStableStakeApr(ctx, &types.QueryStableStakeAprRequest{
				Denom: ptypes.Eden,
			})
		} else {
			// Elys staking, Eden committed, EdenB committed.
			totalStakedSnapshot := k.estakingKeeper.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroDec(), nil
			}

			if stkIncentive == nil || stkIncentive.EdenAmountPerYear.IsNil() {
				return sdk.ZeroDec(), nil
			}

			// Calculate
			stakersEdenAmount := stkIncentive.EdenAmountPerYear.ToLegacyDec().Quo(sdk.NewInt(totalBlocksPerYear).ToLegacyDec())

			// Maximum eden APR - 30% by default
			stakersMaxEdenAmount := estakingParams.MaxEdenRewardAprStakers.
				MulInt(totalStakedSnapshot).
				QuoInt64(totalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = sdk.MinDec(stakersEdenAmount, stakersMaxEdenAmount)

			// For Eden reward Apr for elys staking
			apr := stakersEdenAmount.
				Mul(sdk.NewDec(totalBlocksPerYear)).
				Quo(sdk.NewDecFromInt(totalStakedSnapshot))

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
			if err != nil {
				return sdk.ZeroDec(), err
			}
			apr := params.InterestRate.Mul(res.BorrowRatio)
			return apr, nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			params := k.estakingKeeper.GetParams(ctx)
			amount := params.DexRewardsStakers.Amount
			if amount.IsZero() {
				return sdk.ZeroDec(), nil
			}

			// If no rewards were given.
			if params.DexRewardsStakers.NumBlocks == 0 {
				return sdk.ZeroDec(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return sdk.ZeroDec(), nil
			}

			// Update total committed states
			totalStakedSnapshot := k.estakingKeeper.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroDec(), nil
			}

			usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
			yearlyDexRewardAmount := amount.
				Mul(usdcDenomPrice).
				MulInt64(totalBlocksPerYear).
				QuoInt64(params.DexRewardsStakers.NumBlocks)

			apr := yearlyDexRewardAmount.
				Quo(edenDenomPrice).
				QuoInt(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := estakingParams.EdenBoostApr
		return apr, nil
	}

	return sdk.ZeroDec(), nil
}

// Get total dex rewards amount from the specified pool
func (k Keeper) GetDailyRewardsAmountForPool(ctx sdk.Context, poolId uint64) (sdk.Dec, sdk.Coins) {
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
		return sdk.ZeroDec(), sdk.Coins{}
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, dailyEdenRewardsTotal.RoundInt()))
	rewardCoins = rewardCoins.Add(sdk.NewCoin(baseCurrency, dailyDexRewardsTotal.Add(dailyGasRewardsTotal).RoundInt()))

	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

	totalRewardsUsd := usdcDenomPrice.Mul(dailyDexRewardsTotal.Add(dailyGasRewardsTotal)).
		Add(edenDenomPrice.Mul(dailyEdenRewardsTotal))
	return totalRewardsUsd, rewardCoins
}
