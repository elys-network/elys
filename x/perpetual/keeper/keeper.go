package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	pkeeper "github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

type (
	Keeper struct {
		types.AuthorizationChecker
		types.PositionChecker
		types.OpenChecker
		types.OpenDefineAssetsChecker
		types.ClosePositionChecker
		types.CloseEstimationChecker

		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		memKey             storetypes.StoreKey
		authority          string
		parameterKeeper    *pkeeper.Keeper
		amm                types.AmmKeeper
		bankKeeper         types.BankKeeper
		oracleKeeper       types.OracleKeeper
		assetProfileKeeper types.AssetProfileKeeper

		hooks types.PerpetualHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	authority string,
	amm types.AmmKeeper,
	bk types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	parameterKeeper *pkeeper.Keeper,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	keeper := &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		authority:          authority,
		amm:                amm,
		bankKeeper:         bk,
		oracleKeeper:       oracleKeeper,
		assetProfileKeeper: assetProfileKeeper,
		parameterKeeper:    parameterKeeper,
	}

	keeper.AuthorizationChecker = keeper
	keeper.PositionChecker = keeper
	keeper.OpenChecker = keeper
	keeper.OpenDefineAssetsChecker = keeper
	keeper.ClosePositionChecker = keeper
	keeper.CloseEstimationChecker = keeper

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Borrow(ctx sdk.Context, collateralAmount math.Int, custodyAmount math.Int, mtp *types.MTP, ammPool *ammtypes.Pool, pool *types.Pool, eta sdk.Dec, baseCurrency string, isBroker bool) error {
	senderAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return err
	}
	// if isBroker is true, then retrieve broker address and assign it to senderAddress
	if isBroker {
		brokerAddress, err := sdk.AccAddressFromBech32(k.parameterKeeper.GetParams(ctx).BrokerAddress)
		if err != nil {
			return err
		}
		senderAddress = brokerAddress
	}

	collateralCoin := sdk.NewCoin(mtp.CollateralAsset, collateralAmount)

	if !k.bankKeeper.HasBalance(ctx, senderAddress, collateralCoin) {
		return types.ErrBalanceNotAvailable
	}

	liabilitiesInCollateral := collateralAmount.ToLegacyDec().Mul(eta).TruncateInt()
	liabilities := liabilitiesInCollateral
	// If collateral asset is not base currency, should calculate liability in base currency with the given out.
	// For LONG, Liability has to be in base currency, CollateralAsset can be trading asset or base currency
	// For SHORT, Liability has to be in trading asset and CollateralAsset will be in base currency, so this if case only applies to LONG
	if mtp.CollateralAsset != baseCurrency {
		liabilitiesInCollateralTokenOut := sdk.NewCoin(mtp.CollateralAsset, liabilitiesInCollateral)
		// Calculate base currency amount given atom out amount and we use it liabilty amount in base currency
		liabilities, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesInCollateralTokenOut, baseCurrency, *ammPool)
		if err != nil {
			return err
		}
	}

	// If position is short, CollateralAsset will be in base currency & liabilities should be in trading asset
	if mtp.Position == types.Position_SHORT {
		liabilitiesInCollateralTokenIn := sdk.NewCoin(baseCurrency, liabilities)
		liabilities, _, err = k.EstimateSwap(ctx, liabilitiesInCollateralTokenIn, mtp.LiabilitiesAsset, *ammPool)
		if err != nil {
			return err
		}
	}

	mtp.Collateral = collateralAmount
	mtp.Liabilities = liabilities
	mtp.Custody = custodyAmount

	// calculate mtp take profit custody, delta y_tp_c = delta x_l / take profit price (take profit custody = liabilities / take profit price)
	mtp.TakeProfitCustody = types.CalcMTPTakeProfitCustody(*mtp)

	// calculate mtp take profit liabilities, delta x_tp_l = delta y_tp_c * current price (take profit liabilities = take profit custody * current price)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, mtp, baseCurrency)
	if err != nil {
		return err
	}

	h, err := k.GetMTPHealth(ctx, *mtp, *ammPool, baseCurrency) // set mtp in func or return h?
	if err != nil {
		return err
	}
	mtp.MtpHealth = h

	ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
	if err != nil {
		return err
	}

	collateralCoins := sdk.NewCoins(collateralCoin)
	err = k.bankKeeper.SendCoins(ctx, senderAddress, ammPoolAddr, collateralCoins)
	if err != nil {
		return err
	}
	err = k.amm.AddToPoolBalance(ctx, ammPool, math.ZeroInt(), collateralCoins)
	if err != nil {
		return err
	}

	// All liability has to be in liabilities asset
	err = pool.UpdateLiabilities(mtp.LiabilitiesAsset, mtp.Liabilities, true, mtp.Position)
	if err != nil {
		return err
	}

	// All take profit liability has to be in liabilities asset
	err = pool.UpdateTakeProfitLiabilities(mtp.LiabilitiesAsset, mtp.TakeProfitLiabilities, true, mtp.Position)
	if err != nil {
		return err
	}

	// All take profit custody has to be in custody asset
	err = pool.UpdateTakeProfitCustody(mtp.CustodyAsset, mtp.TakeProfitCustody, true, mtp.Position)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return k.SetMTP(ctx, mtp)
}

func (k Keeper) TakeInCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool) error {
	err := pool.UpdateCustody(mtp.CustodyAsset, mtp.Custody, true, mtp.Position)
	if err != nil {
		return nil
	}

	k.SetPool(ctx, *pool)

	return nil
}

