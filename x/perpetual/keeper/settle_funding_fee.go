package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateFundingFee(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) error {

	err := k.FundingFeeCollection(ctx, mtp, pool, ammPool)
	if err != nil {
		return err
	}

	err = k.FundingFeeDistribution(ctx, mtp, pool, ammPool)
	if err != nil {
		return err
	}

	mtp.LastFundingCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastFundingCalcTime = uint64(ctx.BlockTime().Unix())

	return nil
}

// SettleFunding handles funding fee collection and distribution
func (k Keeper) SettleFunding(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) error {

	err := k.UpdateFundingFee(ctx, mtp, pool, ammPool)
	if err != nil {
		return err
	}

	// apply changes to mtp object
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	// apply changes to pool object
	k.SetPool(ctx, *pool)

	return nil
}
