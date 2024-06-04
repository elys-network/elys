<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### DexRewardsTracker

The `DexRewardsTracker` message is used for tracking rewards for liquidity providers in USDC.

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

### ExternalIncentive

The `ExternalIncentive` message defines the external incentives.

```proto
message ExternalIncentive {
    uint64 id = 1;
    string reward_denom = 2;
    uint64 pool_id = 3;
    uint64 from_block = 4;
    uint64 to_block = 5;
    string amount_per_block = 6 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    string apr = 7 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
}
```

### IncentiveInfo

The `IncentiveInfo` message contains details about the LP incentives.

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

### Params

The `Params` message defines the parameters for the `masterchef` module, including LP incentives and reward portions.

```proto
message Params {
    IncentiveInfo lp_incentives = 1;
    string reward_portion_for_lps = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string reward_portion_for_stakers = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    DexRewardsTracker dex_rewards_lps = 4 [(gogoproto.nullable) = false];
    string max_eden_reward_apr_lps = 5 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    repeated SupportedRewardDenom supported_reward_denoms = 6;
    string protocol_revenue_address = 7;
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

### GenesisState

The `GenesisState` message defines the initial state of the `masterchef` module at genesis.

```proto
message GenesisState {
    Params params = 1 [(gogoproto.nullable) = false];
    repeated ExternalIncentive external_incentives = 2 [(gogoproto.nullable) = false];
    uint64 external_incentive_index = 3;
    repeated PoolInfo pool_infos = 4 [(gogoproto.nullable) = false];
    repeated PoolRewardInfo pool_reward_infos = 5 [(gogoproto.nullable) = false];
    repeated UserRewardInfo user_reward_infos = 6 [(gogoproto.nullable) = false];
}
```

### PoolInfo

The `PoolInfo` message contains information about a specific pool.

```proto
message PoolInfo {
    uint64 pool_id = 1;
    string reward_wallet = 2;
    string multiplier = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string num_blocks = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    string dex_reward_amount_given = 5 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string eden_reward_amount_given =

 6 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    string eden_apr = 7 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string dex_apr = 8 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string external_incentive_apr = 9 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    repeated string external_reward_denoms = 10;
}
```

### PoolRewardInfo

The `PoolRewardInfo` message contains reward information for a specific pool.

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

### UserRewardInfo

The `UserRewardInfo` message contains reward information for a specific user.

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

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `masterchef` module.

```proto
service Msg {
    rpc AddExternalRewardDenom(MsgAddExternalRewardDenom) returns (MsgAddExternalRewardDenomResponse);
    rpc AddExternalIncentive(MsgAddExternalIncentive) returns (MsgAddExternalIncentiveResponse);
    rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
    rpc UpdatePoolMultipliers(MsgUpdatePoolMultipliers) returns (MsgUpdatePoolMultipliersResponse);
    rpc ClaimRewards(MsgClaimRewards) returns (MsgClaimRewardsResponse);
}
```

#### MsgAddExternalRewardDenom

This message adds a new external reward denomination to the supported list.

```proto
message MsgAddExternalRewardDenom {
    string authority = 1;
    string reward_denom = 2;
    string min_amount = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
    bool supported = 4;
}

message MsgAddExternalRewardDenomResponse {}
```

#### MsgAddExternalIncentive

This message adds a new external incentive.

```proto
message MsgAddExternalIncentive {
    string sender = 1;
    string reward_denom = 2;
    uint64 pool_id = 3;
    uint64 from_block = 4;
    uint64 to_block = 5;
    string amount_per_block = 6 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
}

message MsgAddExternalIncentiveResponse {}
```

#### MsgUpdateParams

This message updates the parameters of the `masterchef` module.

```proto
message MsgUpdateParams {
    string authority = 1;
    Params params = 2 [
        (gogoproto.nullable) = false
    ];
}

message MsgUpdateParamsResponse {}
```

#### MsgUpdatePoolMultipliers

This message updates the multipliers for the pools.

```proto
message MsgUpdatePoolMultipliers {
    string authority = 1;
    repeated PoolMultiplier pool_multipliers = 2 [
        (gogoproto.nullable) = false
    ];
}

message PoolMultiplier {
    uint64 pool_id = 1;
    string multiplier = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
}

message MsgUpdatePoolMultipliersResponse {}
```

#### MsgClaimRewards

This message allows a user to claim their pending rewards.

```proto
message MsgClaimRewards {
    string sender = 1;
    repeated uint64 pool_ids = 2;
}

message MsgClaimRewardsResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `masterchef` module.

