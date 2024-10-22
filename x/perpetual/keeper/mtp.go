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
		pnl, err = k.GetEstimatedPnL(ctx, mtp, baseCurrency, false)
		if err != nil {
			return nil, err
		}
		liquidationPrice = k.GetLiquidationPrice(ctx, mtp)
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	// TODO: replace custody amount with liability amount when fees are defined in terms of liability asset
	// calculate total fees in base currency using asset price
	totalFeesInBaseCurrency := mtp.BorrowInterestPaidCustody.Add(mtp.FundingFeePaidCustody)
	borrowInterestFeesInBaseCurrency := mtp.BorrowInterestPaidCustody
	fundingFeesInBaseCurrency := mtp.FundingFeePaidCustody

	if mtp.Position == types.Position_LONG {
		totalFeesInBaseCurrency = totalFeesInBaseCurrency.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
		borrowInterestFeesInBaseCurrency = borrowInterestFeesInBaseCurrency.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
		fundingFeesInBaseCurrency = fundingFeesInBaseCurrency.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
	}

	effectiveLeverage, err := k.UpdatedLeverage(ctx, mtp)
	if err != nil {
		return nil, err
	}

	return &types.MtpAndPrice{
		Mtp:               &mtp,
		TradingAssetPrice: tradingAssetPrice,
		Pnl:               pnl,
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

func (k Keeper) GetEstimatedPnL(ctx sdk.Context, mtp types.MTP, baseCurrency string, useTakeProfitPrice bool) (math.Int, error) {
	// P&L = Custody (in USD) - Total Liability ( in USD) - Collateral ( in USD)

	// Funding rate payment consideration
	// get funding rate
	fundingRate, _ := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.LastFundingCalcTime, mtp.AmmPoolId)
	var fundingAmount sdk.Int
	// if funding rate is zero, return
	if fundingRate.IsZero() {
		fundingAmount = sdk.ZeroInt()
	} else if (fundingRate.IsNegative() && mtp.Position == types.Position_LONG) || (fundingRate.IsPositive() && mtp.Position == types.Position_SHORT) {
		fundingAmount = sdk.ZeroInt()
	} else {
		// Calculate the take amount in custody asset
		fundingAmount = types.CalcTakeAmount(mtp.Custody, fundingRate)
	}

	// Liability should include margin interest and funding fee accrued.
	collateralAmt := mtp.Collateral

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return math.Int{}, err
	}
	if useTakeProfitPrice {
		tradingAssetPrice = mtp.TakeProfitPrice
	}
	if tradingAssetPrice.IsZero() {
		return math.Int{}, fmt.Errorf("trading asset price is zero")
	}

	// in long it's in trading asset ,if short position, custody asset is already in base currency
	custodyAmtAfterFunding := mtp.Custody.Sub(fundingAmount)

	totalLiabilities := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)

	// Calculate estimated PnL
	estimatedPnL := sdk.ZeroInt()

	if mtp.Position == types.Position_SHORT {
		// Estimated PnL for short position:
		// collateral asset is in base currency, custody asset is in base currency but liabilities is in trading asset
		// estimated_pnl = custody_amount - totalLiabilities * market_price - collateral_amount

		// For short position, convert liabilities to base currency
		totalLiabilitiesInBaseCurrency := totalLiabilities.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
		estimatedPnL = custodyAmtAfterFunding.Sub(totalLiabilitiesInBaseCurrency).Sub(collateralAmt)
	} else {
		// Estimated PnL for long position:
		// collateral asset can be base currency or trading asset, custody asset is in trading asset and liabilities is in base currency
		if mtp.CollateralAsset != baseCurrency {
			// estimated_pnl = (custody_amount - collateral_amount) * market_price - totalLiabilities

			// For long position, convert both custody and collateral to base currency
			custodyAfterCollateralInBaseCurrecy := (custodyAmtAfterFunding.Sub(collateralAmt)).ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
			estimatedPnL = custodyAfterCollateralInBaseCurrecy.Sub(totalLiabilities)
		} else {
			// estimated_pnl = custody_amount * market_price - totalLiabilities - collateral_amount

			// For long position, convert custody to base currency
			custodyAmountOutInBaseCurrency := custodyAmtAfterFunding.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
			estimatedPnL = custodyAmountOutInBaseCurrency.Sub(totalLiabilities).Sub(collateralAmt)
		}
	}

	return estimatedPnL, nil
}

func (k Keeper) GetLiquidationPrice(ctx sdk.Context, mtp types.MTP) sdk.Dec {
	liquidationPrice := math.LegacyZeroDec()
	params := k.GetParams(ctx)
	// calculate liquidation price
	if mtp.Position == types.Position_LONG {
		// liquidation_price = (safety_factor * liabilities) / custody
		liquidationPrice = mtp.Liabilities.ToLegacyDec().Quo(params.SafetyFactor.Mul(mtp.Custody.ToLegacyDec()))
	}
	if mtp.Position == types.Position_SHORT {
		// liquidation_price =  Custody / (Liabilities * safety_factor)
		liquidationPrice = mtp.Custody.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec().Mul(params.SafetyFactor))
	}

	return liquidationPrice
}
