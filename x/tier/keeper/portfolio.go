package keeper

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"

	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/tier/types"
)

func (k Keeper) RetreiveAllPortfolio(ctx sdk.Context, user string) {
	// set today + user -> amount
	sender := sdk.MustAccAddressFromBech32(user)
	todayDate := k.GetDateFromBlock(ctx.BlockTime())

	_, found := k.GetPortfolio(ctx, user, todayDate)
	if found {
		return
	}

	totalValue := sdk.NewDec(0)

	// Liquid assets
	liq := k.RetreiveLiquidAssetsTotal(ctx, sender)
	totalValue = totalValue.Add(liq)

	// Rewards
	rew := k.RetreiveRewardsTotal(ctx, sender)
	totalValue = totalValue.Add(rew)

	// Perpetual
	perp := k.RetreivePerpetualTotal(ctx, sender)
	totalValue = totalValue.Add(perp)

	// Staked+Pool assets
	staked := k.RetreiveStakedAndPoolTotal(ctx, sender)
	totalValue = totalValue.Add(staked)

	// LeverageLp
	lev := k.RetreiveLeverageLpTotal(ctx, sender)
	totalValue = totalValue.Add(lev)

	k.SetPortfolio(ctx, todayDate, sender.String(), types.Portfolio{
		Creator:   user,
		Portfolio: totalValue,
	})
}

func (k Keeper) RetreiveStakedAndPoolTotal(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	totalValue := sdk.NewDec(0)
	commitments := k.commitement.GetCommitments(ctx, user.String())
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
			info := k.amm.PoolExtraInfo(ctx, pool)
			amount := commitment.Amount.ToLegacyDec()
			totalValue = totalValue.Add(amount.Mul(info.LpTokenPrice).QuoInt(ammtypes.OneShare))
		} else {
			if commitment.Denom == "ueden" {
				commitment.Denom = "uelys"
			}
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, commitment.Denom)
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, commitment.Denom)
			if !found {
				continue
			}
			if tokenPrice == sdk.ZeroDec() {
				tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := commitment.Amount.ToLegacyDec()
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}

	// Delegations
	delegations := k.stakingKeeper.GetAllDelegatorDelegations(ctx, user)
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, bondDenom)
	asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, bondDenom)
	if tokenPrice == sdk.ZeroDec() {
		tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
	}
	if found {
		for _, delegation := range delegations {
			amount := delegation.Shares
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}
	// Max could be 7 for an account
	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegations(ctx, user, 100)
	if found {
		for _, delegation := range unbondingDelegations {
			for _, entry := range delegation.Entries {
				amount := entry.Balance.ToLegacyDec()
				totalValue = totalValue.Add(amount.Mul(tokenPrice))
			}
		}
	}
	return totalValue
}

func (k Keeper) RetreiveRewardsTotal(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	totalValue := sdk.NewDec(0)
	estaking, err1 := k.estaking.Rewards(ctx, &estakingtypes.QueryRewardsRequest{Address: user.String()})
	masterchef, err2 := k.masterchef.UserPendingReward(ctx, &mastercheftypes.QueryUserPendingRewardRequest{User: user.String()})

	if err1 == nil {
		for _, balance := range estaking.Total {
			if balance.Denom == "ueden" {
				balance.Denom = "uelys"
			}
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, balance.Denom)
			if !found {
				continue
			}
			if tokenPrice == sdk.ZeroDec() {
				tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := balance.Amount.ToLegacyDec()
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}

	if err2 == nil {
		for _, balance := range masterchef.TotalRewards {
			if balance.Denom == "ueden" {
				balance.Denom = "uelys"
			}
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, balance.Denom)
			if !found {
				continue
			}
			if tokenPrice == sdk.ZeroDec() {
				tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := balance.Amount.ToLegacyDec()
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}
	return totalValue
}

func (k Keeper) RetreivePerpetualTotal(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	totalValue := sdk.NewDec(0)
	perpetuals, _, err := k.perpetual.GetMTPsForAddress(ctx, user, &query.PageRequest{})
	if err == nil {
		for _, perpetual := range perpetuals {
			asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, perpetual.GetTradingAsset())
			if !found {
				continue
			}
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
			if tokenPrice == sdk.ZeroDec() {
				tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
			}
			amount := perpetual.Custody.ToLegacyDec()
			totalValue = totalValue.Add((amount.Mul(tokenPrice)))
		}
	}
	return totalValue
}

func (k Keeper) RetreiveLiquidAssetsTotal(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	balances := k.bankKeeper.GetAllBalances(ctx, user)
	totalValue := sdk.NewDec(0)
	// Get eden from AmmBalance
	edenBal, err := k.amm.Balance(ctx, &ammtypes.QueryBalanceRequest{Denom: ptypes.Eden, Address: user.String()})
	if err == nil {
		balances = balances.Add(edenBal.Balance)
	}
	for _, balance := range balances {
		if balance.Denom == "ueden" {
			balance.Denom = "uelys"
		}
		tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
		asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, balance.Denom)
		if !found {
			continue
		}
		if tokenPrice == sdk.ZeroDec() {
			tokenPrice = k.CalcAmmPrice(ctx, balance.Denom, asset.Decimals)
		}
		amount := balance.Amount.ToLegacyDec()
		totalValue = totalValue.Add(amount.Mul(tokenPrice))
	}
	return totalValue
}

