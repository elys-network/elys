package keeper

import (
	sdkmath "cosmossdk.io/math"
	"math"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/v7/x/commitment/types"
	estakingtypes "github.com/elys-network/elys/v7/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/v7/x/masterchef/types"
	perpetualtypes "github.com/elys-network/elys/v7/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/v7/x/stablestake/types"
	tradeshieldtypes "github.com/elys-network/elys/v7/x/tradeshield/types"

	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/tier/types"
)

func (k Keeper) RetrieveAllPortfolio(ctx sdk.Context, user sdk.AccAddress) {
	// set today + user -> amount
	todayDate := k.GetDateFromContext(ctx)

	_, found := k.GetPortfolio(ctx, user, todayDate)
	if found {
		return
	}

	totalValue := osmomath.ZeroBigDec()

	// Liquid assets
	liq := k.RetrieveLiquidAssetsTotal(ctx, user)
	totalValue = totalValue.Add(liq)

	// Rewards
	rew := k.RetrieveRewardsTotal(ctx, user)
	totalValue = totalValue.Add(rew)

	// Perpetual
	_, _, perp := k.RetrievePerpetualTotal(ctx, user)
	totalValue = totalValue.Add(perp)

	// Pool assets
	staked := k.RetrievePoolTotal(ctx, user)
	totalValue = totalValue.Add(staked)

	// Staked assets
	commit, delegations, unbondings, totalVesting := k.RetrieveStaked(ctx, user)
	// convert vesting to usd
	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if found {
		edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
		totalVesting = totalVesting.Mul(edenDenomPrice)
	}
	totalValue = totalValue.Add(commit).Add(delegations).Add(unbondings).Add(totalVesting)

	// LeverageLp
	_, _, lev := k.RetrieveLeverageLpTotal(ctx, user)
	totalValue = totalValue.Add(lev)

	// Tradeshield assets
	tradeshieldTotal := k.RetrieveTradeshieldTotal(ctx, user)
	totalValue = totalValue.Add(tradeshieldTotal)

	k.SetPortfolio(ctx, types.NewPortfolioWithContextDate(todayDate, user, totalValue.Dec()))
}

func (k Keeper) RetrievePoolTotal(ctx sdk.Context, user sdk.AccAddress) osmomath.BigDec {
	totalValue := osmomath.ZeroBigDec()
	commitments := k.commitment.GetCommitments(ctx, user)
	for _, commitment := range commitments.CommittedTokens {
		// Pool balance
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				continue
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				continue
			}
			info := k.amm.PoolExtraInfo(ctx, pool, types.OneDay)
			amount := osmomath.BigDecFromSDKInt(commitment.Amount)
			totalValue = totalValue.Add(amount.Mul(info.GetBigDecLpTokenPrice()).Quo(ammtypes.OneShareBigDec))
		}
	}

	return totalValue
}

func (k Keeper) RetrieveStaked(ctx sdk.Context, user sdk.AccAddress) (osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec) {
	totalCommit := osmomath.ZeroBigDec()
	commitments := k.commitment.GetCommitments(ctx, user)
	totalVested := osmomath.ZeroBigDec()
	vestingResp, vestErr := k.commitment.CommitmentVestingInfo(ctx, &commitmenttypes.QueryCommitmentVestingInfoRequest{Address: user.String()})
	if vestErr == nil {
		totalVested = osmomath.BigDecFromSDKInt(vestingResp.Total)
	}
	for _, commitment := range commitments.CommittedTokens {
		if !strings.HasPrefix(commitment.Denom, "amm/pool") {
			if strings.HasPrefix(commitment.Denom, "stablestake/share") {
				poolId, err := stablestaketypes.GetPoolIDFromPath(commitment.Denom)
				if err != nil {
					continue
				}
				pool, found := k.stablestakeKeeper.GetPool(ctx, poolId)
				if !found {
					continue
				}
				redemptionRate := k.stablestakeKeeper.CalculateRedemptionRateForPool(ctx, pool)
				tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, pool.GetDepositDenom())
				usdValue := commitment.GetBigDecAmount().Mul(redemptionRate).Mul(tokenPrice)
				totalCommit = totalCommit.Add(usdValue)
				continue
			}
			if commitment.Denom == ptypes.Eden {
				commitment.Denom = ptypes.Elys
			}
			tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, commitment.Denom)
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, commitment.Denom)
			if !found {
				continue
			}
			if tokenPrice.Equal(osmomath.ZeroBigDec()) {
				tokenPrice = k.amm.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := commitment.GetBigDecAmount()
			totalCommit = totalCommit.Add(amount.Mul(tokenPrice))
		}
	}

	// Delegations
	totalDelegations := osmomath.ZeroBigDec()
	delegations, err := k.stakingKeeper.GetAllDelegatorDelegations(ctx, user)
	if err != nil {
		ctx.Logger().Error(err.Error())
		delegations = []stakingtypes.Delegation{}
	}
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		panic(err)
	}
	tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, bondDenom)
	asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, bondDenom)
	if tokenPrice.Equal(osmomath.ZeroBigDec()) {
		tokenPrice = k.amm.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
	}
	if found {
		for _, delegation := range delegations {
			amount := osmomath.BigDecFromDec(delegation.Shares)
			totalDelegations = totalDelegations.Add(amount.Mul(tokenPrice))
		}
	}

	// Max could be 7 for an account
	totalUnbondings := osmomath.ZeroBigDec()
	unbondingDelegations, err := k.stakingKeeper.GetUnbondingDelegations(ctx, user, 100)
	if err != nil {
		ctx.Logger().Error(err.Error())
		unbondingDelegations = []stakingtypes.UnbondingDelegation{}
	}
	if found {
		for _, delegation := range unbondingDelegations {
			for _, entry := range delegation.Entries {
				amount := osmomath.BigDecFromSDKInt(entry.Balance)
				totalUnbondings = totalUnbondings.Add(amount.Mul(tokenPrice))
			}
		}
	}
	return totalCommit, totalDelegations, totalUnbondings, totalVested
}

