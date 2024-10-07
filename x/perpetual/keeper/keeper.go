package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
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
		types.PoolChecker
		types.OpenChecker
		types.OpenDefineAssetsChecker
		types.ClosePositionChecker
		types.CloseEstimationChecker

		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
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
	storeService store.KVStoreService,
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
		storeService:       storeService,
		authority:          authority,
		amm:                amm,
		bankKeeper:         bk,
		oracleKeeper:       oracleKeeper,
		assetProfileKeeper: assetProfileKeeper,
		parameterKeeper:    parameterKeeper,
	}

	keeper.AuthorizationChecker = keeper
	keeper.PositionChecker = keeper
	keeper.PoolChecker = keeper
	keeper.OpenChecker = keeper
	keeper.OpenDefineAssetsChecker = keeper
	keeper.ClosePositionChecker = keeper
	keeper.CloseEstimationChecker = keeper

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (math.Int, error) {
	perpetualEnabled := k.IsPoolEnabled(ctx, ammPool.PoolId)
	if !perpetualEnabled {
		return math.ZeroInt(), errorsmod.Wrap(types.ErrPerpetualDisabled, "Perpetual disabled pool")
	}

	tokensOut := sdk.Coins{tokenOutAmount}
	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	swapResult, _, err := k.amm.CalcInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, math.LegacyZeroDec())
	if err != nil {
		return math.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return math.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}

func (k Keeper) Borrow(ctx sdk.Context, collateralAmount math.Int, custodyAmount math.Int, mtp *types.MTP, ammPool *ammtypes.Pool, pool *types.Pool, eta math.LegacyDec, baseCurrency string, isBroker bool) error {
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

	collateralAmountDec := math.LegacyNewDecFromBigInt(collateralAmount.BigInt())
	liabilitiesDec := collateralAmountDec.Mul(eta)

	// If collateral asset is not base currency, should calculate liability in base currency with the given out.
	// Liability has to be in base currency
	if mtp.CollateralAsset != baseCurrency {
		// ATOM amount
		etaAmt := liabilitiesDec.TruncateInt()
		etaAmtToken := sdk.NewCoin(mtp.CollateralAsset, etaAmt)
		// Calculate base currency amount given atom out amount and we use it liabilty amount in base currency
		liabilityAmt, err := k.OpenDefineAssetsChecker.EstimateSwapGivenOut(ctx, etaAmtToken, baseCurrency, *ammPool)
		if err != nil {
			return err
		}

		liabilitiesDec = math.LegacyNewDecFromInt(liabilityAmt)
	}

	// If position is short, liabilities should be swapped to liabilities asset
	if mtp.Position == types.Position_SHORT {
		liabilitiesAmtTokenIn := sdk.NewCoin(baseCurrency, liabilitiesDec.TruncateInt())
		liabilitiesAmt, err := k.OpenDefineAssetsChecker.EstimateSwap(ctx, liabilitiesAmtTokenIn, mtp.LiabilitiesAsset, *ammPool)
		if err != nil {
			return err
		}

		liabilitiesDec = math.LegacyNewDecFromInt(liabilitiesAmt)
	}

	mtp.Collateral = collateralAmount
	mtp.Liabilities = math.NewIntFromBigInt(liabilitiesDec.TruncateInt().BigInt())
	mtp.Custody = custodyAmount

	// calculate mtp take profit custody, delta y_tp_c = delta x_l / take profit price (take profit custody = liabilities / take profit price)
	mtp.TakeProfitCustody = types.CalcMTPTakeProfitCustody(mtp)

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

	err = pool.UpdateBalance(ctx, mtp.CollateralAsset, collateralAmount, true, mtp.Position)
	if err != nil {
		return err
	}

	// All liability has to be in liabilities asset
	err = pool.UpdateLiabilities(ctx, mtp.LiabilitiesAsset, mtp.Liabilities, true, mtp.Position)
	if err != nil {
		return err
	}

	// All take profit liability has to be in liabilities asset
	err = pool.UpdateTakeProfitLiabilities(ctx, mtp.LiabilitiesAsset, mtp.TakeProfitLiabilities, true, mtp.Position)
	if err != nil {
		return err
	}

	// All take profit custody has to be in custody asset
	err = pool.UpdateTakeProfitCustody(ctx, mtp.CustodyAsset, mtp.TakeProfitCustody, true, mtp.Position)
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

func (k Keeper) CalculatePoolHealthByPosition(ctx sdk.Context, pool *types.Pool, ammPool ammtypes.Pool, position types.Position) math.LegacyDec {
	poolAssets := pool.GetPoolAssets(position)
	H := math.LegacyNewDec(1)
	for _, asset := range *poolAssets {
		ammBalance, err := types.GetAmmPoolBalance(ammPool, asset.AssetDenom)
		if err != nil {
			return math.LegacyZeroDec()
		}

		balance := math.LegacyNewDecFromInt(asset.AssetBalance.Add(ammBalance))

		// X_L = X_P_L - X_TP_L (pool liabilities = pool synthetic liabilities - pool take profit liabilities)
		liabilities := math.LegacyNewDecFromInt(asset.Liabilities.Sub(asset.TakeProfitLiabilities))

		if balance.Add(liabilities).IsZero() {
			return math.LegacyZeroDec()
		}

		mul := balance.Quo(balance.Add(liabilities))
		H = H.Mul(mul)
	}
	return H
}

func (k Keeper) CalculatePoolHealth(ctx sdk.Context, pool *types.Pool) math.LegacyDec {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return math.LegacyZeroDec()
	}

	H := k.CalculatePoolHealthByPosition(ctx, pool, ammPool, types.Position_LONG)
	H = H.Mul(k.CalculatePoolHealthByPosition(ctx, pool, ammPool, types.Position_SHORT))

	return H
}

