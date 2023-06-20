package keeper_test

// TODO: implement
// func (server msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool)
// Oracle based pool - Exit
// 1) scenario1
// - JUNO/USDC pool
// - JUNO price $1
// - $1000 JUNO / $1000 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Exit JUNO
// - Check weight breaking fee
// - Check slippage
// 2) scenario2
// - JUNO/USDC pool
// - JUNO price $1
// - $1000 JUNO / $1000 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Exit JUNO+USDC
// - Check weight breaking fee
// - Check slippage
// 3) scenario3
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Exit JUNO
// - Check weight breaking fee zero
// - Check slippage
// 4) scenario4
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Exit USDC
// - Check weight breaking fee
// - Check slippage
// 5) scenario5
// - JUNO/USDC pool
// - JUNO price $1
// - $500 JUNO / $1500 USDC
// - External liquidity on Osmosis $10000 JUNO / $10000 USDC
// - Slippage reduction: 90%
// - Target Weight 50% / 50%
// - Weight Broken Threshold: 20%
// - Exit JUNO+USDC
// - Check weight breaking fee
// - Check slippage

// TODO: check combined scenario - $500 JUNO sell, $500 JUNO buy (weight breaking)
// TODO: check combined scenario - $500 JUNO sell, $500 JUNO buy (weight not breaking)
// TODO: run simulation test with a lot of traffic, and see pool status after the execution
// TODO: Check maximum weight breaking fee applied
// TODO: Check maximum weight recovery bonus applied
// TODO: Check weight recovery treasury empty case
// TODO: handle case bonus pool does not have enough amount
// TODO: check fee distribution
// TODO: write table driven data on spec folder for various cases to show the comparison with Osmosis
