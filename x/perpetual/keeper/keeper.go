package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	pkeeper "github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	tierkeeper "github.com/elys-network/elys/x/tier/keeper"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
		authority          string
		parameterKeeper    *pkeeper.Keeper
		amm                types.AmmKeeper
		bankKeeper         types.BankKeeper
		oracleKeeper       types.OracleKeeper
		assetProfileKeeper types.AssetProfileKeeper
		tierKeeper         *tierkeeper.Keeper

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
	tierKeeper *tierkeeper.Keeper,
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
		tierKeeper:         tierKeeper,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetTierKeeper(tk *tierkeeper.Keeper) {
	k.tierKeeper = tk
}

func (k *Keeper) GetTierKeeper() *tierkeeper.Keeper {
	return k.tierKeeper
}

func (k Keeper) Borrow(ctx sdk.Context, collateralAmount math.Int, custodyAmount math.Int, mtp *types.MTP, ammPool *ammtypes.Pool, pool *types.Pool, proxyLeverage osmomath.BigDec, baseCurrency string) error {
	senderAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return err
	}

	collateralCoin := sdk.NewCoin(mtp.CollateralAsset, collateralAmount)

	if !k.bankKeeper.HasBalance(ctx, senderAddress, collateralCoin) {
		return types.ErrBalanceNotAvailable
	}

	// eta = leverage - 1
	eta := proxyLeverage.Sub(osmomath.OneBigDec())
	liabilitiesInCollateral := osmomath.BigDecFromSDKInt(collateralAmount).Mul(eta).Dec().TruncateInt()
	liabilities := liabilitiesInCollateral
	// If collateral asset is not base currency, should calculate liability in base currency with the given out.
	// For LONG, Liability has to be in base currency, CollateralAsset can be trading asset or base currency
	// For SHORT, Liability has to be in trading asset and CollateralAsset will be in base currency, so this if case only applies to LONG
	if mtp.CollateralAsset != baseCurrency {
		if !liabilities.IsZero() {
			liabilitiesInCollateralTokenOut := sdk.NewCoin(mtp.CollateralAsset, liabilitiesInCollateral)
			// Calculate base currency amount given atom out amount and we use it liabilty amount in base currency
			liabilities, _, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesInCollateralTokenOut, baseCurrency, *ammPool, mtp.Address)
			if err != nil {
				return err
			}
		}
	}

	// If position is short, CollateralAsset will be in base currency & liabilities should be in trading asset
	if mtp.Position == types.Position_SHORT {
		// liabilities.IsZero() happens when we are consolidating with leverage 1 as eta = 0
		if !liabilities.IsZero() {
			liabilitiesInCollateralTokenIn := sdk.NewCoin(baseCurrency, liabilities)
			liabilities, _, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesInCollateralTokenIn, mtp.LiabilitiesAsset, *ammPool, mtp.Address)
			if err != nil {
				return err
			}
		}
	}

	mtp.Collateral = collateralAmount
	mtp.Liabilities = liabilities
	mtp.Custody = custodyAmount

	// calculate mtp take profit custody, delta y_tp_c = delta x_l / take profit price (take profit custody = liabilities / take profit price: LONG, profit custody = liabilities * take profit price: SHORT)
	mtp.TakeProfitCustody, err = k.CalcMTPTakeProfitCustody(ctx, *mtp)
	if err != nil {
		return err
	}

	// calculate mtp take profit liabilities, delta x_tp_l = delta y_tp_c * current price (take profit liabilities = take profit custody * current price LONG, take profit custody / current price SHORT)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, *mtp)
	if err != nil {
		return err
	}

	h, err := k.GetMTPHealth(ctx, *mtp, *ammPool, baseCurrency) // set mtp in func or return h?
	if err != nil {
		return err
	}
	mtp.MtpHealth = h.Dec()

	collateralCoins := sdk.NewCoins(collateralCoin)
	err = k.SendToAmmPool(ctx, senderAddress, ammPool, collateralCoins)
	if err != nil {
		return err
	}

	err = pool.UpdateCustody(mtp.CustodyAsset, mtp.Custody, true, mtp.Position)
	if err != nil {
		return nil
	}

	// All liability has to be in liabilities asset
	err = pool.UpdateLiabilities(mtp.LiabilitiesAsset, mtp.Liabilities, true, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateCollateral(mtp.CollateralAsset, mtp.Collateral, true, mtp.Position)
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

func (k Keeper) SendToAmmPool(ctx sdk.Context, senderAddress sdk.AccAddress, ammPool *ammtypes.Pool, coins sdk.Coins) error {
	ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoins(ctx, senderAddress, ammPoolAddr, coins)
	if err != nil {
		return err
	}
	err = k.amm.AddToPoolBalanceAndUpdateLiquidity(ctx, ammPool, math.ZeroInt(), coins)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) SendFromAmmPool(ctx sdk.Context, ammPool *ammtypes.Pool, receiverAddress sdk.AccAddress, coins sdk.Coins) error {
	ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoins(ctx, ammPoolAddr, receiverAddress, coins)
	if err != nil {
		return err
	}
	err = k.amm.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, ammPool, math.ZeroInt(), coins)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) BorrowInterestRateComputationByPosition(pool types.Pool, ammPool ammtypes.Pool, position types.Position) (osmomath.BigDec, error) {
	poolAssets := pool.GetPoolAssets(position)
	targetBorrowInterestRate := osmomath.OneBigDec()
	for _, asset := range *poolAssets {
		ammBalance, err := ammPool.GetAmmPoolBalance(asset.AssetDenom)
		if err != nil {
			return osmomath.ZeroBigDec(), err
		}

		balance := osmomath.BigDecFromSDKInt(ammBalance.Sub(asset.Custody))
		liabilities := asset.GetBigDecLiabilities()

		// Ensure balance is not zero to avoid division by zero
		if balance.IsZero() {
			return osmomath.ZeroBigDec(), nil
		}
		if balance.Add(liabilities).IsZero() {
			return osmomath.ZeroBigDec(), nil
		}

		mul := balance.Add(liabilities).Quo(balance)
		targetBorrowInterestRate = targetBorrowInterestRate.Mul(mul)
	}
	return targetBorrowInterestRate, nil
}

