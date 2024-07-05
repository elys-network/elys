<!--
order: 4
-->

# Protobuf Definitions

## Types

### Params

The `Params` message defines the structure for module parameters, including contract addresses and gas limits.

```proto
message Params {
  repeated string contract_addresses = 1;
  uint64 contract_gas_limit = 2;
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `clock` module at genesis.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

## Messages

### MsgUpdateParams

The `MsgUpdateParams` message allows updating the module parameters via governance.

```proto
message MsgUpdateParams {
  string authority = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  Params params = 2 [(gogoproto.nullable) = false];
}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `clock` module.

```proto
service Query {
  rpc ClockContracts (QueryClockContracts) returns (QueryClockContractsResponse) {
    option (google.api.http).get = "/elys/clock/v1/contracts";
  }
  rpc Params (QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys/clock/v1/params";
  }
}
```

#### QueryClockContractsRequest

This message requests a list of all contract addresses.

```proto
message QueryClockContracts {}
```

#### QueryClockContractsResponse

This message responds with the list of all contract addresses.

```proto
message QueryClockContractsResponse {
  repeated string contract_addresses = 1;
}
```

#### QueryParamsRequest

This message requests the current module parameters.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the current module parameters.

```proto
message QueryParamsResponse {
  Params params = 1;
}
```
