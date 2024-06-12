package keeper

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"

	"github.com/elys-network/elys/x/membershiptier/types"
)

func (k Keeper) RetreiveAllPortfolio(ctx sdk.Context, user string) {
	// set today + user -> amount
	sender := sdk.MustAccAddressFromBech32(user)
	todayDate := k.GetDateFromBlock(ctx.BlockTime())

	_, found := k.GetPortfolio(ctx, user, todayDate)
	if found {
		return
	}

	// Liquid assets
	balances := k.bankKeeper.GetAllBalances(ctx, sender)
	totalValue := sdk.NewDec(0)
	for _, balance := range balances {
		tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
		asset, found := k.assetProfileKeeper.GetEntry(ctx, balance.Denom)
		if !found {
			continue
		}
		amount := balance.Amount.ToLegacyDec().Quo(Pow10(asset.Decimals))
		totalValue = totalValue.Add(amount.Mul(tokenPrice))
	}

	// Rewards
	estaking, err1 := k.estaking.Rewards(ctx, &estakingtypes.QueryRewardsRequest{Address: user})
	masterchef, err2 := k.masterchef.UserPendingReward(ctx, &mastercheftypes.QueryUserPendingRewardRequest{User: user})

	if err1 == nil {
		for _, balance := range estaking.Total {
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
			asset, found := k.assetProfileKeeper.GetEntry(ctx, balance.Denom)
			if !found {
				continue
			}
			amount := balance.Amount.ToLegacyDec().Quo(Pow10(asset.Decimals))
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}

	if err2 == nil {
		for _, balance := range masterchef.TotalRewards {
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, balance.Denom)
			asset, found := k.assetProfileKeeper.GetEntry(ctx, balance.Denom)
			if !found {
				continue
			}
			amount := balance.Amount.ToLegacyDec().Quo(Pow10(asset.Decimals))
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}

	// Perpetual
	perpetuals, _, err := k.perpetual.GetMTPsForAddress(ctx, sender, &query.PageRequest{})
	if err == nil {
		for _, perpetual := range perpetuals {
			asset, found := k.assetProfileKeeper.GetEntry(ctx, perpetual.GetTradingAsset())
			if !found {
				continue
			}
			amount := perpetual.Custody.ToLegacyDec().Quo(Pow10(asset.Decimals))
			totalValue = totalValue.Add(amount)
		}
	}

	// Staked assets
	commitments := k.commitement.GetCommitments(ctx, user)
	for _, commitment := range commitments.CommittedTokens {
		// Pool balance
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				continue
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				return
			}
			info := k.amm.PoolExtraInfo(ctx, pool)
			amount := commitment.Amount.ToLegacyDec().Quo(Pow10(18))
			totalValue = totalValue.Add(amount.Mul(info.LpTokenPrice))
		} else {
			tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, commitment.Denom)
			asset, found := k.assetProfileKeeper.GetEntry(ctx, commitment.Denom)
			if !found {
				continue
			}
			amount := commitment.Amount.ToLegacyDec().Quo(Pow10(asset.Decimals))
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}

	// Delegations
	delegations := k.stakingKeeper.GetAllDelegatorDelegations(ctx, sender)
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, bondDenom)
	asset, found := k.assetProfileKeeper.GetEntry(ctx, bondDenom)
	if found {
		for _, delegation := range delegations {
			amount := delegation.Shares.Quo(Pow10(asset.Decimals))
			totalValue = totalValue.Add(amount.Mul(tokenPrice))
		}
	}
	// Max could be 7 for an account
	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegations(ctx, sender, 100)
	if found {
		for _, delegation := range unbondingDelegations {
			for _, entry := range delegation.Entries {
				amount := entry.Balance.ToLegacyDec().Quo(Pow10(asset.Decimals))
				totalValue = totalValue.Add(amount.Mul(tokenPrice))
			}
		}
	}

	k.SetPortfolio(ctx, todayDate, sender.String(), types.Portfolio{
		Creator:   user,
		Portfolio: totalValue,
	})
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