func (k Keeper) RetreiveLeverageLpTotal(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	positions, _, err := k.leveragelp.GetPositionsForAddress(ctx, user, &query.PageRequest{})
	totalValue := sdk.NewDec(0)
	if err == nil {
		for _, position := range positions {
			pool, found := k.amm.GetPool(ctx, position.Position.AmmPoolId)
			if !found {
				continue
			}
			info := k.amm.PoolExtraInfo(ctx, pool)
			amount := position.Position.LeveragedLpAmount.ToLegacyDec()
			totalValue = totalValue.Add(amount.Mul(info.LpTokenPrice).QuoInt(ammtypes.OneShare))
		}
	}
	return totalValue
}

func (k Keeper) RetreiveConsolidatedPrice(ctx sdk.Context, denom string) (sdk.Dec, error) {
	tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, denom)
	asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, denom)
	if !found {
		return sdk.ZeroDec(), types.ErrNotFound
	}
	if tokenPrice == sdk.ZeroDec() {
		tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
	}
	return tokenPrice, nil
}

func (k Keeper) CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) sdk.Dec {
	usdcDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return sdk.ZeroDec()
	}
	usdcPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, usdcDenom)
	resp, err := k.amm.InRouteByDenom(sdk.WrapSDKContext(ctx), &ammtypes.QueryInRouteByDenomRequest{DenomIn: denom, DenomOut: usdcDenom})
	if err != nil {
		return sdk.ZeroDec()
	}

	routes := resp.InRoute
	tokenIn := sdk.NewCoin(denom, sdk.NewInt(Pow10(decimal).TruncateInt64()))
	discount := sdk.NewDec(1)
	spotPrice, _, _, _, _, _, _, _, err := k.amm.CalcInRouteSpotPrice(ctx, tokenIn, routes, discount, sdk.ZeroDec())
	if err != nil {
		return sdk.ZeroDec()
	}
	return spotPrice.Mul(usdcPrice)
}

// SetPortfolio set a specific portfolio in the store from its index
func (k Keeper) SetPortfolio(ctx sdk.Context, todayDate string, user string, portfolio types.Portfolio) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(todayDate))
	b := k.cdc.MustMarshal(&portfolio)
	store.Set(types.PortfolioKey(
		user,
	), b)
}

// GetPortfolio returns a portfolio from its index
func (k Keeper) GetPortfolio(
	ctx sdk.Context,
	user string,
	timestamp string,
) (sdk.Dec, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(timestamp))

	found := store.Has(types.PortfolioKey(
		user,
	))

	if !found {
		return sdk.NewDec(0), false
	}

	portfolio := store.Get(types.PortfolioKey(
		user,
	))
	var val types.Portfolio
	k.cdc.MustUnmarshal(portfolio, &val)
	return val.Portfolio, true
}

func (k Keeper) GetMembershipTier(ctx sdk.Context, user string) (total_portfoilio sdk.Dec, tier string, discount uint64) {
	year, month, day := ctx.BlockTime().Date()
	dateToday := time.Date(year, month, day, 0, 0, 0, 0, ctx.BlockTime().Location())
	startDate := dateToday.AddDate(0, 0, -7)
	minTotal := sdk.NewDec(math.MaxInt64)
	for d := startDate; !d.After(dateToday); d = d.AddDate(0, 0, 1) {
		// Traverse all possible portfolio data
		totalPort, found := k.GetPortfolio(ctx, user, d.Format("2006-01-02"))
		if found && totalPort.LT(minTotal) {
			minTotal = totalPort
		}
	}

	if minTotal.Equal(sdk.NewDec(math.MaxInt64)) {
		return sdk.NewDec(0), "bronze", 0
	}

	// TODO: Make tier discount and minimum balance configurable
	if minTotal.GTE(sdk.NewDec(500000)) {
		return minTotal, "platinum", 30
	}

	if minTotal.GTE(sdk.NewDec(250000)) {
		return minTotal, "gold", 20
	}

	if minTotal.GTE(sdk.NewDec(50000)) {
		return minTotal, "silver", 10
	}

	return minTotal, "bronze", 0
}

// RemovePortfolio removes a portfolio from the store
func (k Keeper) RemovePortfolio(
	ctx sdk.Context,
	user string,
	timestamp string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(timestamp))
	store.Delete(types.PortfolioKey(
		user,
	))
}

// RemovePortfolioLast removes a portfolio from the store with a specific date
func (k Keeper) RemovePortfolioLast(
	ctx sdk.Context,
	timestamp string,
	num uint64,
) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(timestamp))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

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
func (k Keeper) GetAllPortfolio(ctx sdk.Context, timestamp string) (list []types.Portfolio) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(timestamp))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Portfolio
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetDateFromBlock(blockTime time.Time) string {
	// Extract the year, month, and day
	year, month, day := blockTime.Date()
	// Create a new time.Time object with the extracted date and time set to midnight
	blockDate := time.Date(year, month, day, 0, 0, 0, 0, blockTime.Location())
	// Format the date as a string in the "%Y-%m-%d" format
	return blockDate.Format("2006-01-02")
}

func GetPoolIdFromShareDenom(shareDenom string) (uint64, error) {
	poolId, err := strconv.Atoi(strings.TrimPrefix(shareDenom, "amm/pool/"))
	if err != nil {
		return 0, err
	}
	return uint64(poolId), nil
}

func Pow10(decimal uint64) (value sdk.Dec) {
	value = sdk.NewDec(1)
	for i := 0; i < int(decimal); i++ {
		value = value.Mul(sdk.NewDec(10))
	}
	return
}
