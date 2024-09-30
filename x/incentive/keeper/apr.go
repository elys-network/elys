package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (sdkmath.LegacyDec, error) {
	masterchefParams := k.masterchef.GetParams(ctx)
	estakingParams := k.estaking.GetParams(ctx)

	// If we don't have enough params
	if estakingParams.StakeIncentives == nil || masterchefParams.LpIncentives == nil {
		return sdkmath.LegacyZeroDec(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return sdkmath.LegacyZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	stkIncentive := estakingParams.StakeIncentives

	totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear

	if query.Denom == ptypes.Eden {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			return k.masterchef.CalculateStableStakeApr(ctx, &mastercheftypes.QueryStableStakeAprRequest{
				Denom: ptypes.Eden,
			})
		} else {
			// Elys staking, Eden committed, EdenB committed.
			totalStakedSnapshot := k.estaking.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdkmath.LegacyZeroDec(), nil
			}

			if stkIncentive == nil || stkIncentive.EdenAmountPerYear.IsNil() {
				return sdkmath.LegacyZeroDec(), nil
			}

			// Calculate
			stakersEdenAmount := stkIncentive.EdenAmountPerYear.Quo(sdkmath.NewInt(totalBlocksPerYear)).ToLegacyDec()

			// Maximum eden APR - 30% by default
			stakersMaxEdenAmount := estakingParams.MaxEdenRewardAprStakers.
				MulInt(totalStakedSnapshot).
				QuoInt64(totalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = sdkmath.LegacyMinDec(stakersEdenAmount, stakersMaxEdenAmount)

			// For Eden reward Apr for elys staking
			apr := stakersEdenAmount.
				Mul(sdkmath.LegacyNewDec(totalBlocksPerYear)).
				Quo(sdkmath.LegacyNewDecFromInt(totalStakedSnapshot))

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
			if err != nil {
				return sdkmath.LegacyZeroDec(), err
			}
			apr := params.InterestRate.Mul(res.BorrowRatio)
			return apr, nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			params := k.estaking.GetParams(ctx)
			amount := params.DexRewardsStakers.Amount
			if amount.IsZero() {
				return sdkmath.LegacyZeroDec(), nil
			}

			// If no rewards were given.
			if params.DexRewardsStakers.NumBlocks == 0 {
				return sdkmath.LegacyZeroDec(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return sdkmath.LegacyZeroDec(), nil
			}

			// Update total committed states
			totalStakedSnapshot := k.estaking.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdkmath.LegacyZeroDec(), nil
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

	return sdkmath.LegacyZeroDec(), nil
}
