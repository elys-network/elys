package keeper_test

// TODO: implement
// func (server msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool)
// Oracle based pool - JoinPool
// 1) scenario1
// - JUNO/USDC pool
// - JUNO price $1
// - $1000 JUNO / $1000 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Join with JUNO token
// - Check weight breaking fee
// - Check slippage
// - Check bonus being zero
// 2) scenario2
// - JUNO/USDC pool
// - JUNO price $1
// - $1000 JUNO / $1000 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Join with JUNO+USDC token ($50 JUNO /$50 USDC)
// - Check weight breaking fee
// - Check slippage
// - Check bonus being zero
// 3) scenario3
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Join JUNO ($50)
// - Check weight breaking fee zero
// - Check slippage
// - Check bonus not be zero
// 4) scenario4
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Join USDC ($50)
// - Check weight breaking fee
// - Check slippage
// - Check bonus be zero
// 5) scenario5
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Join JUNO+USDC ($50 JUNO /$50 USDC)
// - Check weight breaking fee
// - Check slippage
// - Check bonus be zero
