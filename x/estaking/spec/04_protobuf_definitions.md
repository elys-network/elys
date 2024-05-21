<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### ElysStaked

The `ElysStaked` message tracks the amount of Elys staked by a delegator. This is crucial for managing EdenBoost tokens, which need to be burned when an Elys unstake event occurs.

```proto
message ElysStaked {
    string address = 1;
    string amount = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
}
```

### IncentiveInfo

The `IncentiveInfo` message contains details about the staking incentives, including the reward amount per year, the starting block of the distribution, and the total number of blocks per year.

```proto
message IncentiveInfo {
    string eden_amount_per_year = 1 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    string distribution_start_block = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    string total_blocks_per_year = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    string blocks_distributed = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
}
```

### DexRewardsTracker

The `DexRewardsTracker` message tracks rewards for stakers and liquidity providers. The amounts are denominated in USDC.

```proto
message DexRewardsTracker {
  string num_blocks = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

### Params

The `Params` message defines the parameters for the `estaking` module, including staking incentives, maximum reward APR for stakers, EdenBoost APR, and tracking dex rewards.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;

  IncentiveInfo stake_incentives = 1;
  string eden_commit_val = 2;
  string edenb_commit_val = 3;
  string max_eden_reward_apr_stakers = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string eden_boost_apr = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  DexRewardsTracker dex_rewards_stakers = 7 [
    (gogoproto.nullable) = false
  ];
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `estaking` module at genesis.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated ElysStaked staking_snapshots = 2 [ (gogoproto.nullable) = false ];
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `estaking` module.

```proto
service Msg {
  rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
  rpc WithdrawReward(MsgWithdrawReward) returns (MsgWithdrawRewardResponse);
  rpc WithdrawElysStakingRewards(MsgWithdrawElysStakingRewards) returns (MsgWithdrawElysStakingRewardsResponse);
  rpc WithdrawAllRewards(MsgWithdrawAllRewards) returns (MsgWithdrawAllRewardsResponse);
}
```

#### MsgUpdateParams

This message updates the parameters of the `estaking` module.

```proto
message MsgUpdateParams {
  string authority = 1;
  Params params = 2 [ (gogoproto.nullable) = false ];
}

message MsgUpdateParamsResponse {}
```

#### MsgWithdrawReward

This message allows a delegator to withdraw rewards from a single validator.

```proto
message MsgWithdrawReward {
  string delegator_address = 1;
  string validator_address = 2;
}

message MsgWithdrawRewardResponse {
  repeated cosmos.base.v1beta1.Coin amount = 1 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

#### MsgWithdrawElysStakingRewards

This message allows a delegator to withdraw rewards from all validators.

```proto
message MsgWithdrawElysStakingRewards {
  string delegator_address = 1;
}

message MsgWithdrawElysStakingRewardsResponse {
  repeated cosmos.base.v1beta1.Coin amount = 1 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

#### MsgWithdrawAllRewards

This message allows a delegator to withdraw rewards from all validators and Eden/EdenB commitments.

```proto
message MsgWithdrawAllRewards {
  string delegator_address = 1;
}

message MsgWithdrawAllRewardsResponse {
  repeated cosmos.base.v1beta1.Coin amount = 1 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `estaking` module.

```proto
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/estaking/params";
  }
  rpc Rewards(QueryRewardsRequest) returns (QueryRewardsResponse) {
    option (google.api.http).get = "/elys-network/elys/estaking/rewards/{address}";
  }
}
```

#### QueryParamsRequest

This message requests the parameters of the `estaking` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `estaking` module.

```proto
message QueryParamsResponse {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

#### QueryRewardsRequest

This message requests the total rewards accrued by a delegation.

```proto
message QueryRewardsRequest {
  string address = 1;
}
```

#### QueryRewardsResponse

This message responds with the rewards accrued by a delegator.

```proto
message QueryRewardsResponse {
  repeated DelegationDelegatorReward rewards = 1 [(gogoproto.nullable) = false];
  repeated cosmos.base.v1beta1.Coin total = 2 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

#### DelegationDelegatorReward

This message defines the rewards for a delegator from a specific validator.

```proto
message DelegationDelegatorReward {
  string validator_address = 1;
  repeated cosmos.base.v1beta1.Coin reward = 2 [
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins",
    (gogoproto.nullable) = false
  ];
}
```
