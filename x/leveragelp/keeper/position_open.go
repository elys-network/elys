package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	// Initialize a new Leveragelp Trading Position (MTP).
	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, leverage, poolId)
	mtp.Id = k.GetMTPCount(ctx) + 1
	k.SetMTPCount(ctx, mtp.Id)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, mtp, leverage, collateralAmountDec, poolId, msg)
}

func (k Keeper) OpenConsolidate(ctx sdk.Context, mtp *types.MTP, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {

	poolId := mtp.AmmPoolId
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	if !k.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, fmt.Sprintf("poolId: %d", poolId))
	}

	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	mtp, err = k.OpenConsolidateLong(ctx, poolId, mtp, msg)
	if err != nil {
		return nil, err
	}

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral", mtp.Collateral.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	)
	ctx.EventManager().EmitEvent(event)

	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionModified(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}

func (k Keeper) OpenConsolidateLong(ctx sdk.Context, poolId uint64, mtp *types.MTP, msg *types.MsgOpen) (*types.MTP, error) {

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	mtp.Collateral = mtp.Collateral.Add(sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount))
	// TODO: Leverage won't be required probably
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	mtp.Leverage = leverage

	return k.ProcessOpenLong(ctx, mtp, leverage, collateralAmountDec, poolId, msg)
}

func (k Keeper) ProcessOpenLong(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
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

	// send collateral coins to MTP address from MTP owner address
	mtpOwner := sdk.MustAccAddressFromBech32(mtp.Address)
	err = k.bankKeeper.SendCoins(ctx, mtpOwner, mtp.GetMTPAddress(), sdk.Coins{sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)})
	if err != nil {
		return nil, err
	}
	leverageCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)

	// borrow leveragedAmount - collateralAmount
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
	mtp.Liabilities = borrowCoin.Amount
	lr, err := k.GetMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Set MTP
	k.SetMTP(ctx, mtp)

	return mtp, nil
}
