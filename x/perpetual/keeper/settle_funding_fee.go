package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateFundingFee(ctx sdk.Context, mtp *types.MTP, pool *types.Pool) (bool, math.Int, error) {

	fullFundingFeePayment, fundingFeeAmt, err := k.FundingFeeCollection(ctx, mtp, pool)
	if err != nil {
		return fullFundingFeePayment, fundingFeeAmt, err
	}

	err = k.FundingFeeDistribution(ctx, mtp, pool)
	if err != nil {
		return fullFundingFeePayment, fundingFeeAmt, err
	}

	mtp.LastFundingCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastFundingCalcTime = uint64(ctx.BlockTime().Unix())

	return fullFundingFeePayment, fundingFeeAmt, nil
}

// SettleFunding handles funding fee collection and distribution
func (k Keeper) SettleFunding(ctx sdk.Context, mtp *types.MTP, pool *types.Pool) (bool, math.Int, error) {

	fundingFeeFullyPaid, fundingFeeAmt, err := k.UpdateFundingFee(ctx, mtp, pool)
	if err != nil {
		return fundingFeeFullyPaid, fundingFeeAmt, err
	}

	// apply changes to mtp object
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return fundingFeeFullyPaid, fundingFeeAmt, err
	}

	// apply changes to pool object
	k.SetPool(ctx, *pool)

	return fundingFeeFullyPaid, fundingFeeAmt, nil
}