```proto
service Query {
    rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/params";
    }
    rpc ExternalIncentive(QueryExternalIncentiveRequest) returns (QueryExternalIncentiveResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/external-incentive";
    }
    rpc PoolInfo(QueryPoolInfoRequest) returns (QueryPoolInfoResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/pool-info";
    }
    rpc PoolRewardInfo(QueryPoolRewardInfoRequest) returns (QueryPoolRewardInfoResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/pool-reward-info";
    }
    rpc UserRewardInfo(QueryUserRewardInfoRequest) returns (QueryUserRewardInfoResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/user-reward-info";
    }
    rpc UserPendingReward(QueryUserPendingRewardRequest) returns (QueryUserPendingRewardResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/pending-reward";
    }
    rpc StableStakeApr(QueryStableStakeAprRequest) returns (QueryStableStakeAprResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/stable-stake-apr/{denom}";
    }
    rpc PoolAprs(QueryPoolAprsRequest) returns (QueryPoolAprsResponse) {
        option (google.api.http).get = "/elys-network/elys/masterchef/pool-aprs";
    }
}
```

#### QueryParamsRequest

This message requests the parameters of the `masterchef` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `masterchef` module.

```proto
message QueryParamsResponse {
    Params params = 1 [
        (gogoproto.nullable) = false
    ];
}
```

#### QueryExternalIncentiveRequest

This message requests information about an external incentive.

```proto
message QueryExternalIncentiveRequest {
    uint64 id = 1;
}
```

#### QueryExternalIncentiveResponse

This message responds with information about an external incentive.

```proto
message QueryExternalIncentiveResponse {
    ExternalIncentive external_incentive = 1 [
        (gogoproto.nullable) = false
    ];
}
```

#### QueryPoolInfoRequest

This message requests information about a specific pool.

```proto
message QueryPoolInfoRequest {
    uint64 pool_id = 1;
}
```

#### QueryPoolInfoResponse

This message responds with information about a specific pool.

```proto
message QueryPoolInfoResponse {
    PoolInfo pool_info = 1 [
        (gogoproto.nullable) = false
    ];
}
```

#### QueryPoolRewardInfoRequest

This message requests reward information for a specific pool.

```proto
message QueryPoolRewardInfoRequest {
    uint64 pool_id = 1;
    string reward_denom = 2;
}
```

#### QueryPoolRewardInfoResponse

This message responds with reward information for a specific pool.

```proto
message QueryPoolRewardInfoResponse {
    PoolRewardInfo pool_reward_info = 1 [
        (gogoproto.nullable) = false
    ];
}
```

#### QueryUserRewardInfoRequest

This message requests reward information for a specific user.

```proto
message QueryUserRewardInfoRequest {
    string user = 1;
    uint64 pool_id = 2;
    string reward_denom = 3;
}
```

#### QueryUserRewardInfoResponse

This message responds with reward information for a specific user.

```proto
message QueryUserRewardInfoResponse {
    UserRewardInfo user_reward_info = 1 [
        (gogoproto.nullable) = false
    ];
}
```

#### QueryUserPendingRewardRequest

This message requests the pending rewards for a specific user.

```proto
message QueryUserPendingRewardRequest {
    string user = 1;
}
```

#### QueryUserPendingRewardResponse

This message responds with the pending rewards for a specific user.

```proto
message QueryUserPendingRewardResponse {
    repeated RewardInfo rewards = 1;
    repeated cosmos.base.v1beta1.Coin total_rewards = 2 [
        (gogoproto.nullable) = false,
        (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"


    ];
}
```

#### RewardInfo

This message defines the rewards for a user from a specific pool.

```proto
message RewardInfo {
    uint64 pool_id = 1;
    repeated cosmos.base.v1beta1.Coin reward = 2 [
        (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins",
        (gogoproto.nullable) = false
    ];
}
```

#### QueryStableStakeAprRequest

This message requests the stable stake APR for a specific denomination.

```proto
message QueryStableStakeAprRequest {
    string denom = 1;
}
```

#### QueryStableStakeAprResponse

This message responds with the stable stake APR for a specific denomination.

```proto
message QueryStableStakeAprResponse {
    string apr = 1 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.nullable) = false
    ];
}
```

#### QueryPoolAprsRequest

This message requests the APRs for multiple pools.

```proto
message QueryPoolAprsRequest {
    repeated uint64 pool_ids = 1;
}
```

#### QueryPoolAprsResponse

This message responds with the APRs for multiple pools.

```proto
message QueryPoolAprsResponse {
    repeated PoolApr data = 1 [
        (gogoproto.nullable) = false
    ];
}
```

#### PoolApr

This message defines the APR for a specific pool.

```proto
message PoolApr {
    uint64 pool_id = 1;
    string eden_apr = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string usdc_apr = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
    string total_apr = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
}
```