func (k Keeper) RetrieveRewardsTotal(ctx sdk.Context, user sdk.AccAddress) osmomath.BigDec {
	totalValue := osmomath.ZeroBigDec()
	estaking, err1 := k.estaking.Rewards(ctx, &estakingtypes.QueryRewardsRequest{Address: user.String()})
	masterchef, err2 := k.masterchef.UserPendingReward(ctx, &mastercheftypes.QueryUserPendingRewardRequest{User: user.String()})

	if err1 == nil {
		for _, balance := range estaking.Total {
			if balance.Denom == ptypes.Eden {
				balance.Denom = ptypes.Elys
			}
			tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, balance.Denom)
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, balance.Denom)
			if !found {
				continue
			}
			if tokenPrice.Equal(osmomath.ZeroBigDec()) {
				tokenPrice = k.amm.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := osmomath.BigDecFromSDKInt(balance.Amount)
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}

	if err2 == nil {
		for _, balance := range masterchef.TotalRewards {
			if balance.Denom == ptypes.Eden {
				balance.Denom = ptypes.Elys
			}
			tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, balance.Denom)
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, balance.Denom)
			if !found {
				continue
			}
			if tokenPrice.Equal(osmomath.ZeroBigDec()) {
				tokenPrice = k.amm.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := osmomath.BigDecFromSDKInt(balance.Amount)
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}
	return totalValue
}

func (k Keeper) RetrievePerpetualTotal(ctx sdk.Context, user sdk.AccAddress) (osmomath.BigDec, osmomath.BigDec, osmomath.BigDec) {
	totalAssets := osmomath.ZeroBigDec()
	totalLiability := osmomath.ZeroBigDec()
	var netValue osmomath.BigDec
	perpetuals, _, err := k.perpetual.GetMTPsForAddressWithPagination(ctx, user, nil)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec()
	}
	usdcDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec()
	}

	for _, perpetual := range perpetuals {
		if perpetual.Mtp.Position == perpetualtypes.Position_LONG {
			totalAssets = totalAssets.Add(k.amm.CalculateUSDValue(ctx, perpetual.Mtp.GetTradingAsset(), perpetual.Mtp.Custody))
			totalLiability = totalLiability.Add(k.amm.CalculateUSDValue(ctx, usdcDenom, perpetual.Mtp.Liabilities.Add(perpetual.Mtp.BorrowInterestUnpaidLiability)))
		} else {
			totalAssets = totalAssets.Add(k.amm.CalculateUSDValue(ctx, usdcDenom, perpetual.Mtp.Custody))
			totalLiability = totalLiability.Add(k.amm.CalculateUSDValue(ctx, perpetual.Mtp.LiabilitiesAsset, perpetual.Mtp.Liabilities.Add(perpetual.Mtp.BorrowInterestUnpaidLiability)))
		}
	}
	netValue = totalAssets.Sub(totalLiability)
	return totalAssets, totalLiability, netValue
}