func (k Keeper) TakeInCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool) error {
	err := pool.UpdateBalance(ctx, mtp.CustodyAsset, mtp.Custody, false, mtp.Position)
	if err != nil {
		return nil
	}
	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, mtp.Custody, true, mtp.Position)
	if err != nil {
		return nil
	}

	k.SetPool(ctx, *pool)

	return nil
}

func (k Keeper) IncrementalBorrowInterestPayment(ctx sdk.Context, borrowInterestPayment math.Int, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) (math.Int, error) {
	// if mtp has unpaid borrow interest, add to payment
	// convert it into base currency
	if mtp.BorrowInterestUnpaidCollateral.IsPositive() {
		if mtp.CollateralAsset == baseCurrency {
			borrowInterestPayment = borrowInterestPayment.Add(mtp.BorrowInterestUnpaidCollateral)
		} else {
			unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
			C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
			if err != nil {
				return math.ZeroInt(), err
			}

			borrowInterestPayment = borrowInterestPayment.Add(C)
		}
	}

	borrowInterestPaymentTokenIn := sdk.NewCoin(baseCurrency, borrowInterestPayment)
	// swap borrow interest payment to custody asset for payment
	borrowInterestPaymentCustody, err := k.EstimateSwap(ctx, borrowInterestPaymentTokenIn, mtp.CustodyAsset, ammPool)
	if err != nil {
		return math.ZeroInt(), err
	}

	// If collateralAsset is not in base currency, convert it to original asset format
	if mtp.CollateralAsset != baseCurrency {
		// swap custody amount to collateral for updating borrow interest unpaid
		amtTokenIn := sdk.NewCoin(baseCurrency, borrowInterestPayment)
		borrowInterestPayment, err = k.EstimateSwap(ctx, amtTokenIn, mtp.CollateralAsset, ammPool) // may need spot price here to not deduct fee
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	// if paying unpaid borrow interest reset to 0
	mtp.BorrowInterestUnpaidCollateral = math.ZeroInt()

	// edge case, not enough custody to cover payment
	if borrowInterestPaymentCustody.GT(mtp.Custody) {
		// swap custody amount to collateral for updating borrow interest unpaid
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.Custody)
		custodyAmountCollateral, err := k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.CollateralAsset, ammPool) // may need spot price here to not deduct fee
		if err != nil {
			return math.ZeroInt(), err
		}
		mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(borrowInterestPayment).Sub(custodyAmountCollateral)

		borrowInterestPayment = custodyAmountCollateral
		borrowInterestPaymentCustody = mtp.Custody
	}

	// add payment to total paid - collateral
	mtp.BorrowInterestPaidCollateral = mtp.BorrowInterestPaidCollateral.Add(borrowInterestPayment)

	// add payment to total paid - custody
	mtp.BorrowInterestPaidCustody = mtp.BorrowInterestPaidCustody.Add(borrowInterestPaymentCustody)

	// deduct borrow interest payment from custody amount
	mtp.Custody = mtp.Custody.Sub(borrowInterestPaymentCustody)

	takePercentage := k.GetIncrementalBorrowInterestPaymentFundPercentage(ctx)
	fundAddr := k.GetIncrementalBorrowInterestPaymentFundAddress(ctx)
	takeAmount, err := k.TakeFundPayment(ctx, borrowInterestPaymentCustody, mtp.CustodyAsset, takePercentage, fundAddr, &ammPool)
	if err != nil {
		return math.ZeroInt(), err
	}
	actualBorrowInterestPaymentCustody := borrowInterestPaymentCustody.Sub(takeAmount)

	if !takeAmount.IsZero() {
		k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CustodyAsset, types.EventIncrementalPayFund)
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, borrowInterestPaymentCustody, false, mtp.Position)
	if err != nil {
		return math.ZeroInt(), err
	}

	err = pool.UpdateBalance(ctx, mtp.CustodyAsset, actualBorrowInterestPaymentCustody, true, mtp.Position)
	if err != nil {
		return math.ZeroInt(), err
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return math.ZeroInt(), err
	}

	k.SetPool(ctx, *pool)

	return actualBorrowInterestPaymentCustody, nil
}

