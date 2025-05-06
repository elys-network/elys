package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) PerpetualUpdates(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, EnableTakeProfitCustodyLiabilities bool) error {
	// Get accounted pool
	accountedPool, found := k.GetAccountedPool(ctx, ammPool.PoolId)
	if !found {
		return types.ErrPoolDoesNotExist
	}

	//var netTotalLiabilities, netTotalCustody, netTotalTakeProfitCustody, netTotalTakeProfitLiabilities sdk.Coins
	// Accounted pool balance = amm pool + (long liability - long profit taking liability) - (long custody - long profit taking custody) + (short liability - short profit taking liability ) - ( short custody - short profit taking custody)
	// Accounted pool balance = amm pool + totalLiabilities - totalCustody + total profit taking custody - total profit taking liability
	for i, asset := range accountedPool.TotalTokens {
		ammBalance, err := ammPool.GetAmmPoolBalance(asset.Denom)
		if err != nil {
			return err
		}
		totalLiabilities, totalCustody, totalTakeProfitCustody, totalTakeProfitLiabilities := perpetualPool.GetPerpetualPoolBalances(asset.Denom)
		accountedPoolAmt := ammBalance.Add(totalLiabilities).Sub(totalCustody)
		// if this is enabled then we need to consider impact of weight balance bonus of swap while calculating takeProfitLiabilities
		if EnableTakeProfitCustodyLiabilities {
			accountedPoolAmt = accountedPoolAmt.Add(totalTakeProfitCustody).Sub(totalTakeProfitLiabilities)
		}
		if !totalLiabilities.IsZero() {

		}
		//netTotalCustody = netTotalCustody.Add(sdk.NewCoin(asset.Denom, totalCustody))
		//netTotalLiabilities = netTotalLiabilities.Add(sdk.NewCoin(asset.Denom, totalLiabilities))
		//netTotalTakeProfitCustody = netTotalTakeProfitCustody.Add(sdk.NewCoin(asset.Denom, totalTakeProfitCustody))
		//netTotalTakeProfitLiabilities = netTotalTakeProfitLiabilities.Add(sdk.NewCoin(asset.Denom, totalTakeProfitLiabilities))

		accountedPool.TotalTokens[i] = sdk.NewCoin(asset.Denom, accountedPoolAmt)

		for j, nonAmmToken := range accountedPool.NonAmmPoolTokens {
			if nonAmmToken.Denom == asset.Denom {
				accountedPool.NonAmmPoolTokens[j].Amount = accountedPoolAmt.Sub(ammBalance)
				break
			}
		}
	}

	// Update accounted pool
	k.SetAccountedPool(ctx, accountedPool)

	//ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventPerpetualUpdates,
	//	sdk.NewAttribute("pool_id", strconv.FormatUint(ammPool.PoolId, 10)),
	//	sdk.NewAttribute("net_total_custody", netTotalCustody.String()),
	//	sdk.NewAttribute("net_total_take_profit_liabilities", netTotalLiabilities.String()),
	//	sdk.NewAttribute("net_total_take_profit_custody", netTotalTakeProfitCustody.String()),
	//	sdk.NewAttribute("net_total_take_profit_liabilities", netTotalTakeProfitLiabilities.String()),
	//	sdk.NewAttribute("non_amm_token_balance", sdk.Coins(accountedPool.NonAmmPoolTokens).String()),
	//	sdk.NewAttribute("total_tokens", sdk.Coins(accountedPool.TotalTokens).String()),
	//))

	return nil
}

// Hooks wrapper struct for tvl keeper
type PerpetualHooks struct {
	k Keeper
}

var _ perpetualtypes.PerpetualHooks = PerpetualHooks{}

// Return the wrapper struct
func (k Keeper) PerpetualHooks() PerpetualHooks {
	return PerpetualHooks{k}
}

func (h PerpetualHooks) AfterParamsChange(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, EnableTakeProfitCustodyLiabilities bool) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool, EnableTakeProfitCustodyLiabilities)
}

func (h PerpetualHooks) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool, EnableTakeProfitCustodyLiabilities)
}

func (h PerpetualHooks) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool, EnableTakeProfitCustodyLiabilities)
}

func (h PerpetualHooks) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool, EnableTakeProfitCustodyLiabilities)
}