func (k Keeper) RetrieveLiquidAssetsTotal(ctx sdk.Context, user sdk.AccAddress) osmomath.BigDec {
	balances := k.bankKeeper.GetAllBalances(ctx, user)
	totalValue := osmomath.ZeroBigDec()
	// Get eden from AmmBalance
	edenBal, err := k.amm.Balance(ctx, &ammtypes.QueryBalanceRequest{Denom: ptypes.Eden, Address: user.String()})
	if err == nil {
		balances = balances.Add(edenBal.Balance)
	}
	for _, balance := range balances {
		if balance.Denom == ptypes.Eden {
			balance.Denom = ptypes.Elys
		}
		tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, balance.Denom)
		asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, balance.Denom)
		if !found {
			continue
		}
		if tokenPrice.Equal(osmomath.ZeroBigDec()) {
			tokenPrice = k.amm.CalcAmmPrice(ctx, balance.Denom, asset.Decimals)
		}
		amount := osmomath.BigDecFromSDKInt(balance.Amount)
		totalValue = totalValue.Add(amount.Mul(tokenPrice))
	}
	return totalValue
}

func (k Keeper) RetrieveLeverageLpTotal(ctx sdk.Context, user sdk.AccAddress) (osmomath.BigDec, osmomath.BigDec, osmomath.BigDec) {
	positions := k.leveragelp.GetPositionsForAddress(ctx, user)
	totalValue := osmomath.ZeroBigDec()
	totalBorrow := osmomath.ZeroBigDec()
	netValue := osmomath.ZeroBigDec()
	for _, position := range positions {
		pool, found := k.amm.GetPool(ctx, position.AmmPoolId)
		if !found {
			continue
		}
		info := k.amm.PoolExtraInfo(ctx, pool, types.OneDay)
		amount := position.GetBigDecLeveragedLpAmount()
		totalValue = totalValue.Add(amount.Mul(info.GetBigDecLpTokenPrice()).Quo(ammtypes.OneShareBigDec))
		// USD value of debt
		debt := k.stablestakeKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
		collateralDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, position.Collateral.Denom)
		liab := debt.GetBigDecTotalLiablities()
		totalBorrow = totalBorrow.Add(liab.Mul(collateralDenomPrice))
	}
	netValue = totalValue.Sub(totalBorrow)
	return totalValue, totalBorrow, netValue
}

func (k Keeper) RetrieveTradeshieldTotal(ctx sdk.Context, user sdk.AccAddress) osmomath.BigDec {
	pendingStatus := tradeshieldtypes.Status_PENDING
	// Perpetual orders total
	perpetualOrders, _, err := k.tradeshieldKeeper.GetPendingPerpetualOrdersForAddress(ctx, user.String(), &pendingStatus, &query.PageRequest{})
	totalValue := osmomath.ZeroBigDec()
	if err == nil {
		for _, order := range perpetualOrders {
			balances := k.bankKeeper.GetAllBalances(ctx, order.GetOrderAddress())
			for _, balance := range balances {
				totalValue = totalValue.Add(k.amm.CalculateUSDValue(ctx, balance.Denom, balance.Amount))
			}
		}
	}

	// Spot orders total
	spotOrders, _, err := k.tradeshieldKeeper.GetPendingSpotOrdersForAddress(ctx, user.String(), &pendingStatus, &query.PageRequest{})
	if err == nil {
		for _, order := range spotOrders {
			balances := k.bankKeeper.GetAllBalances(ctx, order.GetOrderAddress())
			for _, balance := range balances {
				totalValue = totalValue.Add(k.amm.CalculateUSDValue(ctx, balance.Denom, balance.Amount))
			}
		}
	}
	return totalValue
}

func (k Keeper) RetrieveConsolidatedPrice(ctx sdk.Context, denom string) (osmomath.BigDec, osmomath.BigDec, sdkmath.LegacyDec) {
	if denom == ptypes.Eden {
		denom = ptypes.Elys
	}
	tokenPriceOracle := k.oracleKeeper.GetDenomPrice(ctx, denom)
	asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, denom)
	if !found {
		tokenPriceOracle = osmomath.ZeroBigDec()
	}
	tokenPriceAmm := k.amm.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
	info, found := k.oracleKeeper.GetAssetInfo(ctx, denom)
	tokenPriceOracleDec := sdkmath.LegacyZeroDec()
	if found {
		tokenPriceOracleD, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
		if found {
			tokenPriceOracleDec = tokenPriceOracleD
		}
	}

	return tokenPriceOracle, tokenPriceAmm, tokenPriceOracleDec
}

// SetPortfolio set a specific portfolio in the store from its index
func (k Keeper) SetPortfolio(ctx sdk.Context, portfolio types.Portfolio) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPortfolioKey(portfolio.Date, portfolio.GetCreatorAddress())
	b := k.cdc.MustMarshal(&portfolio)
	store.Set(key, b)
}

