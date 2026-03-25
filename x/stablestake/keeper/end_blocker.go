package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	//allPools := k.GetAllPools(ctx)
	//msgServerImp := NewMsgServerImpl(k)
	//maxCount := 50
	//
	//gasMeter := ctx.BlockGasMeter()
	//maxBlockGas := gasMeter.Limit()
	//safeGasThreshold := uint64(0)
	//
	//if maxBlockGas > 0 {
	//	safeGasThreshold = maxBlockGas * 2 / 5
	//}
	//
	//for _, pool := range allPools {
	//	denom := types.GetShareDenomForPool(pool.Id)
	//	startAddr := k.GetLastProccessed(ctx, pool.Id)
	//
	//	scannedCount := 0
	//	var list []commitmenttypes.Commitments
	//	var lastScannedAddr sdk.AccAddress // Tracks the furthest valid address scanned
	//
	//	// Iterate and collect commitments safely
	//	k.GetCommitmentKeeper().IterateCommitmentsFromAddress(ctx, startAddr, func(commitment commitmenttypes.Commitments) (stop bool) {
	//		// Emergency gas check during iteration
	//		if maxBlockGas > 0 && gasMeter.GasConsumed() > safeGasThreshold {
	//			return true
	//		}
	//
	//		currentAddr, err := sdk.AccAddressFromBech32(commitment.Creator)
	//		if err != nil {
	//			k.Logger(ctx).Error("Invalid creator address skipped", "creator", commitment.Creator)
	//			scannedCount++ // MUST increment to prevent infinite loop
	//			return scannedCount >= maxCount
	//		}
	//
	//		lastScannedAddr = currentAddr
	//
	//		if startAddr != nil && bytes.Equal(startAddr, currentAddr) {
	//			return false
	//		}
	//
	//		scannedCount++
	//
	//		// Filter logic
	//		for _, v := range commitment.CommittedTokens {
	//			if v.Denom == denom {
	//				list = append(list, commitment)
	//				break
	//			}
	//		}
	//
	//		return scannedCount >= maxCount
	//	})
	//
	//	var lastSuccessfullyProcessed sdk.AccAddress = startAddr
	//	gasExceeded := false
	//
	//	// Execute the Unbonds
	//	for _, commitment := range list {
	//		// Gas Brake Check before executing heavy logic
	//		if maxBlockGas > 0 && gasMeter.GasConsumed() > safeGasThreshold {
	//			gasExceeded = true
	//			break
	//		}
	//
	//		amount := math.ZeroInt()
	//		for _, v := range commitment.CommittedTokens {
	//			if v.Denom == denom {
	//				amount = v.Amount
	//				break
	//			}
	//		}
	//
	//		if amount.IsPositive() && amount.GT(math.OneInt()) {
	//			func() {
	//				defer func() {
	//					if r := recover(); r != nil {
	//						k.Logger(ctx).Error("Panic recovered during Unbond", "poolId", pool.Id, "creator", commitment.Creator, "panic", r)
	//					}
	//				}()
	//
	//				cacheCtx, write := ctx.CacheContext()
	//				_, err := msgServerImp.Unbond(cacheCtx, types.NewMsgUnbond(commitment.Creator, amount.QuoRaw(2), pool.Id))
	//
	//				if err == nil {
	//					write()
	//				} else {
	//					k.Logger(ctx).Error("Unbond failed gracefully", "poolId", pool.Id, "creator", commitment.Creator, "err", err)
	//				}
	//			}()
	//		}
	//
	//		lastSuccessfullyProcessed, _ = sdk.AccAddressFromBech32(commitment.Creator)
	//	}
	//
	//	if gasExceeded {
	//		// If we hit the gas limit mid-list, save the last successful item.
	//		if lastSuccessfullyProcessed != nil && !bytes.Equal(startAddr, lastSuccessfullyProcessed) {
	//			k.SetLastProccessed(ctx, pool.Id, lastSuccessfullyProcessed)
	//		}
	//	} else {
	//		// We finished the list successfully (or it was empty).
	//		if scannedCount < maxCount {
	//			// Reached the end of the KVStore for this pool
	//			k.DeleteLastProccessed(ctx, pool.Id)
	//		} else if lastScannedAddr != nil {
	//			// Hit the 50 item limit. Save the furthest address scanned.
	//			k.SetLastProccessed(ctx, pool.Id, lastScannedAddr)
	//		}
	//	}
	//
	//	// 3. Break the main pool loop if gas is low
	//	if gasExceeded || (maxBlockGas > 0 && gasMeter.GasConsumed() > safeGasThreshold) {
	//		k.Logger(ctx).Info("Gas safety threshold reached, pausing EndBlocker")
	//		return
	//	}
	//}
}
