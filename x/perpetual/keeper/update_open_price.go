package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateOpenPrice(ctx sdk.Context, mtp *types.MTP, ammPool ammtypes.Pool, baseCurrency string) error {
	collateralAmountInBaseCurrency := mtp.Collateral
	if mtp.CollateralAsset != baseCurrency {
		C, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral), baseCurrency, ammPool)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error estimating swap: %s %s", mtp.CollateralAsset, mtp.Collateral))
		}
		collateralAmountInBaseCurrency = C
	}

	liabilitiesAmountInBaseCurrency := mtp.Liabilities
	if mtp.LiabilitiesAsset != baseCurrency {
		L, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities), baseCurrency, ammPool)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error estimating swap: %s %s", mtp.LiabilitiesAsset, mtp.Liabilities))
		}
		liabilitiesAmountInBaseCurrency = L
	}

	custodyAmountInTradingAsset := mtp.Custody
	if mtp.CustodyAsset != mtp.TradingAsset {
		C, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CustodyAsset, mtp.Custody), mtp.TradingAsset, ammPool)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error estimating swap: %s %s", mtp.CustodyAsset, mtp.Custody))
		}
		custodyAmountInTradingAsset = C
	}

	// open price = (collateral + liabilities) / custody
	mtp.OpenPrice = math.LegacyNewDecFromBigInt(
		collateralAmountInBaseCurrency.Add(liabilitiesAmountInBaseCurrency).BigInt(),
	).Quo(
		math.LegacyNewDecFromBigInt(custodyAmountInTradingAsset.BigInt()),
	)

	err := k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	return nil
}
