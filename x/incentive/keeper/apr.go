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

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (math.Int, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	masterchefParams := k.masterchef.GetParams(ctx)
	estakingParams := k.estaking.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if estakingParams.StakeIncentives == nil || masterchefParams.LpIncentives == nil {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroInt(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	baseCurrency := entry.Denom
	lpIncentive := masterchefParams.LpIncentives
	stkIncentive := estakingParams.StakeIncentives

	if lpIncentive.TotalBlocksPerYear.IsZero() || stkIncentive.TotalBlocksPerYear.IsZero() {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

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
				return sdk.ZeroInt(), nil
			}

			// Calculate
			stakersEdenAmount := stkIncentive.EdenAmountPerYear.
				Quo(stkIncentive.TotalBlocksPerYear)

			// Maximum eden APR - 30% by default
			// Allocated for staking per day = (0.3/365)* ( total elys staked + total Eden committed + total Eden boost committed)
			stakersMaxEdenAmount := estakingParams.MaxEdenRewardAprStakers.
				MulInt(totalStakedSnapshot).
				QuoInt(stkIncentive.TotalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			stakersEdenAmount = sdk.MinInt(stakersEdenAmount, stakersMaxEdenAmount.TruncateInt())

			// For Eden reward Apr for elys staking = {(amount of Eden allocated for staking per day)*365/( total elys staked + total Eden committed + total Eden boost committed)}*100
			apr := stakersEdenAmount.
				Mul(lpIncentive.TotalBlocksPerYear).
				Mul(sdk.NewInt(100)).
				Quo(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
			if err != nil {
				return sdk.ZeroInt(), err
			}
			apr := params.InterestRate.Mul(res.BorrowRatio).MulInt(sdk.NewInt(100))
			return apr.TruncateInt(), nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			params := k.estaking.GetParams(ctx)
			amount := params.DexRewardsStakers.Amount
			if amount.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// If no rewards were given.
			if params.DexRewardsStakers.NumBlocks.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
			if edenDenomPrice.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Update total committed states
			totalStakedSnapshot := k.estaking.TotalBondedTokens(ctx)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// DexReward amount per day = amount distributed / duration(in seconds) * total seconds per day.
			yearlyDexRewardAmount := amount.MulInt(lpIncentive.TotalBlocksPerYear).QuoInt(params.DexRewardsStakers.NumBlocks)

			// Usdc apr for elys staking = (24 hour dex rewards in USDC generated for stakers) * 365*100/ {price ( elys/usdc)*( sum of (elys staked, Eden committed, Eden boost committed))}
			// we multiply 10 as we have use 10elys as input in the price estimation
			apr := yearlyDexRewardAmount.
				MulInt(sdk.NewInt(100)).
				Quo(edenDenomPrice).
				QuoInt(totalStakedSnapshot)

			return apr.TruncateInt(), nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := estakingParams.EdenBoostApr.MulInt(sdk.NewInt(100)).TruncateInt()
		return apr, nil
	}

	return sdk.ZeroInt(), nil
}
