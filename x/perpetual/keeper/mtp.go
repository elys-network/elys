package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	if mtp.Id == 0 {
		perpetualCounter := k.GetPerpetualCounter(ctx, mtp.AmmPoolId)
		// increment global id count
		perpetualCounter.Counter++
		mtp.Id = perpetualCounter.Counter
		// increment open mtp count
		perpetualCounter.TotalOpen++
		k.SetPerpetualCounter(ctx, perpetualCounter)
	}

	// TODO Do we need validate MTP every single time we set it?
	if err := mtp.Validate(); err != nil {
		return err
	}
	key := types.GetMTPKey(mtp.GetAccountAddress(), mtp.Id)
	store.Set(key, k.cdc.MustMarshal(mtp))
	return nil
}

func (k Keeper) DestroyMTP(ctx sdk.Context, mtp types.MTP) {
	key := types.GetMTPKey(mtp.GetAccountAddress(), mtp.Id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(key)

	// decrement open mtp count
	perpetualCounter := k.GetPerpetualCounter(ctx, mtp.AmmPoolId)
	perpetualCounter.TotalOpen--
	k.SetPerpetualCounter(ctx, perpetualCounter)
}

func (k Keeper) GetMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) (types.MTP, error) {
	var mtp types.MTP
	key := types.GetMTPKey(mtpAddress, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	if !store.Has(key) {
		return mtp, types.ErrMTPDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &mtp)
	return mtp, nil
}

func (k Keeper) CheckMTPExist(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) bool {
	key := types.GetMTPKey(mtpAddress, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(key)
}

func (k Keeper) GetMTPIterator(ctx sdk.Context) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return storetypes.KVStorePrefixIterator(store, types.MTPPrefix)
}

func (k Keeper) GetAllMTPs(ctx sdk.Context) []types.MTP {
	var mtpList []types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer func(iterator storetypes.Iterator) {
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var mtpStore storetypes.KVStore

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
	if !found {
		return nil, nil, status.Error(codes.NotFound, "base currency not found")
	}
	baseCurrency := entry.Denom

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)

		if ammPoolId != nil && mtp.AmmPoolId != *ammPoolId {
			return nil
		}

		mtpAndPrice, err := k.fillMTPData(ctx, mtp, baseCurrency)
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

func (k Keeper) fillMTPData(ctx sdk.Context, mtp types.MTP, baseCurrency string) (*types.MtpAndPrice, error) {
	ammPool, found := k.amm.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return &types.MtpAndPrice{}, fmt.Errorf("amm pool %d not found", mtp.AmmPoolId)
	}

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return &types.MtpAndPrice{}, fmt.Errorf("perpetual pool %d not found", mtp.AmmPoolId)
	}

	// Update interest first and then calculate health
	err := k.UpdateMTPBorrowInterestUnpaidLiability(ctx, &mtp)
	if err != nil {
		return &types.MtpAndPrice{}, err
	}
	_, _, _, err = k.UpdateFundingFee(ctx, &mtp, &pool)
	if err != nil {
		return nil, err
	}

	mtp.MtpHealth, err = k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}
	pnl, err := k.GetEstimatedPnL(ctx, mtp, baseCurrency, false)
	if err != nil {
		return nil, err
	}
	liquidationPrice, err := k.GetLiquidationPrice(ctx, mtp)
	if err != nil {
		return nil, err
	}

	tradingAssetPrice, tradingAssetPriceDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	// TODO: replace custody amount with liability amount when fees are defined in terms of liability asset
	// calculate total fees in base currency using asset price
	totalFeesInBaseCurrency := mtp.BorrowInterestPaidCustody.Add(mtp.FundingFeePaidCustody)
	borrowInterestFeesInBaseCurrency := mtp.BorrowInterestPaidCustody
	fundingFeesInBaseCurrency := mtp.FundingFeePaidCustody

	if mtp.Position == types.Position_LONG {
		totalFeesInBaseCurrency = osmomath.BigDecFromSDKInt(totalFeesInBaseCurrency).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
		borrowInterestFeesInBaseCurrency = osmomath.BigDecFromSDKInt(borrowInterestFeesInBaseCurrency).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
		fundingFeesInBaseCurrency = osmomath.BigDecFromSDKInt(fundingFeesInBaseCurrency).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
	}

	effectiveLeverage, err := k.GetEffectiveLeverage(ctx, mtp)
	if err != nil {
		return nil, err
	}

	// Show updated liability
	mtp.Liabilities = mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)

	return &types.MtpAndPrice{
		Mtp:               &mtp,
		TradingAssetPrice: tradingAssetPrice,
		Pnl:               sdk.Coin{baseCurrency, pnl},
		LiquidationPrice:  liquidationPrice,
		EffectiveLeverage: effectiveLeverage,
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

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.GetMTPPrefixForAddress(mtpAddress))

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

// Delete all to pay if any
func (k Keeper) DeleteAllToPay(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{0x09})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	return nil
}

