package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
	store := ctx.KVStore(k.storeKey)
	count := k.GetMTPCount(ctx)
	openCount := k.GetOpenMTPCount(ctx)

	if mtp.Id == 0 {
		// increment global id count
		count++
		mtp.Id = count
		k.SetMTPCount(ctx, count)
		// increment open mtp count
		openCount++
		k.SetOpenMTPCount(ctx, openCount)
	}

	if err := mtp.Validate(); err != nil {
		return err
	}
	key := types.GetMTPKey(mtp.GetAccountAddress(), mtp.Id)
	store.Set(key, k.cdc.MustMarshal(mtp))
	return nil
}

func (k Keeper) DestroyMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) error {
	key := types.GetMTPKey(mtpAddress, id)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.ErrMTPDoesNotExist
	}
	store.Delete(key)
	// decrement open mtp count
	openCount := k.GetOpenMTPCount(ctx)
	openCount--

	// Set open MTP count
	k.SetOpenMTPCount(ctx, openCount)

	return nil
}

func (k Keeper) GetMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) (types.MTP, error) {
	var mtp types.MTP
	key := types.GetMTPKey(mtpAddress, id)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return mtp, types.ErrMTPDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &mtp)
	return mtp, nil
}

func (k Keeper) CheckMTPExist(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) bool {
	key := types.GetMTPKey(mtpAddress, id)
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

func (k Keeper) GetMTPIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.MTPPrefix)
}

func (k Keeper) GetAllMTPs(ctx sdk.Context) []types.MTP {
	var mtpList []types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		mtpList = append(mtpList, mtp)
	}
	return mtpList
}

func (k Keeper) GetMTPData(ctx sdk.Context, pagination *query.PageRequest, address sdk.AccAddress, ammPoolId *uint64) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	var mtps []*types.MtpAndPrice
	store := ctx.KVStore(k.storeKey)
	var mtpStore sdk.KVStore

	if address != nil {
		mtpStore = prefix.NewStore(store, types.GetMTPPrefixForAddress(address))
	} else {
		mtpStore = prefix.NewStore(store, types.MTPPrefix)
	}

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: types.MaxPageLimit,
		}
	}

	if pagination.Limit > types.MaxPageLimit {
		return nil, nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	realTime := true
	if !found {
		realTime = false
	}
	baseCurrency := entry.Denom

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)

		if ammPoolId != nil && mtp.AmmPoolId != *ammPoolId {
			return nil
		}

		mtpAndPrice, err := k.fillMTPData(ctx, mtp, ammPoolId, realTime, baseCurrency)
		if err != nil {
			return err
		}

		mtps = append(mtps, mtpAndPrice)

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return mtps, pageRes, nil
}

func (k Keeper) fillMTPData(ctx sdk.Context, mtp types.MTP, ammPoolId *uint64, realTime bool, baseCurrency string) (*types.MtpAndPrice, error) {
	var ammPool ammtypes.Pool
	var poolFound bool
	if ammPoolId != nil {
		ammPool, poolFound = k.amm.GetPool(ctx, *ammPoolId)
	} else {
		ammPool, poolFound = k.amm.GetPool(ctx, mtp.AmmPoolId)
	}
	if !poolFound {
		realTime = false
	}

	pnl := math.ZeroInt()
	liquidationPrice := sdk.ZeroDec()
	if realTime {
		mtpHealth, err := k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
		if err == nil {
			mtp.MtpHealth = mtpHealth
		}

		k.UpdateMTPBorrowInterestUnpaidLiability(ctx, &mtp)
		pnl, err = k.GetPnL(ctx, mtp, ammPool, baseCurrency)
		if err != nil {
			return nil, err
		}
		liquidationPrice = k.GetLiquidationPrice(ctx, mtp, ammPool, baseCurrency)
	}

	info, found := k.oracleKeeper.GetAssetInfo(ctx, mtp.TradingAsset)
	if !found {
		return nil, fmt.Errorf("asset not found")
	}
	tradingAssetPrice, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	assetPrice := sdk.ZeroDec()
	if found {
		assetPrice = tradingAssetPrice.Price
	}

	// TODO: replace custody amount with liability amount when fees are defined in terms of liability asset
	// calculate total fees in base currency using asset price
	totalFeesInBaseCurrency := mtp.BorrowInterestPaidCustody.Add(mtp.FundingFeePaidCustody)
	borrowInterestFeesInBaseCurrency := mtp.BorrowInterestPaidCustody
	fundingFeesInBaseCurrency := mtp.FundingFeePaidCustody

	if mtp.Position == types.Position_LONG {
		totalFeesInBaseCurrency = totalFeesInBaseCurrency.ToLegacyDec().Mul(assetPrice).TruncateInt()
		borrowInterestFeesInBaseCurrency = borrowInterestFeesInBaseCurrency.ToLegacyDec().Mul(assetPrice).TruncateInt()
		fundingFeesInBaseCurrency = fundingFeesInBaseCurrency.ToLegacyDec().Mul(assetPrice).TruncateInt()
	}

	return &types.MtpAndPrice{
		Mtp:               &mtp,
		TradingAssetPrice: assetPrice,
		Pnl:               pnl,
		LiquidationPrice:  liquidationPrice,
		Fees: &types.Fees{
			TotalFeesBaseCurrency:            totalFeesInBaseCurrency,
			BorrowInterestFeesLiabilityAsset: mtp.BorrowInterestPaidCustody,
			BorrowInterestFeesBaseCurrency:   borrowInterestFeesInBaseCurrency,
			FundingFeesLiquidityAsset:        mtp.FundingFeePaidCustody,
			FundingFeesBaseCurrency:          fundingFeesInBaseCurrency,
		},
	}, nil
}

