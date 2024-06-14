<!--
order: 4
-->

# Endpoints

## Gov Proposals

```proto
service Msg {
    rpc AddExternalRewardDenom(MsgAddExternalRewardDenom) returns (MsgAddExternalRewardDenomResponse);
    rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
    rpc UpdatePoolMultipliers(MsgUpdatePoolMultipliers) returns (MsgUpdatePoolMultipliersResponse);
}
```

### UpdateParams

`MsgUpdateParams` is used by governance to update masterchef module params.

```proto
message MsgUpdateParams {
  string authority = 1;
  Params params = 2 [(gogoproto.nullable) = false];
}
```

### UpdatePoolMultipliers

`MsgUpdatePoolMultipliers` is used by governance to update masterchef pool info multipliers.

```proto
message PoolMultiplier {
  uint64 pool_id = 1;
  string multiplier = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}

message MsgUpdatePoolMultipliers {
  string authority = 1;
  repeated PoolMultiplier pool_multipliers = 2 [ (gogoproto.nullable) = false ];
}
```

### AddExternalRewardDenom

`MsgAddExternalRewardDenom` is used by governance to approve a reward denom.

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
```

## Msgs

The `Msg` service defines the transactions available in the `masterchef` module.

```proto
service Msg {
    rpc AddExternalIncentive(MsgAddExternalIncentive) returns (MsgAddExternalIncentiveResponse);
    rpc ClaimRewards(MsgClaimRewards) returns (MsgClaimRewardsResponse);
}
```

### MsgAddExternalIncentive

`MsgAddExternalIncentive` is used to add external incentive on a specific pool from `from_block` to `to_block` with same amount of `amount_per_block` per block for `reward_denom` denom.

```proto
message MsgAddExternalIncentive {
    string sender = 1;
    string reward_denom = 2;
    uint64 pool_id = 3;
    uint64 from_block = 4;
    uint64 to_block = 5;
    string amount_per_block = 6
        [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
}
```

### MsgClaimRewards

`MsgClaimRewards` is used to claim pending rewards on specified pool ids (`pool_ids`).

```proto
message MsgClaimRewards {
    string sender = 1;
    repeated uint64 pool_ids = 2;
}
```

## Query endpoints

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

### Params

Query module params

### ExternalIncentive

Query external incentive by id.

### PoolInfo

Query pool info by pool id

### PoolRewardInfo

Query pool reward info by pool id and reward denom.

### UserRewardInfo

Query user reward info by user address, pool id and reward denom.

### UserPendingReward

Query user pending rewards per pool from user address.

### StableStakeApr

Query stablestake APR by reward denom.

### PoolAprs

Query pool APRs from specified pool ids. If nothing's put in pool ids, full list of pool APRs returned.

## Wasmbindings

### Messages

#### MastercheClaimRewards

Connect wasmbinding to `MsgClaimRewards`

### Queries

#### MasterchefParams

Connect wasmbinding to `Params` query.

#### MasterchefExternalIncentive

Connect wasmbinding to `ExternalIncentive` query.

#### MasterchefPoolInfo

Connect wasmbinding to `PoolInfo` query.

#### MasterchefPoolRewardInfo

Connect wasmbinding to `PoolRewardInfo` query.

#### MasterchefUserRewardInfo

Connect wasmbinding to `UserRewardInfo` query.

#### MasterchefUserPendingReward

Connect wasmbinding to `UserPendingReward` query.

#### MasterchefStableStakeApr

Connect wasmbinding to `StableStakeApr` query.

#### MasterchefPoolAprs

Connect wasmbinding to `PoolAprs` query.
