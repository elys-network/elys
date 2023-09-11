package keeper

import (
	"fmt"
	"math/big"

	"math"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		types.AuthorizationChecker
		types.PositionChecker
		types.PoolChecker
		types.OpenChecker
		types.OpenLongChecker
		types.CloseLongChecker
		types.CloseShortChecker
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		memKey       storetypes.StoreKey
		authority    string
		amm          types.AmmKeeper
		bankKeeper   types.BankKeeper
		oracleKeeper ammtypes.OracleKeeper

		hooks types.MarginHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	authority string,
	amm types.AmmKeeper,
	bk types.BankKeeper,
	oracleKeeper ammtypes.OracleKeeper,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	keeper := &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		memKey:       memKey,
		authority:    authority,
		amm:          amm,
		bankKeeper:   bk,
		oracleKeeper: oracleKeeper,
	}

	keeper.AuthorizationChecker = keeper
	keeper.PositionChecker = keeper
	keeper.PoolChecker = keeper
	keeper.OpenChecker = keeper
	keeper.OpenLongChecker = keeper
	keeper.CloseLongChecker = keeper
	keeper.CloseShortChecker = keeper

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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

func (k Keeper) CheckIfWhitelisted(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWhitelistKey(address))
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (sdk.Int, error) {
	marginEnabled := k.IsPoolEnabled(ctx, ammPool.PoolId)
	if !marginEnabled {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrMarginDisabled, "Margin disabled pool")
	}

	tokensOut := sdk.Coins{tokenOutAmount}
	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	swapResult, err := k.amm.CalcInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, sdk.ZeroDec())

	if err != nil {
		return sdk.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return sdk.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}

func (k Keeper) Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Int, custodyAmount sdk.Int, mtp *types.MTP, ammPool *ammtypes.Pool, pool *types.Pool, eta sdk.Dec) error {
	mtpAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return err
	}
	collateralCoin := sdk.NewCoin(collateralAsset, collateralAmount)

	if !k.bankKeeper.HasBalance(ctx, mtpAddress, collateralCoin) {
		return types.ErrBalanceNotAvailable
	}

	collateralAmountDec := sdk.NewDecFromBigInt(collateralAmount.BigInt())
	liabilitiesDec := collateralAmountDec.Mul(eta)

	// If collateral asset is not usdc, should calculate liability in usdc with the given out.
	if collateralAsset != ptypes.USDC {
		// ATOM amount
		etaAmt := liabilitiesDec.TruncateInt()
		etaAmtToken := sdk.NewCoin(collateralAsset, etaAmt)
		// Calculate usdc amount given atom out amount and we use it liabilty amount in usdc
		liabilityAmt, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, etaAmtToken, ptypes.USDC, *ammPool)
		if err != nil {
			return err
		}

		liabilitiesDec = sdk.NewDecFromInt(liabilityAmt)
	}

	mtp.CollateralAmount = mtp.CollateralAmount.Add(collateralAmount)
	mtp.Liabilities = mtp.Liabilities.Add(sdk.NewIntFromBigInt(liabilitiesDec.TruncateInt().BigInt()))
	mtp.CustodyAmount = mtp.CustodyAmount.Add(custodyAmount)
	mtp.Leverage = eta.Add(sdk.OneDec())

	// print mtp.CustodyAmount
	ctx.Logger().Info(fmt.Sprintf("mtp.CustodyAmount: %s", mtp.CustodyAmount.String()))

	h, err := k.UpdateMTPHealth(ctx, *mtp, *ammPool) // set mtp in func or return h?
	if err != nil {
		return err
	}
	mtp.MtpHealth = h

	ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
	if err != nil {
		return err
	}

	collateralCoins := sdk.NewCoins(collateralCoin)
	err = k.bankKeeper.SendCoins(ctx, mtpAddress, ammPoolAddr, collateralCoins)

	if err != nil {
		return err
	}

	err = pool.UpdateBalance(ctx, collateralAsset, collateralAmount, true)
	if err != nil {
		return err
	}

	// All liability has to be in usdc
	err = pool.UpdateLiabilities(ctx, ptypes.USDC, mtp.Liabilities, true)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return k.SetMTP(ctx, mtp)
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) error {
	pool.Health = k.CalculatePoolHealth(ctx, pool)
	k.SetPool(ctx, *pool)

	return nil
}

