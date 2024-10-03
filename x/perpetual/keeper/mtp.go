package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	gomath "math"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	if !store.Has(key) {
		return mtp, types.ErrMTPDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &mtp)
	ammPool, found := k.amm.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return mtp, nil
	}
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return mtp, nil
	}
	baseCurrency := entry.Denom

	mtpHealth, err := k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
	if err == nil {
		mtp.MtpHealth = mtpHealth
	}

	pendingBorrowInterest := k.GetBorrowInterest(ctx, &mtp, ammPool)
	mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(pendingBorrowInterest)

	return mtp, nil
}

func (k Keeper) DoesMTPExist(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) bool {
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

func (k Keeper) GetAllLegacyMTPs(ctx sdk.Context) []types.LegacyMTP {
	var mtpList []types.LegacyMTP
	iterator := k.GetMTPIterator(ctx)
	defer func(iterator storetypes.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var mtp types.LegacyMTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		mtpList = append(mtpList, mtp)
	}
	return mtpList
}

func (k Keeper) GetMTPs(ctx sdk.Context, pagination *query.PageRequest) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	var mtpList []*types.MtpAndPrice
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: types.MaxPageLimit,
		}
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

		ammPool, found := k.amm.GetPool(ctx, mtp.AmmPoolId)
		if !found {
			realTime = false
		}

		pnl := math.LegacyZeroDec()
		if realTime {
			mtpHealth, err := k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
			if err == nil {
				mtp.MtpHealth = mtpHealth
			}

			pendingBorrowInterest := k.GetBorrowInterest(ctx, &mtp, ammPool)
			mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(pendingBorrowInterest)
			pnl = k.GetPnL(ctx, mtp, ammPool, baseCurrency)
		}

		info, found := k.oracleKeeper.GetAssetInfo(ctx, mtp.TradingAsset)
		if !found {
			return fmt.Errorf("asset not found")
		}
		trading_asset_price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
		asset_price := math.LegacyZeroDec()
		// If not found set trading_asset_price to zero
		if found {
			asset_price = trading_asset_price.Price
		}

		mtpList = append(mtpList, &types.MtpAndPrice{
			Mtp:               &mtp,
			TradingAssetPrice: asset_price,
			Pnl:               pnl,
		})
		return nil
	})

	return mtpList, pageRes, err
}

func (k Keeper) GetMTPsForPool(ctx sdk.Context, ammPoolId uint64, pagination *query.PageRequest) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	var mtps []*types.MtpAndPrice

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: gomath.MaxUint64 - 1,
		}
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	realTime := true
	if !found {
		realTime = false
	}
	baseCurrency := entry.Denom

	ammPool, found := k.amm.GetPool(ctx, ammPoolId)
	if !found {
		realTime = false
	}

	pageRes, err := query.FilteredPaginate(mtpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		pnl := math.LegacyZeroDec()
		if accumulate && mtp.AmmPoolId == ammPoolId {
			if realTime {
				// Interest
				mtpHealth, err := k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
				if err == nil {
					mtp.MtpHealth = mtpHealth
				}

				pendingBorrowInterest := k.GetBorrowInterest(ctx, &mtp, ammPool)
				mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(pendingBorrowInterest)
				pnl = k.GetPnL(ctx, mtp, ammPool, baseCurrency)
			}

			info, found := k.oracleKeeper.GetAssetInfo(ctx, mtp.TradingAsset)
			if !found {
				return false, fmt.Errorf("asset not found")
			}
			trading_asset_price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
			asset_price := math.LegacyZeroDec()
			// If not found set trading_asset_price to zero
			if found {
				asset_price = trading_asset_price.Price
			}
			mtps = append(mtps, &types.MtpAndPrice{
				Mtp:               &mtp,
				TradingAssetPrice: asset_price,
				Pnl:               pnl,
			})
			return true, nil
		}

		return false, nil
	})

	return mtps, pageRes, err
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

func (k Keeper) GetMTPsForAddressWithPagination(ctx sdk.Context, mtpAddress sdk.AccAddress, pagination *query.PageRequest) ([]*types.MtpAndPrice, *query.PageResponse, error) {
	var mtps []*types.MtpAndPrice

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	mtpStore := prefix.NewStore(store, types.GetMTPPrefixForAddress(mtpAddress))

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
		ammPool, found := k.amm.GetPool(ctx, mtp.AmmPoolId)
		if !found {
			realTime = false
		}

		pnl := math.LegacyZeroDec()
		if realTime {
			mtpHealth, err := k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
			if err == nil {
				mtp.MtpHealth = mtpHealth
			}

			pendingBorrowInterest := k.GetBorrowInterest(ctx, &mtp, ammPool)
			mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(pendingBorrowInterest)
			pnl = k.GetPnL(ctx, mtp, ammPool, baseCurrency)
		}

		info, found := k.oracleKeeper.GetAssetInfo(ctx, mtp.TradingAsset)
		if !found {
			return fmt.Errorf("asset not found")
		}
		trading_asset_price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
		asset_price := math.LegacyZeroDec()
		// If not found set trading_asset_price to zero
		if found {
			asset_price = trading_asset_price.Price
		}

		mtps = append(mtps, &types.MtpAndPrice{
			Mtp:               &mtp,
			TradingAssetPrice: asset_price,
			Pnl:               pnl,
		})
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return mtps, pageRes, nil
}

