package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	pkeeper "github.com/elys-network/elys/v6/x/parameter/keeper"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	tierkeeper "github.com/elys-network/elys/v6/x/tier/keeper"
	"github.com/osmosis-labs/osmosis/osmomath"
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

func (k Keeper) Borrow(ctx sdk.Context, collateralAmount math.Int, custodyAmount math.Int, mtp *types.MTP, ammPool *ammtypes.Pool, pool *types.Pool, proxyLeverage math.LegacyDec, baseCurrency string) (types.PerpetualFees, error) {
	senderAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	collateralCoin := sdk.NewCoin(mtp.CollateralAsset, collateralAmount)

	if !k.bankKeeper.HasBalance(ctx, senderAddress, collateralCoin) {
		return types.PerpetualFees{}, types.ErrBalanceNotAvailable
	}

	// eta = leverage - 1
	eta := proxyLeverage.Sub(math.LegacyOneDec())
	liabilitiesInCollateral := eta.MulInt(collateralAmount).TruncateInt()
	liabilities := liabilitiesInCollateral
	// If collateral asset is not base currency, should calculate liability in base currency with the given out.
	// For LONG, Liability has to be in base currency, CollateralAsset can be trading asset or base currency
	// For SHORT, Liability has to be in trading asset and CollateralAsset will be in base currency, so this if case only applies to LONG
	slippageAmount, weightBreakingFee, liabilitiesOracleAmount := osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec()
	perpetualFees, takerFees := math.LegacyZeroDec(), math.LegacyZeroDec()

	totalPerpFees := types.NewPerpetualFeesWithEmptyCoins()

	if mtp.CollateralAsset != baseCurrency {
		if !liabilities.IsZero() {
			liabilitiesInCollateralTokenOut := sdk.NewCoin(mtp.CollateralAsset, liabilitiesInCollateral)
			// Calculate base currency amount given atom out amount and we use it liabilty amount in base currency
			liabilities, _, slippageAmount, weightBreakingFee, liabilitiesOracleAmount, perpetualFees, takerFees, err = k.EstimateSwapGivenOut(ctx, liabilitiesInCollateralTokenOut, baseCurrency, *ammPool, mtp.Address)
			if err != nil {
				return types.PerpetualFees{}, err
			}
			totalPerpFees = k.CalculatePerpetualFees(ctx, ammPool.PoolParams.UseOracle, sdk.NewCoin(baseCurrency, liabilities), liabilitiesInCollateralTokenOut, slippageAmount, weightBreakingFee, perpetualFees, takerFees, liabilitiesOracleAmount, false, false)
		}
	}

	// If position is short, CollateralAsset will be in base currency & liabilities should be in trading asset
	if mtp.Position == types.Position_SHORT {
		// liabilities.IsZero() happens when we are consolidating with leverage 1 as eta = 0
		if !liabilities.IsZero() {
			liabilitiesInCollateralTokenIn := sdk.NewCoin(baseCurrency, liabilities)
			liabilities, _, slippageAmount, weightBreakingFee, liabilitiesOracleAmount, perpetualFees, takerFees, err = k.EstimateSwapGivenOut(ctx, liabilitiesInCollateralTokenIn, mtp.LiabilitiesAsset, *ammPool, mtp.Address)
			if err != nil {
				return types.PerpetualFees{}, err
			}
			perpFees := k.CalculatePerpetualFees(ctx, ammPool.PoolParams.UseOracle, sdk.NewCoin(mtp.LiabilitiesAsset, liabilities), liabilitiesInCollateralTokenIn, slippageAmount, weightBreakingFee, perpetualFees, takerFees, liabilitiesOracleAmount, false, false)
			totalPerpFees = totalPerpFees.Add(perpFees)
		}
	}

	// Track slippage and weight breaking fee slippage in amm via perpetual
	for _, coin := range totalPerpFees.SlippageFees {
		k.amm.TrackSlippage(ctx, ammPool.PoolId, coin)
		k.amm.TrackWeightBreakingSlippage(ctx, ammPool.PoolId, coin)
	}
	for _, coin := range totalPerpFees.WeightBreakingFees {
		if coin.Amount.IsPositive() {
			k.amm.TrackWeightBreakingSlippage(ctx, ammPool.PoolId, coin)
		}
	}

	mtp.Collateral = collateralAmount
	mtp.Liabilities = liabilities
	mtp.Custody = custodyAmount

	mtp.MtpHealth, err = k.GetMTPHealth(ctx, *mtp)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	collateralCoins := sdk.NewCoins(collateralCoin)
	err = k.SendToAmmPool(ctx, senderAddress, ammPool, collateralCoins)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	// send fees to masterchef and taker collection address
	var totalAmount math.Int
	if mtp.Position == types.Position_LONG && mtp.CollateralAsset == baseCurrency {
		totalAmount = proxyLeverage.MulInt(collateralAmount).TruncateInt()
	} else {
		totalAmount = liabilitiesInCollateral
	}
	_, err = k.SendFeesToMasterchefAndTakerCollection(ctx, senderAddress, mtp.Address, totalAmount, mtp.CollateralAsset, ammPool, &totalPerpFees)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	err = pool.UpdateCustody(mtp.CustodyAsset, mtp.Custody, true, mtp.Position)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	// All liability has to be in liabilities asset
	err = pool.UpdateLiabilities(mtp.LiabilitiesAsset, mtp.Liabilities, true, mtp.Position)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	err = pool.UpdateCollateral(mtp.CollateralAsset, mtp.Collateral, true, mtp.Position)
	if err != nil {
		return types.PerpetualFees{}, err
	}

	k.SetPool(ctx, *pool)

	return totalPerpFees, k.SetMTP(ctx, mtp)
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

func (k Keeper) BorrowInterestRateComputationByPosition(pool types.Pool, ammPool ammtypes.Pool, position types.Position) (math.LegacyDec, error) {
	poolAssets := pool.GetPoolAssets(position)
	targetBorrowInterestRate := math.LegacyOneDec()
	for _, asset := range *poolAssets {
		ammBalance, err := ammPool.GetAmmPoolBalance(asset.AssetDenom)
		if err != nil {
			return math.LegacyZeroDec(), err
		}

		balance := ammBalance.Sub(asset.Custody)
		liabilities := asset.Liabilities

		// Ensure balance is not zero to avoid division by zero
		if balance.IsZero() {
			return math.LegacyZeroDec(), nil
		}
		if balance.Add(liabilities).IsZero() {
			return math.LegacyZeroDec(), nil
		}

		mul := balance.Add(liabilities).ToLegacyDec().Quo(balance.ToLegacyDec())
		targetBorrowInterestRate = targetBorrowInterestRate.Mul(mul)
	}
	return targetBorrowInterestRate, nil
}

func (k Keeper) BorrowInterestRateComputation(ctx sdk.Context, pool types.Pool) (math.LegacyDec, error) {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return math.LegacyZeroDec(), errors.New("amm pool not found")
	}

	params := k.GetParams(ctx)
	borrowInterestRateMax := params.BorrowInterestRateMax
	borrowInterestRateMin := params.BorrowInterestRateMin
	borrowInterestRateIncrease := params.BorrowInterestRateIncrease
	borrowInterestRateDecrease := params.BorrowInterestRateDecrease
	healthGainFactor := params.HealthGainFactor

	prevBorrowInterestRate := pool.BorrowInterestRate

	targetBorrowInterestRate := healthGainFactor
	targetBorrowInterestRateLong, err := k.BorrowInterestRateComputationByPosition(pool, ammPool, types.Position_LONG)
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	targetBorrowInterestRateShort, err := k.BorrowInterestRateComputationByPosition(pool, ammPool, types.Position_SHORT)
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

func (k Keeper) CollectInsuranceFund(ctx sdk.Context, amount math.Int, returnAsset string, ammPool *ammtypes.Pool, pool types.Pool) (math.Int, error) {
	params := k.GetParams(ctx)
	insuranceAmount := params.BorrowInterestPaymentFundPercentage.MulInt(amount).TruncateInt()

	if !insuranceAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, insuranceAmount))

		err := k.SendFromAmmPool(ctx, ammPool, pool.GetInsuranceAccount(), takeCoins)
		if err != nil {
			return math.ZeroInt(), err
		}

	}
	return insuranceAmount, nil
}