func (k Keeper) GetAllMTPsForAddress(ctx sdk.Context, mtpAddress sdk.AccAddress) []*types.MTP {
	var mtps []*types.MTP

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetMTPPrefixForAddress(mtpAddress))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		mtps = append(mtps, &mtp)
	}
	return mtps
}

func (k Keeper) GetMTPs(ctx sdk.Context, pagination *query.PageRequest) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	return k.GetMTPData(ctx, pagination, nil, nil)
}

func (k Keeper) GetMTPsForPool(ctx sdk.Context, ammPoolId uint64, pagination *query.PageRequest) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	return k.GetMTPData(ctx, pagination, nil, &ammPoolId)
}

func (k Keeper) GetMTPsForAddressWithPagination(ctx sdk.Context, mtpAddress sdk.AccAddress, pagination *query.PageRequest) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	return k.GetMTPData(ctx, pagination, mtpAddress, nil)
}

// Set MTP count
func (k Keeper) SetMTPCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MTPCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetMTPCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := ctx.KVStore(k.storeKey).Get(types.MTPCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

// Set Open MTP count
func (k Keeper) SetOpenMTPCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OpenMTPCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetOpenMTPCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := ctx.KVStore(k.storeKey).Get(types.OpenMTPCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) SetToPay(ctx sdk.Context, toPay *types.ToPay) error {
	store := ctx.KVStore(k.storeKey)
	address := sdk.MustAccAddressFromBech32(toPay.Address)

	key := types.GetToPayKey(address, toPay.Id)
	store.Set(key, k.cdc.MustMarshal(toPay))
	return nil
}

func (k Keeper) GetAllToPayStore(ctx sdk.Context) []types.ToPay {
	var toPays []types.ToPay
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ToPayPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var toPay types.ToPay
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &toPay)
		toPays = append(toPays, toPay)
	}
	return toPays
}

func (k Keeper) DeleteToPay(ctx sdk.Context, address sdk.AccAddress, id uint64) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetToPayKey(address, id)
	if !store.Has(key) {
		return types.ErrToPayDoesNotExist
	}
	store.Delete(key)
	return nil
}

