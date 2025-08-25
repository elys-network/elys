package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) AddCollateral(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, collateral sdk.Coin, ammPool *ammtypes.Pool) (sdk.Coin, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.Coin{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	if mtp.Position == types.Position_LONG && mtp.CollateralAsset == baseCurrency {
		if collateral.Denom != mtp.CollateralAsset {
			return sdk.Coin{}, errors.New("denom not same as collateral asset")
		}

		if collateral.Denom != mtp.LiabilitiesAsset {
			return sdk.Coin{}, errors.New("denom not same as liabilities asset")
		}

		creator := mtp.GetAccountAddress()

		// interest amount has been paid from custody
		params := k.GetParams(ctx)
		maxAmount := mtp.Liabilities.Sub(params.LongMinimumLiabilityAmount)
		if !maxAmount.IsPositive() {
			return sdk.Coin{}, fmt.Errorf("cannot reduce liabilties less than %s", params.LongMinimumLiabilityAmount.String())
		}

		var finalAmount math.Int
		if collateral.Amount.LT(maxAmount) {
			finalAmount = collateral.Amount
		} else {
			finalAmount = maxAmount
		}

		mtp.Liabilities = mtp.Liabilities.Sub(finalAmount)
		err := pool.UpdateLiabilities(mtp.LiabilitiesAsset, finalAmount, false, mtp.Position)
		if err != nil {
			return sdk.Coin{}, err
		}

		mtp.Collateral = mtp.Collateral.Add(finalAmount)
		err = pool.UpdateCollateral(mtp.CollateralAsset, finalAmount, true, mtp.Position)
		if err != nil {
			return sdk.Coin{}, err
		}

		finalCollateralCoin := sdk.NewCoin(collateral.Denom, finalAmount)
		err = k.SendToAmmPool(ctx, creator, ammPool, sdk.NewCoins(finalCollateralCoin))
		if err != nil {
			return sdk.Coin{}, err
		}

		err = k.SetMTP(ctx, mtp)
		if err != nil {
			return sdk.Coin{}, err
		}
		k.SetPool(ctx, *pool)

		if k.hooks != nil {
			err = k.hooks.AfterPerpetualPositionModified(ctx, *ammPool, *pool, creator)
			if err != nil {
				return sdk.Coin{}, err
			}
		}
		return finalCollateralCoin, nil

	} else {
		msgOpen := types.MsgOpen{
			Creator:         mtp.Address,
			Position:        mtp.Position,
			Leverage:        math.LegacyZeroDec(),
			Collateral:      collateral,
			TakeProfitPrice: math.LegacyZeroDec(),
			StopLossPrice:   math.LegacyZeroDec(),
			PoolId:          mtp.AmmPoolId,
		}
		if err := msgOpen.ValidateBasic(); err != nil {
			return sdk.Coin{}, err
		}
		_, err := k.Open(ctx, &msgOpen)
		if err != nil {
			return sdk.Coin{}, err
		}
		return sdk.Coin{}, nil
	}
}
