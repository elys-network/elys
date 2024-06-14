<!--
order: 2
-->

# State

The `GenesisState` message defines the initial state of the `masterchef` module at genesis.

```proto
// GenesisState defines the masterchef module's genesis state.
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated ExternalIncentive external_incentives = 2 [(gogoproto.nullable) = false];
  uint64 external_incentive_index = 3;
  repeated PoolInfo pool_infos = 4 [(gogoproto.nullable) = false];
  repeated PoolRewardInfo pool_reward_infos = 5 [(gogoproto.nullable) = false];
  repeated UserRewardInfo user_reward_infos = 6 [(gogoproto.nullable) = false];
  repeated PoolRewardsAccum pool_rewards_accum = 7 [(gogoproto.nullable) = false];
}
```

## Params

`Params` holds module parameters like protocol revenue address, supported reward denoms, max eden apr, reward portion for lps and stakers and yearly incentive data.

```proto
// Params defines the parameters for the module.
message Params {
  option (gogoproto.goproto_stringer) = false;
  IncentiveInfo lp_incentives = 1;

  // gas fees and swap fees portion for lps, `100 - reward_portion_for_lps - reward_portion_for_stakers = revenue percent for protocol`.
  string reward_portion_for_lps = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];

  // gas fees and swap fees portion for stakers, `100 - reward_portion_for_lps - reward_portion_for_stakers = revenue percent for protocol`.
  string reward_portion_for_stakers = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];

  // Maximum eden reward apr for lps - [0 - 0.3]
  string max_eden_reward_apr_lps = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];

  repeated SupportedRewardDenom supported_reward_denoms = 5;

  string protocol_revenue_address = 6;
}
```

### SupportedRewardDenom

The `SupportedRewardDenom` message defines the supported reward denominations and their minimum amounts.

```proto
message SupportedRewardDenom {
    string denom = 1;
    string min_amount = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
}
```

## IncentiveInfo

`IncentiveInfo` holds eden reward distribution info for current year.

```proto
// Incentive Info
message IncentiveInfo {
    // reward amount in eden for 1 year
  string eden_amount_per_year = 1
        [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    // starting block height of the distribution
    string distribution_start_block = 2
        [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    // distribution duration - block number per year
  string total_blocks_per_year = 3
        [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    // blocks distributed
    string blocks_distributed = 4
        [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
}
```

## ExternalIncentive

`ExternalIncentive` holds external incentive record which is created by third party.

```proto
// ExternalIncentive defines the external incentives.
message ExternalIncentive {
  uint64 id = 1;
  string reward_denom = 2;
  uint64 pool_id = 3;
  uint64 from_block = 4;
  uint64 to_block = 5;
  string amount_per_block = 6
    [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
  string apr = 7
    [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
}
```

## PoolInfo

`PoolInfo` holds APR information per reward source, proxy TVL multiplier, and external reward denoms allowed on the pool.

```proto
message PoolInfo {
  uint64 pool_id = 1;
    // reward wallet address
    string reward_wallet = 2;
    // multiplier for lp rewards
    string multiplier = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
      ];
    // Eden APR, updated at every distribution
    string eden_apr = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    // Dex APR, updated at every distribution
    string dex_apr = 5 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    // Gas APR, updated at every distribution
    string gas_apr = 6 [
      (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
      (gogoproto.nullable) = false
    ];
    // External Incentive APR, updated at every distribution
    string external_incentive_apr = 7 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    // external reward denoms on the pool
    repeated string external_reward_denoms = 8;
}
```

## PoolRewardInfo

`PoolRewardInfo` is used to track pool's reward growth per denom.

```proto
message PoolRewardInfo {
    uint64 pool_id = 1;
    string reward_denom = 2;
    string pool_acc_reward_per_share = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
      ];
    uint64 last_updated_block = 4;
}
```

## UserRewardInfo

`UserRewardInfo` is used to track user's reward information per reward denom and pool id.

```proto
message UserRewardInfo {
    string user = 1;
    uint64 pool_id = 2;
    string reward_denom = 3;
    string reward_debt = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
      ];
    string reward_pending = 5 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
      ];
}
```

## PoolRewardsAccum

`PoolRewardsAccum` is introduced to track 24h dex rewards, gas rewards, and eden rewards.

```proto
message PoolRewardsAccum {
  uint64 pool_id = 1;
  int64 block_height = 2;
  uint64 timestamp = 3;
  string dex_reward = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string gas_reward = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string eden_reward = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```