func (k Keeper) GetEstimatedPnL(ctx sdk.Context, mtp types.MTP, baseCurrency string, useTakeProfitPrice bool) (math.Int, error) {

	if useTakeProfitPrice && !mtp.TakeProfitPrice.IsPositive() {
		return math.ZeroInt(), nil
	}

	// P&L = Custody (in USD) - Total Liability ( in USD) - Collateral (in USD)
	// Liability should include margin interest and funding fee accrued.
	collateralAmt := mtp.Collateral

	var tradingAssetPriceDenomRatio osmomath.BigDec
	var err error
	if useTakeProfitPrice {
		tradingAssetPriceDenomRatio, err = k.ConvertPriceToAssetUsdcDenomRatio(ctx, mtp.TradingAsset, mtp.TakeProfitPrice)
		if err != nil {
			return math.Int{}, err
		}
	} else {
		_, tradingAssetPriceDenomRatio, err = k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
		if err != nil {
			return math.Int{}, err
		}
	}

	// in long it's in trading asset ,if short position, custody asset is already in base currency
	custodyAmtAfterFunding := mtp.Custody

	totalLiabilities := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)

	// Calculate estimated PnL
	var estimatedPnL math.Int

	if mtp.Position == types.Position_SHORT {
		// Estimated PnL for short position:
		// collateral asset is in base currency, custody asset is in base currency but liabilities is in trading asset
		// estimated_pnl = custody_amount - totalLiabilities * market_price - collateral_amount

		// For short position, convert liabilities to base currency
		totalLiabilitiesInBaseCurrency := osmomath.BigDecFromSDKInt(totalLiabilities).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
		estimatedPnL = custodyAmtAfterFunding.Sub(totalLiabilitiesInBaseCurrency).Sub(collateralAmt)
	} else {
		// Estimated PnL for long position:
		// collateral asset can be base currency or trading asset, custody asset is in trading asset and liabilities is in base currency
		if mtp.CollateralAsset != baseCurrency {
			// estimated_pnl = (custody_amount - collateral_amount) * market_price - totalLiabilities

			// For long position, convert both custody and collateral to base currency
			custodyAfterCollateralInBaseCurrency := osmomath.BigDecFromSDKInt(custodyAmtAfterFunding.Sub(collateralAmt)).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
			estimatedPnL = custodyAfterCollateralInBaseCurrency.Sub(totalLiabilities)
		} else {
			// estimated_pnl = custody_amount * market_price - totalLiabilities - collateral_amount

			// For long position, convert custody to base currency
			custodyAmountOutInBaseCurrency := osmomath.BigDecFromSDKInt(custodyAmtAfterFunding).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
			estimatedPnL = custodyAmountOutInBaseCurrency.Sub(totalLiabilities).Sub(collateralAmt)
		}
	}

	return estimatedPnL, nil
}

func (k Keeper) GetLiquidationPrice(ctx sdk.Context, mtp types.MTP) (math.LegacyDec, error) {
	liquidationPriceDenomRatio := osmomath.ZeroBigDec()
	params := k.GetParams(ctx)
	// calculate liquidation price
	if mtp.Position == types.Position_LONG {
		// liquidation_price = (safety_factor * liabilities) / custody
		if !mtp.Custody.IsZero() {
			liquidationPriceDenomRatio = params.GetBigDecSafetyFactor().Mul(mtp.GetBigDecLiabilities()).Quo(mtp.GetBigDecCustody())
		}
	}
	if mtp.Position == types.Position_SHORT {
		// liquidation_price =  Custody / (Liabilities * safety_factor)
		if !mtp.Liabilities.IsZero() {
			liquidationPriceDenomRatio = mtp.GetBigDecCustody().Quo(mtp.GetBigDecLiabilities().Mul(params.GetBigDecSafetyFactor()))
		}
	}

	liquidationPrice, err := k.ConvertDenomRatioPriceToUSDPrice(ctx, liquidationPriceDenomRatio, mtp.TradingAsset)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	return liquidationPrice, nil
}

func (k Keeper) CalcMTPTakeProfitCustody(ctx sdk.Context, mtp types.MTP) (math.Int, error) {
	if mtp.IsTakeProfitPriceInfinite() || mtp.TakeProfitPrice.IsZero() {
		return math.ZeroInt(), nil
	}
	takeProfitPriceInDenomRatio, err := k.ConvertPriceToAssetUsdcDenomRatio(ctx, mtp.TradingAsset, mtp.TakeProfitPrice)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("error converting price to base units, asset info %s not found", ptypes.BaseCurrency)
	}
	if mtp.Position == types.Position_LONG {
		return osmomath.BigDecFromSDKInt(mtp.Liabilities).Quo(takeProfitPriceInDenomRatio).Dec().TruncateInt(), nil
	} else {
		return osmomath.BigDecFromSDKInt(mtp.Liabilities).Mul(takeProfitPriceInDenomRatio).Dec().TruncateInt(), nil
	}
}

func (k Keeper) UpdateMTPTakeProfitBorrowFactor(ctx sdk.Context, mtp *types.MTP) error {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return types.ErrZeroCustodyAmount
	}

	// infinite for long, 0 for short
	if mtp.IsTakeProfitPriceInfinite() || mtp.TakeProfitPrice.IsZero() {
		mtp.TakeProfitBorrowFactor = math.LegacyOneDec()
		return nil
	}

	takeProfitPriceDenomRatio, err := k.ConvertPriceToAssetUsdcDenomRatio(ctx, mtp.TradingAsset, mtp.TakeProfitPrice)
	if err != nil {
		return err
	}
	takeProfitBorrowFactor := osmomath.OneBigDec()
	if mtp.Position == types.Position_LONG {
		// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
		takeProfitBorrowFactor = osmomath.OneBigDec().Sub(mtp.GetBigDecLiabilities().Quo(mtp.GetBigDecCustody().Mul(takeProfitPriceDenomRatio)))
	} else {
		// takeProfitBorrowFactor = 1 - ((liabilities  * take profit price) / custody)
		takeProfitBorrowFactor = osmomath.OneBigDec().Sub((mtp.GetBigDecLiabilities().Mul(takeProfitPriceDenomRatio)).Quo(mtp.GetBigDecCustody()))
	}

	mtp.TakeProfitBorrowFactor = takeProfitBorrowFactor.Dec()
	return nil
}
