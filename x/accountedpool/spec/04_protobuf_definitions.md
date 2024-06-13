<!--
order: 4
-->

# Protobuf Definitions

## Types

### AccountedPool

The `AccountedPool` message defines the structure of an accounted pool, including its ID, total shares, pool assets, and total weight.

```proto
message AccountedPool {
  uint64 pool_id = 1;
  cosmos.base.v1beta1.Coin total_shares = 2 [(gogoproto.nullable) = false];
  repeated elys.amm.PoolAsset pool_assets = 3 [(gogoproto.nullable) = false];
  string total_weight = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `accountedpool` module at genesis.

```proto
message GenesisState {
  repeated AccountedPool accounted_pool_list = 1 [(gogoproto.nullable) = false];
}
```

## Messages

The `accountedpool` module does not define any specific messages for transactions.

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `accountedpool` module.

```proto
service Query {
  rpc AccountedPool (QueryGetAccountedPoolRequest) returns (QueryGetAccountedPoolResponse) {
    option (google.api.http).get = "/elys-network/elys/accountedpool/accounted_pool/{pool_id}";
  }
  rpc AccountedPoolAll (QueryAllAccountedPoolRequest) returns (QueryAllAccountedPoolResponse) {
    option (google.api.http).get = "/elys-network/elys/accountedpool/accounted_pool";
  }
}
```

#### QueryGetAccountedPoolRequest

This message requests a specific accounted pool by its ID.

```proto
message QueryGetAccountedPoolRequest {
  uint64 pool_id = 1;
}
```

#### QueryGetAccountedPoolResponse

This message responds with the details of a specific accounted pool.

```proto
message QueryGetAccountedPoolResponse {
  AccountedPool accounted_pool = 1 [(gogoproto.nullable) = false];
}
```

#### QueryAllAccountedPoolRequest

This message requests a list of all accounted pools with pagination support.

```proto
message QueryAllAccountedPoolRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllAccountedPoolResponse

This message responds with a list of all accounted pools and pagination details.

```proto
message QueryAllAccountedPoolResponse {
  repeated AccountedPool accounted_pool = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```
