<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### Airdrop

The `Airdrop` message defines the structure of an airdrop entry.

```proto
message Airdrop {
  string intent = 1;
  uint64 amount = 2;
  string authority = 3;
  uint64 expiry = 4;
}
```

### GenesisInflation

The `GenesisInflation` message defines the parameters for genesis inflation.

```proto
message GenesisInflation {
  InflationEntry inflation = 1;
  uint64 seed_vesting = 2;
  uint64 strategic_sales_vesting = 3;
  string authority = 4;
}
```

### InflationEntry

The `InflationEntry` message defines the structure of an inflation entry.

```proto
message InflationEntry {
  uint64 lm_rewards = 1;
  uint64 ics_staking_rewards = 2;
  uint64 community_fund = 3;
  uint64 strategic_reserve = 4;
  uint64 team_tokens_vested = 5;
}
```

### Params

The `Params` message defines the parameters for the `tokenomics` module.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;
}
```

### TimeBasedInflation

The `TimeBasedInflation` message defines the structure for time-based inflation.

```proto
message TimeBasedInflation {
  uint64 start_block_height = 1;
  uint64 end_block_height = 2;
  string description = 3;
  InflationEntry inflation = 4;
  string authority = 5;
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `tokenomics` module at genesis.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated Airdrop airdrop_list = 2 [(gogoproto.nullable) = false];
  GenesisInflation genesis_inflation = 3;
  repeated TimeBasedInflation time_based_inflation_list = 4 [(gogoproto.nullable) = false];
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `tokenomics` module.

```proto
service Msg {
  rpc CreateAirdrop (MsgCreateAirdrop) returns (MsgCreateAirdropResponse);
  rpc UpdateAirdrop (MsgUpdateAirdrop) returns (MsgUpdateAirdropResponse);
  rpc DeleteAirdrop (MsgDeleteAirdrop) returns (MsgDeleteAirdropResponse);
  rpc ClaimAirdrop (MsgClaimAirdrop) returns (MsgClaimAirdropResponse);
  rpc UpdateGenesisInflation (MsgUpdateGenesisInflation) returns (MsgUpdateGenesisInflationResponse);
  rpc CreateTimeBasedInflation (MsgCreateTimeBasedInflation) returns (MsgCreateTimeBasedInflationResponse);
  rpc UpdateTimeBasedInflation (MsgUpdateTimeBasedInflation) returns (MsgUpdateTimeBasedInflationResponse);
  rpc DeleteTimeBasedInflation (MsgDeleteTimeBasedInflation) returns (MsgDeleteTimeBasedInflationResponse);
}
```

#### MsgCreateAirdrop

This message creates an airdrop entry.

```proto
message MsgCreateAirdrop {
  string authority = 1;
  string intent = 2;
  uint64 amount = 3;
  uint64 expiry = 4;
}

message MsgCreateAirdropResponse {}
```

#### MsgUpdateAirdrop

This message updates an airdrop entry.

```proto
message MsgUpdateAirdrop {
  string authority = 1;
  string intent = 2;
  uint64 amount = 3;
  uint64 expiry = 4;
}

message MsgUpdateAirdropResponse {}
```

#### MsgDeleteAirdrop

This message deletes an airdrop entry.

```proto
message MsgDeleteAirdrop {
  string authority = 1;
  string intent = 2;
}

message MsgDeleteAirdropResponse {}
```

#### MsgClaimAirdrop

This message allows a user to claim an airdrop.

```proto
message MsgClaimAirdrop {
  string sender = 1;
}

message MsgClaimAirdropResponse {}
```

#### MsgUpdateGenesisInflation

This message updates the genesis inflation parameters.

```proto
message MsgUpdateGenesisInflation {
  string authority = 1;
  InflationEntry inflation = 3;
  uint64 seed_vesting = 4;
  uint64 strategic_sales_vesting = 5;
}

message MsgUpdateGenesisInflationResponse {}
```

#### MsgCreateTimeBasedInflation

This message creates a time-based inflation entry.

```proto
message MsgCreateTimeBasedInflation {
  string authority = 1;
  uint64 start_block_height = 2;
  uint64 end_block_height = 3;
  string description = 4;
  InflationEntry inflation = 5;
}

message MsgCreateTimeBasedInflationResponse {}
```

#### MsgUpdateTimeBasedInflation

This message updates a time-based inflation entry.

```proto
message MsgUpdateTimeBasedInflation {
  string authority = 1;
  uint64 start_block_height = 2;
  uint64 end_block_height = 3;
  string description = 4;
  InflationEntry inflation = 5;
}

message MsgUpdateTimeBasedInflationResponse {}
```

#### MsgDeleteTimeBasedInflation

This message deletes a time-based inflation entry.

```proto
message MsgDeleteTimeBasedInflation {
  string authority = 1;
  uint64 start_block_height = 2;
  uint64 end block_height = 3;
}

message MsgDeleteTimeBasedInflationResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `tokenomics` module.

```proto
service Query {
  rpc Params (QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/tokenomics/params";
  }
  rpc Airdrop (QueryGetAirdropRequest) returns (QueryGetAirdropResponse) {
    option (google.api.http).get = "/elys-network/elys/tokenomics/airdrop/{intent}";
  }
  rpc AirdropAll (QueryAllAirdropRequest) returns (QueryAllAirdropResponse) {
    option (google.api.http).get = "/elys-network/elys/tokenomics/airdrop";
  }
  rpc GenesisInflation (QueryGetGenesisInflationRequest) returns (QueryGetGenesisInflationResponse) {
    option (google.api.http).get = "/elys-network/elys/tokenomics/genesis_inflation";
  }
  rpc TimeBasedInflation (QueryGetTimeBasedInflationRequest) returns (QueryGetTimeBasedInflationResponse) {
    option (google.api.http).get = "/elys-network/elys/tokenomics/time_based_inflation/{start_block_height}/{end_block_height}";
  }
  rpc TimeBasedInflationAll (QueryAllTimeBasedInflationRequest) returns (QueryAllTimeBasedInflationResponse) {
    option (google.api.http).get = "/elys-network/elys/tokenomics/time_based_inflation";
  }
}
```

#### QueryParamsRequest

This message requests the parameters of the `tokenomics` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `tokenomics` module.

```proto
message QueryParamsResponse {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

#### QueryGetAirdropRequest

This message requests a specific airdrop entry.

```proto
message QueryGetAirdropRequest {
  string intent = 1;
}
```

#### QueryGetAirdrop

Response

This message responds with the details of a specific airdrop entry.

```proto
message QueryGetAirdropResponse {
  Airdrop airdrop = 1 [(gogoproto.nullable) = false];
}
```

#### QueryAllAirdropRequest

This message requests a list of all airdrop entries.

```proto
message QueryAllAirdropRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllAirdropResponse

This message responds with a list of all airdrop entries.

```proto
message QueryAllAirdropResponse {
  repeated Airdrop airdrop = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### QueryGetGenesisInflationRequest

This message requests the genesis inflation parameters.

```proto
message QueryGetGenesisInflationRequest {}
```

#### QueryGetGenesisInflationResponse

This message responds with the genesis inflation parameters.

```proto
message QueryGetGenesisInflationResponse {
  GenesisInflation genesis_inflation = 1 [(gogoproto.nullable) = false];
}
```

#### QueryGetTimeBasedInflationRequest

This message requests specific time-based inflation parameters.

```proto
message QueryGetTimeBasedInflationRequest {
  uint64 start_block_height = 1;
  uint64 end_block_height = 2;
}
```

#### QueryGetTimeBasedInflationResponse

This message responds with specific time-based inflation parameters.

```proto
message QueryGetTimeBasedInflationResponse {
  TimeBasedInflation time_based_inflation = 1 [(gogoproto.nullable) = false];
}
```

#### QueryAllTimeBasedInflationRequest

This message requests a list of all time-based inflation entries.

```proto
message QueryAllTimeBasedInflationRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllTimeBasedInflationResponse

This message responds with a list of all time-based inflation entries.

```proto
message QueryAllTimeBasedInflationResponse {
  repeated TimeBasedInflation time_based_inflation = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```
