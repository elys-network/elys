package keeper_test

// TODO: implement
// func (server msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn)
// Oracle based pool - Swap
// 1) scenario1
// - JUNO/USDC pool
// - JUNO price $1
// - $1000 JUNO / $1000 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Swap JUNO -> USDC
// - Check weight breaking fee
// - Check slippage
// - Check bonus being zero
// 2) scenario2
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Swap JUNO -> USDC
// - Check weight breaking fee zero
// - Check slippage
// - Check bonus not be zero
// 3) scenario3
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Swap USDC -> JUNO
// - Check weight breaking fee
// - Check slippage
// - Check bonus be zero
