package keeper

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	"github.com/elys-network/elys/v6/x/amm/types"
)

func (k Keeper) GetStackedSlippage(ctx sdk.Context, poolId uint64) osmomath.BigDec {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return osmomath.ZeroBigDec()
	}
	snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)
	return pool.StackedRatioFromSnapshot(snapshot)
}

func (k Keeper) ApplySwapRequest(ctx sdk.Context, msg sdk.Msg) error {
	switch msg := msg.(type) {
	case *types.MsgSwapExactAmountIn:
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return err
		}
		recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
		if err != nil {
			recipient = sender
		}
		_, _, _, err = k.RouteExactAmountIn(ctx, sender, recipient, msg.Routes, msg.TokenIn, msg.TokenOutMinAmount)
		if err != nil {
			return err
		}
		return nil
	case *types.MsgSwapExactAmountOut:
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return err
		}
		recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
		if err != nil {
			recipient = sender
		}
		_, _, _, err = k.RouteExactAmountOut(ctx, sender, recipient, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
		if err != nil {
			return err
		}
		return nil
	default:
		return types.ErrInvalidSwapMsgType
	}
}

func (k Keeper) DeleteSwapRequest(ctx sdk.Context, msg sdk.Msg, index uint64) {
	switch msg := msg.(type) {
	case *types.MsgSwapExactAmountIn:
		k.DeleteSwapExactAmountInRequest(ctx, msg, index)
	case *types.MsgSwapExactAmountOut:
		k.DeleteSwapExactAmountOutRequest(ctx, msg, index)
	}
}

func (k Keeper) SelectOneSwapRequest(ctx sdk.Context, sprefix []byte) (sdk.Msg, uint64) {
	msg1, index := k.GetFirstSwapExactAmountInRequest(ctx, sprefix)
	if index != 0 {
		return msg1, index
	}
	msg2, index := k.GetFirstSwapExactAmountOutRequest(ctx, sprefix)
	return msg2, index
}

func (k Keeper) SelectReverseSwapRequest(ctx sdk.Context, msg sdk.Msg) (sdk.Msg, uint64) {
	sprefix := []byte{}
	switch msg := msg.(type) {
	case *types.MsgSwapExactAmountIn:
		sprefix = types.TKeyPrefixSwapExactAmountInPrefix(msg)
	case *types.MsgSwapExactAmountOut:
		sprefix = types.TKeyPrefixSwapExactAmountOutPrefix(msg)
	}

	split := strings.Split(string(sprefix), "/")
	for i, j := 0, len(split)-1; i < j; i, j = i+1, j-1 {
		split[i], split[j] = split[j], split[i]
	}
	rprefix := strings.Join(split, "/")
	return k.SelectOneSwapRequest(ctx, []byte(rprefix))
}

func (k Keeper) FirstPoolId(msg sdk.Msg) uint64 {
	switch msg := msg.(type) {
	case *types.MsgSwapExactAmountIn:
		return types.FirstPoolIdFromSwapExactAmountIn(msg)
	case *types.MsgSwapExactAmountOut:
		return types.FirstPoolIdFromSwapExactAmountOut(msg)
	}
	return 0
}

func (k Keeper) ExecuteSwapRequests(ctx sdk.Context) []sdk.Msg {
	// Algorithm
	// - Select a random swap request
	//   - Try execution on cache context, and check stacked slippage
	//   - Check if opposite direction request exists (Same pool id with opposite in/out tokens)
	//   - If opposite direction request exists, try execution on cache context, and check stacked slippage
	//   - Apply the swap request which as lower stacked slippage
	//   - If one of the swaps fail, not apply any changes and remove the swap request
	// - Repeat the process until the swap requests run-out
	requests := []sdk.Msg{}
	for {
		var index1, index2 uint64
		var msg1, msg2 sdk.Msg
		msg1, index1 = k.SelectOneSwapRequest(ctx, []byte{})
		if index1 == 0 {
			break
		}

		msg2, index2 = k.SelectReverseSwapRequest(ctx, msg1)
		if index2 == 0 {
			cachedCtx, write := ctx.CacheContext()
			err := k.ApplySwapRequest(cachedCtx, msg1)
			if err == nil {
				write()
			}
			// remove msg1 from the store
			k.DeleteSwapRequest(ctx, msg1, index1)
			requests = append(requests, msg1)
			continue
		}

		poolId := k.FirstPoolId(msg1)
		cachedCtx1, write1 := ctx.CacheContext()
		err1 := k.ApplySwapRequest(cachedCtx1, msg1)
		stackedSlippage1 := k.GetStackedSlippage(cachedCtx1, poolId)

		cachedCtx2, write2 := ctx.CacheContext()
		err2 := k.ApplySwapRequest(cachedCtx2, msg2)
		stackedSlippage2 := k.GetStackedSlippage(cachedCtx2, poolId)

		if err1 == nil && err2 == nil {
			if stackedSlippage1.LT(stackedSlippage2) {
				write1()
				// remove msg1 from the store
				k.DeleteSwapRequest(ctx, msg1, index1)
				requests = append(requests, msg1)
			} else {
				write2()
				// remove msg2 from the store
				k.DeleteSwapRequest(ctx, msg2, index2)
				requests = append(requests, msg2)
			}
		} else if err1 == nil {
			// remove msg2 from the store
			k.DeleteSwapRequest(ctx, msg2, index2)
			requests = append(requests, msg2)
		} else if err2 == nil {
			// remove msg1 from the store
			k.DeleteSwapRequest(ctx, msg1, index1)
			requests = append(requests, msg1)
		} else {
			// remove both msg1, msg2 messages
			k.DeleteSwapRequest(ctx, msg1, index1)
			k.DeleteSwapRequest(ctx, msg2, index2)
			requests = append(requests, msg1, msg2)
		}
	}
	return requests
}

func (k Keeper) ClearOutdatedSlippageTrack(ctx sdk.Context) {
	params := k.GetParams(ctx)
	tracks := k.AllSlippageTracks(ctx)
	for _, track := range tracks {
		if track.Timestamp+params.SlippageTrackDuration < uint64(ctx.BlockTime().Unix()) {
			k.DeleteSlippageTrack(ctx, track)
		}
	}
}

// EndBlocker of amm module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	msgs := k.ExecuteSwapRequests(ctx)
	if len(msgs) > 0 {
		bz, _ := json.Marshal(msgs)
		k.Logger(ctx).Debug("Executed swap requests: " + string(bz))
	}

	// Set amm and accounted pools in oracle kv store
	// TODO this is being used for price feeder, migrate to query in price feeder and the remove this
	//ammPools := k.GetAllPool(ctx)
	//for _, ammPool := range ammPools {
	//	if ammPool.PoolParams.UseOracle {
	//		oraclePool := oracletypes.Pool{
	//			PoolId: ammPool.PoolId,
	//		}
	//
	//		oraclePoolAssets := make([]oracletypes.PoolAsset, 0)
	//		for _, poolAsset := range ammPool.PoolAssets {
	//			oraclePoolAssets = append(oraclePoolAssets, oracletypes.PoolAsset{
	//				Token:                  poolAsset.Token,
	//				Weight:                 poolAsset.Weight,
	//				ExternalLiquidityRatio: poolAsset.ExternalLiquidityRatio,
	//			})
	//		}
	//		oraclePool.PoolAssets = oraclePoolAssets
	//		k.oracleKeeper.SetPool(ctx, oraclePool)
	//	}
	//}

	k.ClearOutdatedSlippageTrack(ctx)
}
