#!/usr/bin/env bash

# list pool
elysd query leveragelp list-pool
elysd query leveragelp get-positions-for-address $(elysd keys show -a treasury --keyring-backend=test)
elysd query leveragelp get-positions

# Put pool 2 on leveragelp (Testnet)
elysd tx gov submit-proposal proposal.json --from=t2a --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000 --node=https://rpc.testnet.elys.network:443 --fees=250uelys
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

# Local test
elysd tx gov submit-proposal proposal.json --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
elysd tx gov vote 1 Yes --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# Open 100 USDC position
elysd tx leveragelp open 5 ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65 100000000 2 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
elysd tx leveragelp open [leverage] [collateral-asset] [collateral-amount] [amm-pool-id] [flags]

# Close position
elysd tx leveragelp close 1 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
elysd tx leveragelp close [position-id] [flags]
