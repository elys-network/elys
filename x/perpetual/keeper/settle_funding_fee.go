package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

func (k Keeper) UpdateFundingFee(ctx sdk.Context, mtp *types.MTP, pool *types.Pool) (bool, math.Int, math.Int, error) {

	fullFundingFeePayment, fundingFeeAmt, err := k.FundingFeeCollection(ctx, mtp, pool)
	if err != nil {
		return fullFundingFeePayment, fundingFeeAmt, math.ZeroInt(), err
	}

	amountDistributed, err := k.FundingFeeDistribution(ctx, mtp, pool)
	if err != nil {
		return fullFundingFeePayment, fundingFeeAmt, amountDistributed, err
	}

	mtp.LastFundingCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastFundingCalcTime = uint64(ctx.BlockTime().Unix())

	return fullFundingFeePayment, fundingFeeAmt, amountDistributed, nil
}

// SettleFunding handles funding fee collection and distribution
func (k Keeper) SettleFunding(ctx sdk.Context, mtp *types.MTP, pool *types.Pool) (bool, math.Int, math.Int, error) {

	fundingFeeFullyPaid, fundingFeeAmt, amountDistributed, err := k.UpdateFundingFee(ctx, mtp, pool)
	if err != nil {
		return fundingFeeFullyPaid, fundingFeeAmt, amountDistributed, err
	}

	// apply changes to mtp object
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return fundingFeeFullyPaid, fundingFeeAmt, amountDistributed, err
	}

	// apply changes to pool object
	k.SetPool(ctx, *pool)

	return fundingFeeFullyPaid, fundingFeeAmt, amountDistributed, nil
}
