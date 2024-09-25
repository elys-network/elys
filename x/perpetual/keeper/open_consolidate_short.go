package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidateShort(ctx sdk.Context, poolId uint64, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	existingMtp.Collateral = existingMtp.Collateral.Add(newMtp.Collateral)
	existingMtp.Custody = existingMtp.Custody.Add(newMtp.Custody)
	existingMtp.Liabilities = existingMtp.Liabilities.Add(newMtp.Liabilities)

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, existingMtp.AmmPoolId, existingMtp.CustodyAsset)
	if err != nil {
		return nil, err
	}

	custodyAmtTokenIn := sdk.NewCoin(existingMtp.CustodyAsset, existingMtp.Custody)
	repayAmount, err := k.EstimateSwap(ctx, custodyAmtTokenIn, existingMtp.TradingAsset, ammPool)
	if err != nil {
		return nil, err
	}

	unpaidCollateralIn := sdk.NewCoin(existingMtp.CollateralAsset, existingMtp.BorrowInterestUnpaidCollateral)
	C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
	if err != nil {
		return nil, err
	}

	updated_leverage := sdk.ZeroDec()
	denominator := repayAmount.Sub(C)
	if denominator.IsPositive() {
		updated_leverage = repayAmount.ToLegacyDec().Quo(denominator.ToLegacyDec())
	}
	existingMtp.Leverage = updated_leverage

	// Set existing MTP
	if err := k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	// Destroy new MTP
	if err := k.DestroyMTP(ctx, newMtp.GetAccountAddress(), newMtp.Id); err != nil {
		return nil, err
	}

	return existingMtp, nil
}
