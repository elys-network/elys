package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Int, error) {
	// Retrieve MTP
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Retrieve Pool
	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Exit liquidity with collateral token
	exitCoins, err := k.amm.ExitPool(ctx, mtp.GetMTPAddress(), mtp.AmmPoolId, mtp.LeveragedLpAmount, sdk.Coins{}, mtp.Collateral.Denom)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Repay with interest
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, mtp.GetMTPAddress())
	repayAmount := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	err = k.stableKeeper.Repay(ctx, mtp.GetMTPAddress(), sdk.NewCoin(mtp.Collateral.Denom, repayAmount))
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	userAmount := exitCoins[0].Amount.Sub(repayAmount)
	mtpOwner := sdk.MustAccAddressFromBech32(mtp.Address)
	err = k.bankKeeper.SendCoins(ctx, mtpOwner, mtp.GetMTPAddress(), sdk.Coins{sdk.NewCoin(mtp.Collateral.Denom, userAmount)})
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	err = k.DestroyMTP(ctx, mtp.Address, mtp.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(mtp.LeveragedLpAmount)
	k.SetPool(ctx, pool)

	// Hooks after leveragelp position closed
	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionClosed(ctx, ammPool, pool)
	}

	return &mtp, repayAmount, nil
}
