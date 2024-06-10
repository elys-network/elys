<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### Commitments

The `Commitments` message represents the state of a user's committed tokens, including any vesting schedules.

```proto
message Commitments {
    string creator = 1;
    repeated CommittedTokens committed_tokens = 2;
    repeated cosmos.base.v1beta1.Coin claimed = 3;
    repeated VestingTokens vesting_tokens = 4;
}
```

### Lockup

The `Lockup` message defines a locked token amount with a specific unlock timestamp.

```proto
message Lockup {
    string amount = 1;
    uint64 unlock_timestamp = 2;
}
```

### CommittedTokens

The `CommittedTokens` message holds information about a specific denomination of committed tokens and their lockups.

```proto
message CommittedTokens {
    string denom = 1;
    string amount = 2;
    repeated Lockup lockups = 3;
}
```

### VestingTokens

The `VestingTokens` message represents tokens that are being vested over a certain number of blocks.

```proto
message VestingTokens {
    string denom = 1;
    string total_amount = 2;
    string claimed_amount = 3;
    int64 num_blocks = 5;
    int64 start_block = 6;
    int64 vest_started_timestamp = 7;
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `commitment` module.

```proto
service Msg {
  rpc CommitClaimedRewards(MsgCommitClaimedRewards) returns (MsgCommitClaimedRewardsResponse);
  rpc UncommitTokens(MsgUncommitTokens) returns (MsgUncommitTokensResponse);
  rpc Vest(MsgVest) returns (MsgVestResponse);
  rpc VestNow(MsgVestNow) returns (MsgVestNowResponse);
  rpc VestLiquid(MsgVestLiquid) returns (MsgVestLiquidResponse);
  rpc CancelVest(MsgCancelVest) returns (MsgCancelVestResponse);
  rpc ClaimVesting(MsgClaimVesting) returns (MsgClaimVestingResponse);
  rpc UpdateVestingInfo(MsgUpdateVestingInfo) returns (MsgUpdateVestingInfoResponse);
  rpc Stake(MsgStake) returns (MsgStakeResponse);
  rpc Unstake(MsgUnstake) returns (MsgUnstakeResponse);
}
```

#### MsgCommitClaimedRewards

This message commits the claimed rewards to the committed store.

```proto
message MsgCommitClaimedRewards {
  string creator = 1;
  string amount  = 2;
  string denom   = 3;
}

message MsgCommitClaimedRewardsResponse {}
```

#### MsgUncommitTokens

This message uncommits tokens from the committed store, making them immediately liquid.

```proto
message MsgUncommitTokens {
  string creator = 1;
  string amount  = 2;
  string denom   = 3;
}

message MsgUncommitTokensResponse {}
```

#### MsgVest

This message converts user's committed tokens to vesting tokens.

```proto
message MsgVest {
  string creator = 1;
  string amount  = 2;
  string denom   = 3;
}

message MsgVestResponse {}
```

#### MsgClaimVesting

This message claims already vested tokens.

```proto
message MsgClaimVesting {
  string sender = 1;
}

message MsgClaimVestingResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `commitment` module.

```proto
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse);
  rpc ShowCommitments(QueryShowCommitmentsRequest) returns (QueryShowCommitmentsResponse);
  rpc NumberOfCommitments(QueryNumberOfCommitmentsRequest) returns (QueryNumberOfCommitmentsResponse);
}
```

#### QueryParamsRequest

This message requests the parameters of the `commitment` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `commitment` module.

```proto
message QueryParamsResponse {
  Params params = 1;
}
```

#### QueryShowCommitmentsRequest

This message requests the commitments for a specific creator.

```proto
message QueryShowCommitmentsRequest {
  string creator = 1;
}
```

#### QueryShowCommitmentsResponse

This message responds with the commitments for a specific creator.

```proto
message QueryShowCommitmentsResponse {
  Commitments commitments = 1;
}
```

#### QueryNumberOfCommitmentsRequest

This message requests the total number of commitment items.

```proto
message QueryNumberOfCommitmentsRequest {}
```

#### QueryNumberOfCommitmentsResponse

This message responds with the total number of commitment items.

```proto
message QueryNumberOfCommitmentsResponse {
  int64 number = 1;
}
```
