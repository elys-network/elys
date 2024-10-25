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

	// Accounted pool balance = amm pool + (long liability - long profit taking liability) - (long custody - long profit taking custody) + (short liability - short profit taking liability ) - ( short custody - short profit taking custody)
	// Accounted pool balance = amm pool + totalLiabilities - totalCustody + total profit taking custody - total profit taking liability
	for i, asset := range accountedPool.PoolAssets {
		ammBalance, err := ammPool.GetAmmPoolBalance(asset.Token.Denom)
		if err != nil {
			return err
		}
		totalLiabilities, totalCustody, totalTakeProfitCustody, totalTakeProfitLiabilities := perpetualPool.GetPerpetualPoolBalances(asset.Token.Denom)
		accountedPoolAmt := ammBalance.Add(totalLiabilities).Sub(totalCustody)
		if EnableTakeProfitCustodyLiabilities {
			accountedPoolAmt = accountedPoolAmt.Add(totalTakeProfitCustody).Sub(totalTakeProfitLiabilities)
		}
		accountedPool.PoolAssets[i].Token = sdk.NewCoin(asset.Token.Denom, accountedPoolAmt)

		for j, nonAmmToken := range accountedPool.NonAmmPoolTokens {
			if nonAmmToken.Denom == asset.Token.Denom {
				accountedPool.NonAmmPoolTokens[j].Amount = accountedPoolAmt.Sub(ammBalance)
				break
			}
		}
	}

	// Update accounted pool
	k.SetAccountedPool(ctx, accountedPool)
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