func (k Keeper) BorrowInterestRateComputationByPosition(ctx sdk.Context, pool types.Pool, ammPool ammtypes.Pool, position types.Position) (math.LegacyDec, error) {
	poolAssets := pool.GetPoolAssets(position)
	targetBorrowInterestRate := math.LegacyOneDec()
	for _, asset := range *poolAssets {
		ammBalance, err := types.GetAmmPoolBalance(ammPool, asset.AssetDenom)
		if err != nil {
			return math.LegacyZeroDec(), err
		}

		balance := math.LegacyNewDecFromInt(asset.AssetBalance.Add(ammBalance))
		liabilities := math.LegacyNewDecFromInt(asset.Liabilities)

		// Ensure balance is not zero to avoid division by zero
		if balance.IsZero() {
			return math.LegacyZeroDec(), nil
		}
		if balance.Add(liabilities).IsZero() {
			return math.LegacyZeroDec(), nil
		}

		mul := balance.Add(liabilities).Quo(balance)
		targetBorrowInterestRate = targetBorrowInterestRate.Mul(mul)
	}
	return targetBorrowInterestRate, nil
}

func (k Keeper) BorrowInterestRateComputation(ctx sdk.Context, pool types.Pool) (math.LegacyDec, error) {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return math.LegacyZeroDec(), errorsmod.Wrap(types.ErrBalanceNotAvailable, "Balance not available")
	}

	borrowInterestRateMax := k.GetBorrowInterestRateMax(ctx)
	borrowInterestRateMin := k.GetBorrowInterestRateMin(ctx)
	borrowInterestRateIncrease := k.GetBorrowInterestRateIncrease(ctx)
	borrowInterestRateDecrease := k.GetBorrowInterestRateDecrease(ctx)
	healthGainFactor := k.GetHealthGainFactor(ctx)

	prevBorrowInterestRate := pool.BorrowInterestRate

	targetBorrowInterestRate := healthGainFactor
	targetBorrowInterestRateLong, err := k.BorrowInterestRateComputationByPosition(ctx, pool, ammPool, types.Position_LONG)
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	targetBorrowInterestRateShort, err := k.BorrowInterestRateComputationByPosition(ctx, pool, ammPool, types.Position_SHORT)
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	targetBorrowInterestRate = targetBorrowInterestRate.Mul(targetBorrowInterestRateLong)
	targetBorrowInterestRate = targetBorrowInterestRate.Mul(targetBorrowInterestRateShort)

	borrowInterestRateChange := targetBorrowInterestRate.Sub(prevBorrowInterestRate)
	borrowInterestRate := prevBorrowInterestRate
	if borrowInterestRateChange.GTE(borrowInterestRateDecrease.Mul(math.LegacyNewDec(-1))) && borrowInterestRateChange.LTE(borrowInterestRateIncrease) {
		borrowInterestRate = targetBorrowInterestRate
	} else if borrowInterestRateChange.GT(borrowInterestRateIncrease) {
		borrowInterestRate = prevBorrowInterestRate.Add(borrowInterestRateIncrease)
	} else if borrowInterestRateChange.LT(borrowInterestRateDecrease.Mul(math.LegacyNewDec(-1))) {
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

func (k Keeper) TakeFundPayment(ctx sdk.Context, returnAmount math.Int, returnAsset string, takePercentage math.LegacyDec, fundAddr sdk.AccAddress, ammPool *ammtypes.Pool) (math.Int, error) {
	returnAmountDec := math.LegacyNewDecFromBigInt(returnAmount.BigInt())
	takeAmount := math.NewIntFromBigInt(takePercentage.Mul(returnAmountDec).TruncateInt().BigInt())

	if !takeAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, math.NewIntFromBigInt(takeAmount.BigInt())))

		ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
		if err != nil {
			return math.ZeroInt(), err
		}
		err = k.bankKeeper.SendCoins(ctx, ammPoolAddr, fundAddr, takeCoins)
		if err != nil {
			return math.ZeroInt(), err
		}
	}
	return takeAmount, nil
}

// CalcTakeFundPayment calculates the take fund payment
func (k Keeper) CalcTakeFundPayment(ctx sdk.Context, returnAmount math.Int, returnAsset string, takePercentage math.LegacyDec) math.Int {
	returnAmountDec := math.LegacyNewDecFromBigInt(returnAmount.BigInt())
	takeAmount := math.NewIntFromBigInt(takePercentage.Mul(returnAmountDec).TruncateInt().BigInt())

	return takeAmount
}

// Set the perpetual hooks.
func (k *Keeper) SetHooks(gh types.PerpetualHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set perpetual hooks twice")
	}

	k.hooks = gh

	return k
}
