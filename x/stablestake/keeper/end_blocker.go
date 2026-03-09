package keeper

import (
	"bytes"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/v6/x/commitment/types"
	"github.com/elys-network/elys/v6/x/stablestake/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	allPools := k.GetAllPools(ctx)
	msgServerImp := NewMsgServerImpl(k)
	maxCount := 50
	for _, pool := range allPools {
		denom := types.GetShareDenomForPool(pool.Id)
		startAddr := k.GetLastProccessed(ctx, pool.Id)

		count := 0
		var list []commitmenttypes.Commitments
		var lastProcessedAddr sdk.AccAddress

		k.GetCommitmentKeeper().IterateCommitmentsFromAddress(ctx, startAddr, func(commitment commitmenttypes.Commitments) (stop bool) {
			currentAddr, err := sdk.AccAddressFromBech32(commitment.Creator)
			if err != nil {
				return false // Skip invalid addresses
			}

			if startAddr != nil && bytes.Equal(startAddr, currentAddr) {
				return false
			}

			count++
			lastProcessedAddr = currentAddr

			// Filter logic
			for _, v := range commitment.CommittedTokens {
				if v.Denom == denom {
					list = append(list, commitment)
					break
				}
			}

			if count == maxCount {
				return true
			}
			return false
		})

		for _, commitment := range list {
			amount := math.ZeroInt()
			for _, v := range commitment.CommittedTokens {
				if v.Denom == denom {
					amount = v.Amount
					break
				}
			}

			if amount.IsPositive() && amount.GT(math.OneInt()) {
				cacheCtx, write := ctx.CacheContext()
				_, err := msgServerImp.Unbond(cacheCtx, types.NewMsgUnbond(commitment.Creator, amount.QuoRaw(2), pool.Id))
				if err == nil {
					write()
				}
			}
		}

		if count == maxCount {
			k.SetLastProccessed(ctx, pool.Id, lastProcessedAddr)
		} else {
			k.DeleteLastProccessed(ctx, pool.Id)
		}
	}
}