func (k Keeper) SendFeesToMasterchefAndTakerCollection(ctx sdk.Context, senderAddress sdk.AccAddress, tierAddressStr string, liabilitiesInCollateral math.Int, collateralDenom string, ammPool *ammtypes.Pool, perpFees *types.PerpetualFees) (math.Int, error) {
	tierAddress, err := sdk.AccAddressFromBech32(tierAddressStr)
	if err != nil {
		return math.ZeroInt(), err
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, tierAddress)
	params := k.GetParams(ctx)
	perpetualFee := ammtypes.ApplyDiscount(params.GetBigDecPerpetualSwapFee(), tier.GetBigDecDiscount())
	takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()
	sendToMasterchef := perpetualFee.Dec().Mul(math.LegacyNewDecFromInt(liabilitiesInCollateral)).TruncateInt()
	sendToTakerCollection := takersFee.Dec().Mul(math.LegacyNewDecFromInt(liabilitiesInCollateral)).TruncateInt()

	if sendToMasterchef.IsPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(ammPool.GetRebalanceTreasury())
		sendToMasterchefCoin := sdk.NewCoin(collateralDenom, sendToMasterchef)
		perpFees.PerpFees = perpFees.PerpFees.Add(sendToMasterchefCoin)
		err := k.bankKeeper.SendCoins(ctx, senderAddress, rebalanceTreasury, sdk.NewCoins(sendToMasterchefCoin))
		if err != nil {
			return math.ZeroInt(), err
		}

		err = k.amm.OnCollectFee(ctx, *ammPool, sdk.NewCoins(sendToMasterchefCoin))
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	if sendToTakerCollection.IsPositive() {
		takerAddress, err := sdk.AccAddressFromBech32(k.parameterKeeper.GetParams(ctx).TakerFeeCollectionAddress)
		if err != nil {
			return math.ZeroInt(), err
		}
		sendToTakerCollectionCoin := sdk.NewCoin(collateralDenom, sendToTakerCollection)
		perpFees.TakerFees = perpFees.TakerFees.Add(sendToTakerCollectionCoin)

		err = k.bankKeeper.SendCoins(ctx, senderAddress, takerAddress, sdk.NewCoins(sendToTakerCollectionCoin))
		if err != nil {
			return math.ZeroInt(), err
		}
	}
	return sendToMasterchef.Add(sendToTakerCollection), nil
}

// Set the perpetual hooks.
func (k *Keeper) SetHooks(gh types.PerpetualHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set perpetual hooks twice")
	}

	k.hooks = gh

	return k
}

func (k Keeper) GetPerpFeesInUSD(ctx sdk.Context, perpFeesCoins types.PerpetualFees) (perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd math.LegacyDec) {

	perpFeesInUsd = k.amm.CalculateCoinsUSDValue(ctx, perpFeesCoins.PerpFees).Dec()
	slippageFeesInUsd = k.amm.CalculateCoinsUSDValue(ctx, perpFeesCoins.SlippageFees).Dec()
	weightBreakingFeesInUsd = k.amm.CalculateCoinsUSDValue(ctx, perpFeesCoins.WeightBreakingFees).Dec()
	takerFeesInUsd = k.amm.CalculateCoinsUSDValue(ctx, perpFeesCoins.TakerFees).Dec()

	return perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd
}
