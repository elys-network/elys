package keeper

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetMTP(ctx sdk.Context, mtpAddress string, id uint64) (types.MTP, error) {
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

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) {
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

	key := types.GetMTPKey(mtp.Address, mtp.Id)
	store.Set(key, k.cdc.MustMarshal(mtp))
}

func (k Keeper) DestroyMTP(ctx sdk.Context, mtpAddress string, id uint64) error {
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

// Set Open MTP count
func (k Keeper) SetOpenMTPCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OpenMTPCountPrefix, types.GetUint64Bytes(count))
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

func (k Keeper) GetMTPs(ctx sdk.Context, pagination *query.PageRequest) ([]*types.MTP, *query.PageResponse, error) {
	var mtpList []*types.MTP
	store := ctx.KVStore(k.storeKey)
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		mtpList = append(mtpList, &mtp)
		return nil
	})

	return mtpList, pageRes, err
}

func (k Keeper) GetMTPsForPool(ctx sdk.Context, ammPoolId uint64, pagination *query.PageRequest) ([]*types.MTP, *query.PageResponse, error) {
	var mtps []*types.MTP

	store := ctx.KVStore(k.storeKey)
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.FilteredPaginate(mtpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		if accumulate && mtp.AmmPoolId == ammPoolId {
			mtps = append(mtps, &mtp)
			return true, nil
		}

		return false, nil
	})

	return mtps, pageRes, err
}

func (k Keeper) GetMTPsForAddress(ctx sdk.Context, mtpAddress sdk.Address, pagination *query.PageRequest) ([]*types.MTP, *query.PageResponse, error) {
	var mtps []*types.MTP

	store := ctx.KVStore(k.storeKey)
	mtpStore := prefix.NewStore(store, types.GetMTPPrefixForAddress(mtpAddress.String()))

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: types.MaxPageLimit,
		}
	}

	if pagination.Limit > types.MaxPageLimit {
		return nil, nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		mtps = append(mtps, &mtp)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return mtps, pageRes, nil
}

func (k Keeper) GetMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool) (sdk.Dec, error) {
	debt := k.stableKeeper.GetDebt(ctx, mtp.GetMTPAddress())
	xl := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}

	commitments, found := k.commKeeper.GetCommitments(ctx, mtp.GetMTPAddress().String())
	if !found {
		return sdk.ZeroDec(), nil
	}

	fmt.Println("commitments", commitments)
	fmt.Println("xl", xl.String())

	mtpVal := sdk.ZeroDec()
	for _, commitment := range commitments.CommittedTokens {
		ammPoolTvl, err := ammPool.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		mtpVal = mtpVal.Add(
			ammPoolTvl.
				Mul(sdkmath.LegacyNewDecFromInt(commitment.Amount)).
				Quo(sdkmath.LegacyNewDecFromInt(ammPool.TotalShares.Amount)))

		fmt.Println("ammPoolTvl", ammPoolTvl.String())
		fmt.Println("commitment.Amount", commitment.Amount.String())
		fmt.Println("ammPool.TotalShares", ammPool.TotalShares.String())
	}

	lr := mtpVal.Quo(sdk.NewDecFromBigInt(xl.BigInt()))
	return lr, nil
}
