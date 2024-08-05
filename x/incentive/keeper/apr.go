package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (math.LegacyDec, error) {
	masterchefParams := k.masterchef.GetParams(ctx)
	estakingParams := k.estaking.GetParams(ctx)

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
			return k.masterchef.CalculateStableStakeApr(ctx, &mastercheftypes.QueryStableStakeAprRequest{
				Denom: ptypes.Eden,
			})
		} else {
			// Elys staking, Eden committed, EdenB committed.
			totalStakedSnapshot := k.estaking.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroDec(), nil
			}

			if stkIncentive == nil || stkIncentive.EdenAmountPerYear.IsNil() {
				return sdk.ZeroDec(), nil
			}

			// Calculate
			stakersEdenAmount := stkIncentive.EdenAmountPerYear.Quo(sdk.NewInt(totalBlocksPerYear))

			// Maximum eden APR - 30% by default
			stakersMaxEdenAmount := estakingParams.MaxEdenRewardAprStakers.
				MulInt(totalStakedSnapshot).
				QuoInt64(totalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = sdk.MinInt(stakersEdenAmount, stakersMaxEdenAmount.TruncateInt())

			// For Eden reward Apr for elys staking
			apr := stakersEdenAmount.
				Mul(sdk.NewInt(totalBlocksPerYear)).
				Quo(totalStakedSnapshot)

			return apr.ToLegacyDec(), nil
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
			params := k.estaking.GetParams(ctx)
			amount := params.DexRewardsStakers.Amount
			if amount.IsZero() {
				return sdk.ZeroDec(), nil
			}

			// If no rewards were given.
			if params.DexRewardsStakers.NumBlocks.IsZero() {
				return sdk.ZeroDec(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return sdk.ZeroDec(), nil
			}

			// Update total committed states
			totalStakedSnapshot := k.estaking.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroDec(), nil
			}

			usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
			yearlyDexRewardAmount := amount.
				Mul(usdcDenomPrice).
				MulInt64(totalBlocksPerYear).
				QuoInt(params.DexRewardsStakers.NumBlocks)

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
