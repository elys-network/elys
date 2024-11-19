<!--
order: 4
-->

# Protobuf Definitions

## Types

### Params

The `Params` message defines the structure of the configuration parameters, including commission rates, voting power, self-delegation and broker address.

```proto
message Params {
  string min_commission_rate = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string max_voting_power = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string min_self_delegation = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string broker_address = 4;
  int64 total_blocks_per_year = 5;
  int64 rewards_data_lifetime = 6;
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `parameter` module at genesis.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

## Messages

### MsgUpdateMinCommission

This message updates the minimum commission rate.

```proto
message MsgUpdateMinCommission {
  string creator = 1;
  string min_commission = 2;
}
message MsgUpdateMinCommissionResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `parameter` module.

```proto
service Query {
  rpc Params (QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/parameter/params";
  }
}
```

#### QueryParamsRequest

This message requests the current parameters.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the current parameters.

```proto
message QueryParamsResponse {
  Params params = 1 [(gogoproto.nullable) = false];
}
```
