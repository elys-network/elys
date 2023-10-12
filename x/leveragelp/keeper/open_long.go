package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Calculate the eta value.
	eta := leverage.Sub(sdk.OneDec())

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	// Initialize a new Leveragelp Trading Position (MTP).
	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, leverage, poolId)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg)
}

func (k Keeper) ProcessOpenLong(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	// Check if the pool is enabled.
	if !k.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, fmt.Sprintf("poolId: %d", poolId))
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	if msg.CollateralAsset != ptypes.BaseCurrency {
		return nil, types.ErrOnlyBaseCurrencyAllowed
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	// Check minimum liabilities - TODO: enable if required
	// collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	// err = k.CheckMinLiabilities(ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset)
	// if err != nil {
	// 	return nil, err
	// }

	// Borrow the asset the user wants to long.
	// TODO: borrow leveragedAmount - collateralAmount
	// TODO: send collateral coins to MTP address from MTP owner address
	leverageCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)
	borrowCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount.Sub(msg.CollateralAmount))
	err = k.stableKeeper.Borrow(ctx, mtp.GetMTPAddress(), borrowCoin)
	if err != nil {
		return nil, err
	}

	_, shares, err := k.amm.JoinPoolNoSwap(ctx, mtp.GetMTPAddress(), poolId, sdk.OneInt(), sdk.Coins{leverageCoin}, true)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Add(shares)
	if err = k.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.UpdateMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	k.CalcMTPConsolidateCollateral(ctx, mtp)

	// Calculate consolidate liabiltiy
	k.CalcMTPConsolidateLiability(ctx, mtp)

	// Set MTP
	k.SetMTP(ctx, mtp)

	return mtp, nil
}