func (k Keeper) BorrowInterestRateComputation(ctx sdk.Context, pool types.Pool) (osmomath.BigDec, error) {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return osmomath.ZeroBigDec(), errorsmod.Wrap(types.ErrBalanceNotAvailable, "Balance not available")
	}

	borrowInterestRateMax := k.GetBigDecBorrowInterestRateMax(ctx)
	borrowInterestRateMin := k.GetBigDecBorrowInterestRateMin(ctx)
	borrowInterestRateIncrease := k.GetBigDecBorrowInterestRateIncrease(ctx)
	borrowInterestRateDecrease := k.GetBigDecBorrowInterestRateDecrease(ctx)
	healthGainFactor := k.GetBigDecHealthGainFactor(ctx)

	prevBorrowInterestRate := pool.GetBigDecBorrowInterestRate()

	targetBorrowInterestRate := healthGainFactor
	targetBorrowInterestRateLong, err := k.BorrowInterestRateComputationByPosition(pool, ammPool, types.Position_LONG)
	if err != nil {
		return osmomath.ZeroBigDec(), err
	}
	targetBorrowInterestRateShort, err := k.BorrowInterestRateComputationByPosition(pool, ammPool, types.Position_SHORT)
	if err != nil {
		return osmomath.ZeroBigDec(), err
	}
	targetBorrowInterestRate = targetBorrowInterestRate.Mul(targetBorrowInterestRateLong)
	targetBorrowInterestRate = targetBorrowInterestRate.Mul(targetBorrowInterestRateShort)

	borrowInterestRateChange := targetBorrowInterestRate.Sub(prevBorrowInterestRate)
	borrowInterestRate := prevBorrowInterestRate
	if borrowInterestRateChange.GTE(borrowInterestRateDecrease.Mul(osmomath.NewBigDec(-1))) && borrowInterestRateChange.LTE(borrowInterestRateIncrease) {
		borrowInterestRate = targetBorrowInterestRate
	} else if borrowInterestRateChange.GT(borrowInterestRateIncrease) {
		borrowInterestRate = prevBorrowInterestRate.Add(borrowInterestRateIncrease)
	} else if borrowInterestRateChange.LT(borrowInterestRateDecrease.Mul(osmomath.NewBigDec(-1))) {
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

func (k Keeper) CollectInsuranceFund(ctx sdk.Context, amount math.Int, returnAsset string, ammPool *ammtypes.Pool, pool types.Pool) (math.Int, error) {
	params := k.GetParams(ctx)
	insuranceAmount := osmomath.BigDecFromSDKInt(amount).Mul(params.GetBigDecBorrowInterestPaymentFundPercentage()).Dec().TruncateInt()

	if !insuranceAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, insuranceAmount))

		err := k.SendFromAmmPool(ctx, ammPool, pool.GetInsuranceAccount(), takeCoins)
		if err != nil {
			return math.ZeroInt(), err
		}

	}
	return insuranceAmount, nil
}

// Set the perpetual hooks.
func (k *Keeper) SetHooks(gh types.PerpetualHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set perpetual hooks twice")
	}

	k.hooks = gh

	return k
}