func (k Keeper) BorrowInterestRateComputationByPosition(pool types.Pool, ammPool ammtypes.Pool, position types.Position) (sdk.Dec, error) {
	poolAssets := pool.GetPoolAssets(position)
	targetBorrowInterestRate := sdk.OneDec()
	for _, asset := range *poolAssets {
		ammBalance, err := ammPool.GetAmmPoolBalance(asset.AssetDenom)
		if err != nil {
			return sdk.ZeroDec(), err
		}

		balance := ammBalance.Sub(asset.Custody).ToLegacyDec()
		liabilities := asset.Liabilities.ToLegacyDec()

		// Ensure balance is not zero to avoid division by zero
		if balance.IsZero() {
			return sdk.ZeroDec(), nil
		}
		if balance.Add(liabilities).IsZero() {
			return sdk.ZeroDec(), nil
		}

		mul := balance.Add(liabilities).Quo(balance)
		targetBorrowInterestRate = targetBorrowInterestRate.Mul(mul)
	}
	return targetBorrowInterestRate, nil
}

func (k Keeper) BorrowInterestRateComputation(ctx sdk.Context, pool types.Pool) (sdk.Dec, error) {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return sdk.ZeroDec(), errorsmod.Wrap(types.ErrBalanceNotAvailable, "Balance not available")
	}

	borrowInterestRateMax := k.GetBorrowInterestRateMax(ctx)
	borrowInterestRateMin := k.GetBorrowInterestRateMin(ctx)
	borrowInterestRateIncrease := k.GetBorrowInterestRateIncrease(ctx)
	borrowInterestRateDecrease := k.GetBorrowInterestRateDecrease(ctx)
	healthGainFactor := k.GetHealthGainFactor(ctx)

	prevBorrowInterestRate := pool.BorrowInterestRate

	targetBorrowInterestRate := healthGainFactor
	targetBorrowInterestRateLong, err := k.BorrowInterestRateComputationByPosition(pool, ammPool, types.Position_LONG)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	targetBorrowInterestRateShort, err := k.BorrowInterestRateComputationByPosition(pool, ammPool, types.Position_SHORT)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	targetBorrowInterestRate = targetBorrowInterestRate.Mul(targetBorrowInterestRateLong)
	targetBorrowInterestRate = targetBorrowInterestRate.Mul(targetBorrowInterestRateShort)

	borrowInterestRateChange := targetBorrowInterestRate.Sub(prevBorrowInterestRate)
	borrowInterestRate := prevBorrowInterestRate
	if borrowInterestRateChange.GTE(borrowInterestRateDecrease.Mul(sdk.NewDec(-1))) && borrowInterestRateChange.LTE(borrowInterestRateIncrease) {
		borrowInterestRate = targetBorrowInterestRate
	} else if borrowInterestRateChange.GT(borrowInterestRateIncrease) {
		borrowInterestRate = prevBorrowInterestRate.Add(borrowInterestRateIncrease)
	} else if borrowInterestRateChange.LT(borrowInterestRateDecrease.Mul(sdk.NewDec(-1))) {
		borrowInterestRate = prevBorrowInterestRate.Sub(borrowInterestRateDecrease)
	}

	newBorrowInterestRate := borrowInterestRate

	if borrowInterestRate.GT(borrowInterestRateMin) && borrowInterestRate.LT(borrowInterestRateMax) {
		newBorrowInterestRate = borrowInterestRate
	} else if borrowInterestRate.LTE(borrowInterestRateMin) {
		newBorrowInterestRate = borrowInterestRateMin
	} else if borrowInterestRate.GTE(borrowInterestRateMax) {
		newBorrowInterestRate = borrowInterestRateMax
	}

	return newBorrowInterestRate, nil
}

func (k Keeper) TakeFundPayment(ctx sdk.Context, amount math.Int, returnAsset string, takePercentage sdk.Dec, fundAddr sdk.AccAddress, ammPool *ammtypes.Pool) (math.Int, error) {
	takeAmount := amount.ToLegacyDec().Mul(takePercentage).TruncateInt()

	if !takeAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, takeAmount))

		ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		err = k.bankKeeper.SendCoins(ctx, ammPoolAddr, fundAddr, takeCoins)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		err = k.amm.RemoveFromPoolBalance(ctx, ammPool, math.ZeroInt(), takeCoins)
		if err != nil {
			return math.ZeroInt(), err
		}

	}
	return takeAmount, nil
}

// Set the perpetual hooks.
func (k *Keeper) SetHooks(gh types.PerpetualHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set perpetual hooks twice")
	}

	k.hooks = gh

	return k
}

func (k Keeper) NukeDB(ctx sdk.Context) {
	// delete all mtps
	store := ctx.KVStore(k.storeKey)
	mtpIterator := sdk.KVStorePrefixIterator(store, types.MTPPrefix)
	defer mtpIterator.Close()

	for ; mtpIterator.Valid(); mtpIterator.Next() {
		store.Delete(mtpIterator.Key())
	}

	// delete all pools
	poolIterator := sdk.KVStorePrefixIterator(store, types.PoolKeyPrefix)
	defer poolIterator.Close()
	for ; poolIterator.Valid(); poolIterator.Next() {
		store.Delete(poolIterator.Key())
	}

	k.SetMTPCount(ctx, 0)
	k.SetOpenMTPCount(ctx, 0)

	k.DeleteAllFundingRate(ctx)
	k.DeleteAllInterestRate(ctx)

	store.Delete(types.KeyPrefix(types.ParamsKey))

	return
}
