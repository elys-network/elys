#!/usr/bin/env bash

# list pool
elysd query leveragelp list-pool
elysd query leveragelp get-position $(elysd keys show -a treasury --keyring-backend=test) 1
elysd query leveragelp get-positions-for-address $(elysd keys show -a treasury --keyring-backend=test)
elysd query leveragelp get-positions
elysd query leveragelp params
elysd query stablestake borrow-ratio
elysd query stablestake params

# Open amm pool
elysd tx amm create-pool 10uatom,10uusdc 10000000000uatom,10000000000uusdc --use-oracle=true --from=treasury --yes --gas=1000000

# Put funds on stablestake
elysd tx stablestake bond 100000000 --from=treasury --yes --gas=1000000

# Local test
elysd tx gov submit-proposal proposal.json --from=treasury --yes --gas=1000000
elysd tx gov vote 1 Yes --from=treasury --yes --gas=1000000
elysd tx gov vote 3 Yes --from=treasury --yes --gas=1000000
elysd query gov proposals

# Open position
elysd tx leveragelp open 5.0 uusdc 5000000 1 0.0 --from=treasury --yes --gas=1000000 --fees=250uelys
elysd tx leveragelp open [leverage] [collateral-asset] [collateral-amount] [amm-pool-id] [flags]

# Close position
elysd tx leveragelp close 2 25000000000000000000 --from=treasury --yes --gas=1000000
elysd tx leveragelp close 1 500000000000000000 --from=treasury --yes --gas=1000000
elysd tx leveragelp close [position-id] [flags]

# Query rewards
elysd query leveragelp rewards $(elysd keys show -a treasury --keyring-backend=test) 1 --output=json
# {"rewards":[{"position_id":"1","reward":[{"denom":"ueden","amount":"3086835"}]}],"total_rewards":[{"denom":"ueden","amount":"3086835"}]}

# Claim rewards
elysd tx leveragelp claim-rewards 1 --from=treasury --yes --gas=1000000

elysd query commitment show-commitments $(elysd keys show -a treasury --keyring-backend=test)
elysd query commitment committed-tokens-locked $(elysd keys show -a treasury --keyring-backend=test)

# Testnet
elysd query oracle show-price ATOM  --node=https://rpc.testnet.elys.network:443
elysd query oracle show-price USDC  --node=https://rpc.testnet.elys.network:443
elysd query oracle list-asset-info --node=https://rpc.testnet.elys.network:443
asset_info:
- band_ticker: USDC
  decimal: "6"
  denom: ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65
  display: USDC
  elys_ticker: USDC
- band_ticker: ATOM
  decimal: "6"
  denom: ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4
  display: ATOM
  elys_ticker: ATOM

elysd query amm show-pool 2 --node=https://rpc.testnet.elys.network:443
pool:
  address: elys1t7z4shh8tzvjc2u9exu2fs8rmewlm6hza494x3dna0n7aumm05aq209wy9
  pool_assets:
  - token:
      amount: "95925788760"
      denom: ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65
    weight: "10737418240"
  - token:
      amount: "5466517678"
      denom: ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4
    weight: "10737418240"
  pool_id: "2"
  pool_params:
    exit_fee: "0.000000000000000000"
    external_liquidity_ratio: "228.549464829334853887"
    fee_denom: ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65
    swap_fee: "0.002000000000000000"
    threshold_weight_difference: "0.300000000000000000"
    use_oracle: true
    weight_breaking_fee_exponent: "2.500000000000000000"
    weight_breaking_fee_multiplier: "0.000200000000000000"
    weight_recovery_fee_portion: "0.100000000000000000"
  rebalance_treasury: elys1zfz2hcvzgcg2kgw0xyc9l27nmda6qnxahjrnjrateuejc84zf2kspdd9cl
  total_shares:
    amount: "143020000000000000000000"
    denom: amm/pool/2
  total_weight: "21474836480"

# Open 50 USDC position
elysd tx leveragelp open 5 ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65 50000000 2 --from=treasury --yes --gas=1000000 --node=https://rpc.testnet.elys.network:443 --fees=250uelys


# Put pool 2 on leveragelp (Testnet)
elysd tx gov submit-proposal proposal.json --from=t2a --yes --gas=1000000 --node=https://rpc.testnet.elys.network:443 --fees=250uelys
```
{
  "title": "Enable ATOM/USDC pool (PoolId 2) on leveragelp",
  "summary": "Enable ATOM/USDC pool (PoolId 2) on leveragelp",
  "messages": [
    {
        "@type": "/elys.leveragelp.MsgUpdatePools",
        "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
        "pools": [
            {
                "amm_pool_id": 2,
                "health": "0.0",
                "enabled": true,
                "closed": false,
                "leveraged_lp_amount": "0",
                "leverage_max": "10.0"
            }
        ]
    }
  ],
  "deposit": "10000000uelys"
}
```

elysd query leveragelp liquidation-price elys12tzylat4udvjj56uuhu3vj2n4vgp7cf9fwna9w 1
# price: "0.000000000000080000"

elysd query amm show-pool 1
# extra_info:
#   lp_token_price: "1.000000000000000000"
#   tvl: "110000.000000000000000000"
# pool:
#   address: elys12857uh2f6fus3x3jh07cp5u7xfg2ypeec27dcp6864xhraa33mdscw660j
#   pool_assets:
#   - token:
#       amount: "10000000000"