// TODO might be incorrect, takeFundingFee is being reduced from custody rather than mtp.BorrowInterestUnpaidLiability
func (k Keeper) GetPnL(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) (math.Int, error) {
	// P&L = Custody (in USD) - Liability ( in USD) - Collateral ( in USD)

	// Funding rate payment consideration
	// get funding rate
	fundingRate, _ := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.LastFundingCalcTime, mtp.AmmPoolId)
	var takeAmountCustodyAmount sdk.Int
	// if funding rate is zero, return
	if fundingRate.IsZero() {
		takeAmountCustodyAmount = sdk.ZeroInt()
	} else if (fundingRate.IsNegative() && mtp.Position == types.Position_LONG) || (fundingRate.IsPositive() && mtp.Position == types.Position_SHORT) {
		takeAmountCustodyAmount = sdk.ZeroInt()
	} else {
		// Calculate the take amount in custody asset
		takeAmountCustodyAmount = types.CalcTakeAmount(mtp.Custody, fundingRate)
	}

	// Liability should include margin interest and funding fee accrued.
	collateralAmt := mtp.Collateral

	// in long it's in trading asset ,if short position, custody asset is already in base currency
	custodyAmtAfterTake := mtp.Custody.Sub(takeAmountCustodyAmount)

	totalLiabilities := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)

	// Calculate estimated PnL
	estimatedPnL := sdk.ZeroInt()

	if mtp.Position == types.Position_SHORT {
		// Estimated PnL for short position:
		// collateral asset is in base currency, custody asset is in base currency but liabilities is in trading asset
		// estimated_pnl = custody_amount - totalLiabilities * market_price - collateral_amount

		// For short position, convert liabilities to base currency
		totalLiabilitiesTokenOut := sdk.NewCoin(mtp.LiabilitiesAsset, totalLiabilities)
		totalLiabilitiesInBaseCurrency, _, err := k.EstimateSwapGivenOut(ctx, totalLiabilitiesTokenOut, baseCurrency, ammPool)
		if err != nil {
			totalLiabilitiesInBaseCurrency = sdk.ZeroInt()
		}

		estimatedPnL = custodyAmtAfterTake.Sub(totalLiabilitiesInBaseCurrency).Sub(collateralAmt)
	} else {
		// Estimated PnL for long position:
		// collateral asset can be base currency or trading asset, custody asset is in trading asset and liabilities is in base currency
		if mtp.CollateralAsset != baseCurrency {
			// estimated_pnl = (custody_amount - collateral_amount) * market_price - totalLiabilities

			// For long position, convert both custody and collateral to base currency
			custodyCollateralDiffTokenOut := sdk.NewCoin(mtp.CollateralAsset, sdk.MaxInt(custodyAmtAfterTake.Sub(collateralAmt), sdk.ZeroInt()))
			custodyCollateralDiffInBaseCurrency, _, err := k.EstimateSwapGivenOut(ctx, custodyCollateralDiffTokenOut, baseCurrency, ammPool)
			if err != nil {
				custodyCollateralDiffInBaseCurrency = sdk.ZeroInt()
			}

			estimatedPnL = custodyCollateralDiffInBaseCurrency.Sub(totalLiabilities)
		} else {
			// estimated_pnl = custody_amount * market_price - totalLiabilities - collateral_amount

			// For long position, convert custody to base currency
			custodyAmountOut := sdk.NewCoin(mtp.CustodyAsset, sdk.MaxInt(custodyAmtAfterTake, sdk.ZeroInt()))
			custodyAmountOutInBaseCurrency, _, err := k.EstimateSwapGivenOut(ctx, custodyAmountOut, baseCurrency, ammPool)
			if err != nil {
				custodyAmountOutInBaseCurrency = sdk.ZeroInt()
			}

			estimatedPnL = custodyAmountOutInBaseCurrency.Sub(mtp.Liabilities).Sub(collateralAmt)
		}
	}

	return estimatedPnL, nil
}

func (k Keeper) GetLiquidationPrice(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) sdk.Dec {
	collateralAmountInBaseCurrency := mtp.Collateral
	if mtp.CollateralAsset != baseCurrency {
		amount, _, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral), baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroDec()
		}
		collateralAmountInBaseCurrency = amount
	}

	liabilitiesAmountInBaseCurrency := mtp.Liabilities
	if mtp.LiabilitiesAsset != baseCurrency {
		amount, _, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities), baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroDec()
		}
		liabilitiesAmountInBaseCurrency = amount
	}

	custodyAmountInTradingAsset := mtp.Custody
	if mtp.CustodyAsset != mtp.TradingAsset {
		amount, _, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CustodyAsset, mtp.Custody), mtp.TradingAsset, ammPool)
		if err != nil {
			return sdk.ZeroDec()
		}
		custodyAmountInTradingAsset = amount
	}

	// open price = (collateral + liabilities) / custody
	mtp.OpenPrice = math.LegacyNewDecFromBigInt(
		collateralAmountInBaseCurrency.Add(liabilitiesAmountInBaseCurrency).BigInt(),
	).Quo(
		math.LegacyNewDecFromBigInt(custodyAmountInTradingAsset.BigInt()),
	)

	// calculate liquidation price
	// liquidation_price = open_price_value - collateral_amount / custody_amount
	liquidationPrice := mtp.OpenPrice.Sub(
		sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.BigInt()).Quo(sdk.NewDecFromBigInt(mtp.Custody.BigInt())),
	)

	// if position is short then liquidation price is open price + collateral amount / (custody amount / open price)
	if mtp.Position == types.Position_SHORT {
		positionSizeInTradingAsset := sdk.NewDecFromBigInt(mtp.Custody.BigInt()).Quo(mtp.OpenPrice)
		liquidationPrice = mtp.OpenPrice.Add(
			sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.BigInt()).Quo(positionSizeInTradingAsset),
		)
	}

	return liquidationPrice
}
