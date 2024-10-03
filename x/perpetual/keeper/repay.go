package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, returnAmount math.Int, amount math.Int) (err error) {
	if returnAmount.IsPositive() {
		returnCoin := sdk.NewCoin(mtp.CollateralAsset, sdk.NewIntFromBigInt(returnAmount.BigInt()))
		returnCoins := sdk.NewCoins(returnCoin)
		addr, err := sdk.AccAddressFromBech32(mtp.Address)
		if err != nil {
			return err
		}

		ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoins(ctx, ammPoolAddr, addr, returnCoins)
		if err != nil {
			return err
		}
	}

	err = pool.UpdateBalance(ctx, mtp.CollateralAsset, returnAmount, false, mtp.Position)
	if err != nil {
		return err
	}

	// long position
	err = pool.UpdateLiabilities(ctx, mtp.LiabilitiesAsset, mtp.Liabilities, false, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateTakeProfitLiabilities(ctx, mtp.LiabilitiesAsset, mtp.TakeProfitLiabilities, false, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateTakeProfitCustody(ctx, mtp.CustodyAsset, mtp.TakeProfitCustody, false, mtp.Position)
	if err != nil {
		return err
	}

	mtp.Custody = mtp.Custody.Sub(amount)
	// This is for accounting purposes, mtp.Custody gets reduced by borrowInterestPaymentCustody and funding fee. so msg.Amount is greater than mtp.Custody here. So if it's negative it should be closed
	if mtp.Custody.IsZero() || mtp.Custody.IsNegative() {
		err = k.DestroyMTP(ctx, mtp.GetAccountAddress(), mtp.Id)
		if err != nil {
			return err
		}
	} else {
		err = k.SetMTP(ctx, mtp)
		if err != nil {
			return err
		}
	}

	k.SetPool(ctx, *pool)

	return nil
}