// GetPortfolio returns a portfolio from its index
func (k Keeper) GetPortfolio(ctx sdk.Context, user sdk.AccAddress, date string) (osmomath.BigDec, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPortfolioKey(date, user)
	portfolioBytes := store.Get(key)
	if portfolioBytes == nil {
		return osmomath.ZeroBigDec(), false
	}
	var val types.Portfolio
	k.cdc.MustUnmarshal(portfolioBytes, &val)
	return val.GetBigDecPortFolio(), true
}

func (k Keeper) GetMembershipTier(ctx sdk.Context, user sdk.AccAddress) (total_portfolio osmomath.BigDec, membership_tier types.MembershipTier) {
	year, month, day := ctx.BlockTime().Date()
	dateToday := time.Date(year, month, day, 0, 0, 0, 0, ctx.BlockTime().Location())
	startDate := dateToday.AddDate(0, 0, -7)
	minTotal := osmomath.NewBigDec(math.MaxInt64)
	for d := startDate; !d.After(dateToday); d = d.AddDate(0, 0, 1) {
		// Traverse all possible portfolio data
		totalPort, found := k.GetPortfolio(ctx, user, d.Format("2006-01-02"))
		if found && totalPort.LT(minTotal) {
			minTotal = totalPort
		}
	}

	if minTotal.Equal(osmomath.NewBigDec(math.MaxInt64)) {
		return osmomath.ZeroBigDec(), types.Basic
	}

	if minTotal.GTE(types.Platinum.GetBigDecMinimumPortfolio()) {
		return minTotal, types.Platinum
	}

	if minTotal.GTE(types.Gold.GetBigDecMinimumPortfolio()) {
		return minTotal, types.Gold
	}

	if minTotal.GTE(types.Silver.GetBigDecMinimumPortfolio()) {
		return minTotal, types.Silver
	}

	if minTotal.GTE(types.Bronze.GetBigDecMinimumPortfolio()) {
		return minTotal, types.Bronze
	}

	return minTotal, types.Basic
}

// RemovePortfolioLast removes a portfolio from the store with a specific date
func (k Keeper) RemovePortfolioLast(ctx sdk.Context, date string, num uint64) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.GetPortfolioByDateIteratorKey(date))

	defer iterator.Close()
	count := 0

	for ; iterator.Valid(); iterator.Next() {
		count++
		store.Delete(iterator.Key())
		if count == int(num) {
			break
		}
	}
	return uint64(count)
}

// GetAllPortfolio returns all portfolio
func (k Keeper) GetAllPortfolio(ctx sdk.Context) (list []types.Portfolio) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PortfolioKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Portfolio
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetDateFromContext(ctx sdk.Context) string {
	contextTime := ctx.BlockTime()
	// Extract the year, month, and day
	year, month, day := contextTime.Date()
	// Create a new time.Time object with the extracted date and time set to midnight
	contextDate := time.Date(year, month, day, 0, 0, 0, 0, contextTime.Location())
	// Format the date as a string in the "%Y-%m-%d" format
	return contextDate.Format(types.DateFormat)
}

func (k Keeper) GetDateAfterDaysFromContext(ctx sdk.Context, n int) string {
	contextTime := ctx.BlockTime()
	year, month, day := contextTime.Date()
	contextDate := time.Date(year, month, day, 0, 0, 0, 0, contextTime.Location())
	resultDate := contextDate.AddDate(0, 0, n)
	return resultDate.Format(types.DateFormat)
}

func GetPoolIdFromShareDenom(shareDenom string) (uint64, error) {
	poolId, err := strconv.Atoi(strings.TrimPrefix(shareDenom, "amm/pool/"))
	if err != nil {
		return 0, err
	}
	return uint64(poolId), nil
}

// remove after migrations
func (k Keeper) GetLegacyPortfolios(ctx sdk.Context, date string) (list []types.Portfolio) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(date))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyPortfolio
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, types.NewPortfolioWithContextDate(date, sdk.MustAccAddressFromBech32(val.Creator), val.Portfolio))
	}
	return
}

func (k Keeper) RemoveLegacyPortfolio(ctx sdk.Context, date string, user string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(date))
	store.Delete(types.PortfolioKey(user))
}

func (k Keeper) RemoveLegacyPortfolioCounted(ctx sdk.Context, date string, num uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(date))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	count := uint64(0)
	for ; iterator.Valid(); iterator.Next() {
		count++
		store.Delete(iterator.Key())
		if count == num {
			break
		}
	}
	return
}
