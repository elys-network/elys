<!--
order: 4
-->

# Protobuf Definitions

## Types

### Params

The `Params` message defines the parameters for the `transferhook` module, including whether the AMM functionality is active.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;
  bool amm_active = 1;
}
```

## Messages

### Query Service

The `Query` service defines the gRPC querier service for the `transferhook` module.

```proto
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/transferhook/params";
  }
}
```

#### QueryParamsRequest

This message requests the parameters of the `transferhook` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `transferhook` module.

```proto
message QueryParamsResponse {
  Params params = 1 [ (gogoproto.nullable) = false ];
}
```

## States

### GenesisState

The `GenesisState` message defines the initial state of the `transferhook` module at genesis.

```proto
message GenesisState {
  Params params = 1 [ (gogoproto.nullable) = false ];
}
```