func (k Keeper) CalculatePoolHealth(ctx sdk.Context, pool *types.Pool) sdk.Dec {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return sdk.ZeroDec()
	}

	H := sdk.NewDec(1)
	for _, asset := range pool.PoolAssets {
		ammBalance, err := k.GetAmmPoolBalance(ctx, ammPool, asset.AssetDenom)
		if err != nil {
			return sdk.ZeroDec()
		}

		balance := sdk.NewDecFromInt(asset.AssetBalance.Add(ammBalance))
		liabilities := sdk.NewDecFromInt(asset.Liabilities)

		if balance.Add(liabilities).IsZero() {
			return sdk.ZeroDec()
		}

		mul := balance.Quo(balance.Add(liabilities))
		H = H.Mul(mul)
	}

	return H
}

func (k Keeper) TakeInCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool) error {
	err := pool.UpdateBalance(ctx, mtp.CustodyAsset, mtp.CustodyAmount, false)
	if err != nil {
		return nil
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, mtp.CustodyAmount, true)
	if err != nil {
		return nil
	}

	k.SetPool(ctx, *pool)

	return nil
}

func (k Keeper) IncrementalInterestPayment(ctx sdk.Context, interestPayment sdk.Int, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) (sdk.Int, error) {
	// if mtp has unpaid interest, add to payment
	if mtp.InterestUnpaidCollateral.GT(sdk.ZeroInt()) {
		interestPayment = interestPayment.Add(mtp.InterestUnpaidCollateral)
	}

	interestPaymentTokenIn := sdk.NewCoin(mtp.CollateralAsset, interestPayment)
	// swap interest payment to custody asset for payment
	interestPaymentCustody, err := k.EstimateSwap(ctx, interestPaymentTokenIn, mtp.CustodyAsset, ammPool)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// if paying unpaid interest reset to 0
	mtp.InterestUnpaidCollateral = sdk.ZeroInt()

	// edge case, not enough custody to cover payment
	if interestPaymentCustody.GT(mtp.CustodyAmount) {
		// swap custody amount to collateral for updating interest unpaid
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.CustodyAmount)
		custodyAmountCollateral, err := k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.CollateralAsset, ammPool) // may need spot price here to not deduct fee
		if err != nil {
			return sdk.ZeroInt(), err
		}
		mtp.InterestUnpaidCollateral = interestPayment.Sub(custodyAmountCollateral)
		interestPayment = custodyAmountCollateral
		interestPaymentCustody = mtp.CustodyAmount
	}

	// add payment to total paid - collateral
	mtp.InterestPaidCollateral = mtp.InterestPaidCollateral.Add(interestPayment)

	// add payment to total paid - custody
	mtp.InterestPaidCustody = mtp.InterestPaidCustody.Add(interestPaymentCustody)

	// deduct interest payment from custody amount
	mtp.CustodyAmount = mtp.CustodyAmount.Sub(interestPaymentCustody)

	takePercentage := k.GetIncrementalInterestPaymentFundPercentage(ctx)
	fundAddr := k.GetIncrementalInterestPaymentFundAddress(ctx)
	takeAmount, err := k.TakeFundPayment(ctx, interestPaymentCustody, mtp.CustodyAsset, takePercentage, fundAddr, &ammPool)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	actualInterestPaymentCustody := interestPaymentCustody.Sub(takeAmount)

	if !takeAmount.IsZero() {
		k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CustodyAsset, types.EventIncrementalPayFund)
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, interestPaymentCustody, false)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	err = pool.UpdateBalance(ctx, mtp.CustodyAsset, actualInterestPaymentCustody, true)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	k.SetPool(ctx, *pool)

	return actualInterestPaymentCustody, nil
}

func (k Keeper) InterestRateComputation(ctx sdk.Context, pool types.Pool, ammPool ammtypes.Pool) (sdk.Dec, error) {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return sdk.ZeroDec(), sdkerrors.Wrap(types.ErrBalanceNotAvailable, "Balance not available")
	}

	interestRateMax := k.GetInterestRateMax(ctx)
	interestRateMin := k.GetInterestRateMin(ctx)
	interestRateIncrease := k.GetInterestRateIncrease(ctx)
	interestRateDecrease := k.GetInterestRateDecrease(ctx)
	healthGainFactor := k.GetHealthGainFactor(ctx)

	prevInterestRate := pool.InterestRate

	targetInterestRate := healthGainFactor
	for _, asset := range pool.PoolAssets {
		ammBalance, err := k.GetAmmPoolBalance(ctx, ammPool, asset.AssetDenom)
		if err != nil {
			return sdk.ZeroDec(), err
		}

		balance := sdk.NewDecFromInt(asset.AssetBalance.Add(ammBalance))
		liabilities := sdk.NewDecFromInt(asset.Liabilities)

		if balance.Add(liabilities).IsZero() {
			return sdk.ZeroDec(), err
		}

		mul := balance.Add(liabilities).Quo(balance)
		targetInterestRate = targetInterestRate.Mul(mul)
	}

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.GTE(interestRateDecrease.Mul(sdk.NewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
		interestRate = targetInterestRate
	} else if interestRateChange.GT(interestRateIncrease) {
		interestRate = prevInterestRate.Add(interestRateIncrease)
	} else if interestRateChange.LT(interestRateDecrease.Mul(sdk.NewDec(-1))) {
		interestRate = prevInterestRate.Sub(interestRateDecrease)
	}

	newInterestRate := interestRate

	if interestRate.GT(interestRateMin) && interestRate.LT(interestRateMax) {
		newInterestRate = interestRate
	} else if interestRate.LTE(interestRateMin) {
		newInterestRate = interestRateMin
	} else if interestRate.GTE(interestRateMax) {
		newInterestRate = interestRateMax
	}

	return newInterestRate, nil
}