// Set MTP count
func (k Keeper) SetMTPCount(ctx sdk.Context, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.MTPCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetMTPCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.MTPCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

// Set Open MTP count
func (k Keeper) SetOpenMTPCount(ctx sdk.Context, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.OpenMTPCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetOpenMTPCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.OpenMTPCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) SetToPay(ctx sdk.Context, toPay *types.ToPay) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	address := sdk.MustAccAddressFromBech32(toPay.Address)

	key := types.GetToPayKey(address, toPay.Id)
	store.Set(key, k.cdc.MustMarshal(toPay))
	return nil
}

func (k Keeper) GetAllToPayStore(ctx sdk.Context) []types.ToPay {
	var toPays []types.ToPay
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.ToPayPrefix)
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetToPayKey(address, id)
	if !store.Has(key) {
		return types.ErrToPayDoesNotExist
	}
	store.Delete(key)
	return nil
}

func (k Keeper) GetPnL(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) math.LegacyDec {
	// P&L = Custody (in USD) - Liability ( in USD) - Collateral ( in USD)
	// Liability should include margin interest and funding fee accrued.
	totalLiability := mtp.Liabilities

	pendingBorrowInterest := k.GetBorrowInterest(ctx, &mtp, ammPool)
	mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(pendingBorrowInterest)

	// if short position, convert liabilities to base currency
	if mtp.Position == types.Position_SHORT {
		liabilities := sdk.NewCoin(mtp.LiabilitiesAsset, totalLiability)
		var err error
		totalLiability, err = k.EstimateSwapGivenOut(ctx, liabilities, baseCurrency, ammPool)
		if err != nil {
			totalLiability = math.ZeroInt()
		}
	}

	collateral := mtp.Collateral.Add(mtp.BorrowInterestUnpaidCollateral)
	// include unpaid borrow interest in debt
	if collateral.IsPositive() {
		unpaidCollateral := sdk.NewCoin(mtp.CollateralAsset, collateral)

		if mtp.CollateralAsset == baseCurrency {
			totalLiability = totalLiability.Add(collateral)
		} else {
			C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateral, baseCurrency, ammPool)
			if err != nil {
				C = math.ZeroInt()
			}

			totalLiability = totalLiability.Add(C)
		}
	}

	// Funding rate payment consideration
	// get funding rate
	fundingRate, _, _ := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.AmmPoolId)
	var takeAmountCustodyAmount math.Int
	// if funding rate is zero, return
	if fundingRate.IsZero() {
		takeAmountCustodyAmount = math.ZeroInt()
	} else if (fundingRate.IsNegative() && mtp.Position == types.Position_LONG) || (fundingRate.IsPositive() && mtp.Position == types.Position_SHORT) {
		takeAmountCustodyAmount = math.ZeroInt()
	} else {
		// Calculate the take amount in custody asset
		takeAmountCustodyAmount = types.CalcTakeAmount(mtp.Custody, fundingRate)
	}

	// if short position, custody asset is already in base currency
	custodyAmtInBaseCurrency := mtp.Custody.Sub(takeAmountCustodyAmount)

	if mtp.Position == types.Position_LONG {
		custodyAmt := sdk.NewCoin(mtp.CustodyAsset, mtp.Custody)
		var err error
		custodyAmtInBaseCurrency, err = k.EstimateSwapGivenOut(ctx, custodyAmt, baseCurrency, ammPool)
		if err != nil {
			custodyAmtInBaseCurrency = math.ZeroInt()
		}
	}

	return custodyAmtInBaseCurrency.ToLegacyDec().Sub(totalLiability.ToLegacyDec())
}

func (k Keeper) DeleteLegacyMTP(ctx sdk.Context, mtpaddress string, id uint64) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetMTPKey(sdk.MustAccAddressFromBech32(mtpaddress), id)
	if !store.Has(key) {
		return types.ErrMTPDoesNotExist
	}
	store.Delete(key)
	return nil
}

func (k Keeper) GetAllLegacyMTP(ctx sdk.Context) []types.LegacyMTP {
	var mtps []types.LegacyMTP
	iterator := k.GetMTPIterator(ctx)
	defer func(iterator storetypes.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var mtp types.LegacyMTP
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &mtp)
		if err == nil {
			mtps = append(mtps, mtp)
		}
	}

	return mtps
}

func (k Keeper) DeleteAllNegativeCustomMTP(ctx sdk.Context) {
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var mtp types.LegacyMTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)

		if mtp.Custody.IsNegative() || mtp.Custody.IsZero() {
			k.DestroyMTP(ctx, sdk.MustAccAddressFromBech32(mtp.Address), mtp.Id)
		}
	}

}