func (k Keeper) CheckMinLiabilities(ctx sdk.Context, collateralAmount sdk.Coin, eta sdk.Dec, pool types.Pool, ammPool ammtypes.Pool, custodyAsset string) error {
	var interestRational, liabilitiesRational, rate big.Rat
	minInterestRate := k.GetInterestRateMin(ctx)

	collateralAmountDec := sdk.NewDecFromInt(collateralAmount.Amount)
	liabilitiesDec := collateralAmountDec.Mul(eta)
	liabilities := sdk.NewUint(liabilitiesDec.TruncateInt().Uint64())

	// In Long position, liabilty has to be always in USDC
	if collateralAmount.Denom != ptypes.USDC {
		outAmt := liabilitiesDec.TruncateInt()
		outAmtToken := sdk.NewCoin(collateralAmount.Denom, outAmt)
		inAmt, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, outAmtToken, ptypes.USDC, ammPool)
		if err != nil {
			return types.ErrBorrowTooLow
		}
		liabilities = sdk.NewUint(inAmt.Uint64())
	}
	rate.SetFloat64(minInterestRate.MustFloat64())
	liabilitiesRational.SetInt(liabilities.BigInt())
	interestRational.Mul(&rate, &liabilitiesRational)

	interestNew := interestRational.Num().Quo(interestRational.Num(), interestRational.Denom())
	samplePayment := sdk.NewInt(interestNew.Int64())

	if samplePayment.IsZero() && !minInterestRate.IsZero() {
		return types.ErrBorrowTooLow
	}

	// If collateral is not usdc, custody amount is already checked in HasSufficientBalance function.
	// its liability balance checked in the above if statement, so return
	if collateralAmount.Denom != ptypes.USDC {
		return nil
	}

	samplePaymentTokenIn := sdk.NewCoin(collateralAmount.Denom, samplePayment)
	// swap interest payment to custody asset
	_, err := k.EstimateSwap(ctx, samplePaymentTokenIn, custodyAsset, ammPool)
	if err != nil {
		return types.ErrBorrowTooLow
	}

	return nil
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

func (k Keeper) TakeFundPayment(ctx sdk.Context, returnAmount sdk.Int, returnAsset string, takePercentage sdk.Dec, fundAddr sdk.AccAddress, ammPool *ammtypes.Pool) (sdk.Int, error) {
	returnAmountDec := sdk.NewDecFromBigInt(returnAmount.BigInt())
	takeAmount := sdk.NewIntFromBigInt(takePercentage.Mul(returnAmountDec).TruncateInt().BigInt())

	if !takeAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, sdk.NewIntFromBigInt(takeAmount.BigInt())))
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ammPool.Address, fundAddr, takeCoins)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	}
	return takeAmount, nil
}

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
	key := types.GetMTPKey(mtp.Address, mtp.Id)
	store.Set(key, k.cdc.MustMarshal(mtp))
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

func (k Keeper) GetWhitelistAddressIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WhitelistPrefix)
}

func (k Keeper) GetAllWhitelistedAddress(ctx sdk.Context) []string {
	var list []string
	iterator := k.GetWhitelistAddressIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, (string)(iterator.Value()))
	}

	return list
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

func (k Keeper) GetWhitelistedAddress(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error) {
	var list []string
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.WhitelistPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(prefixStore, pagination, func(key []byte, value []byte) error {
		list = append(list, string(value))
		return nil
	})

	return list, pageRes, err
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWhitelistKey(address), []byte(address))
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWhitelistKey(address))
}

// Set the margin hooks.
func (k *Keeper) SetHooks(gh types.MarginHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set margin hooks twice")
	}

	k.hooks = gh

	return k
}
