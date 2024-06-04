<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### Position

The `Position` message tracks leveraged positions in the module. It includes information on the collateral, leverage, health, and other parameters.

```proto
message Position {
    string address = 1;
    cosmos.base.v1beta1.Coin collateral = 2 [(gogoproto.nullable) = false];
    string liabilities = 3 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    string interest_paid = 4 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    string leverage = 5 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    string leveraged_lp_amount = 6 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    string position_health = 7 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    uint64 id = 8;
    uint64 amm_pool_id = 9;
    string stop_loss_price = 10 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
}
```

### Pool

The `Pool` message contains information about the AMM pools involved in leveraged positions, including their health and the total amount of leveraged LP tokens.

```proto
message Pool {
    uint64 amm_pool_id = 1;
    string health = 2 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    bool enabled = 3;
    bool closed = 4;
    string leveraged_lp_amount = 5 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    string leverage_max = 6 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
}
```

### Params

The `Params` message defines the module's parameters, including leverage limits, whitelisting options, and other configurations.

```proto
message Params {
    option (gogoproto.goproto_stringer) = false;
    string leverage_max = 1 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    int64 max_open_positions = 2;
    string pool_open_threshold = 3 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    string safety_factor = 4 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    bool whitelisting_enabled = 5;
    int64 epoch_length = 6;
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `leveragelp` module at genesis.

```proto
message GenesisState {
    Params params = 1 [(gogoproto.nullable) = false];
    repeated Pool pool_list = 2 [(gogoproto.nullable) = false];
    repeated Position position_list = 3 [(gogoproto.nullable) = false];
    repeated string address_whitelist = 4;
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `leveragelp` module.

```proto
service Msg {
    rpc Open(MsgOpen) returns (MsgOpenResponse);
    rpc Close(MsgClose) returns (MsgCloseResponse);
    rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
    rpc UpdatePools(MsgUpdatePools) returns (MsgUpdatePoolsResponse);
    rpc Whitelist(MsgWhitelist) returns (MsgWhitelistResponse);
    rpc Dewhitelist(MsgDewhitelist) returns (MsgDewhitelistResponse);
}
```

#### MsgOpen

This message opens a leveraged position.

```proto
message MsgOpen {
    string creator = 1;
    string collateral_asset = 2;
    string collateral_amount = 3 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    uint64 amm_pool_id = 4;
    string leverage = 5 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
    string stop_loss_price = 6 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec", (gogoproto.nullable) = false];
}

message MsgOpenResponse {}
```

#### MsgClose

This message closes a leveraged position.

```proto
message MsgClose {
    string creator = 1;
    uint64 id = 2;
    string lp_amount = 3 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
}

message MsgCloseResponse {}
```

#### MsgUpdateParams

This message updates the parameters of the `leveragelp` module.

```proto
message MsgUpdateParams {
    string authority = 1;
    Params params = 2 [(gogoproto.nullable) = false];
}

message MsgUpdateParamsResponse {}
```

#### MsgUpdatePools

This message updates the pools managed by the `leveragelp` module.

```proto
message MsgUpdatePools {
    string authority = 1;
    repeated Pool pools = 2 [(gogoproto.nullable)

 = false];
}

message MsgUpdatePoolsResponse {}
```

#### MsgWhitelist

This message adds an address to the whitelist.

```proto
message MsgWhitelist {
    string authority = 1;
    string whitelisted_address = 2;
}

message MsgWhitelistResponse {}
```

#### MsgDewhitelist

This message removes an address from the whitelist.

```proto
message MsgDewhitelist {
    string authority = 1;
    string whitelisted_address = 2;
}

message MsgDewhitelistResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `leveragelp` module.

```proto
service Query {
    rpc Params(ParamsRequest) returns (ParamsResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/params";
    }
    rpc QueryPositions(PositionsRequest) returns (PositionsResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/positions/{pagination.key}";
    }
    rpc QueryPositionsByPool(PositionsByPoolRequest) returns (PositionsByPoolResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/positions-by-pool/{amm_pool_id}/{pagination.key}";
    }
    rpc GetStatus(StatusRequest) returns (StatusResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/status";
    }
    rpc QueryPositionsForAddress(PositionsForAddressRequest) returns (PositionsForAddressResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/positions-for-address/{address}/{pagination.key}";
    }
    rpc GetWhitelist(WhitelistRequest) returns (WhitelistResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/whitelist/{pagination.key}";
    }
    rpc IsWhitelisted(IsWhitelistedRequest) returns (IsWhitelistedResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/is-whitelisted";
    }
    rpc Pool(QueryGetPoolRequest) returns (QueryGetPoolResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/pool/{index}";
    }
    rpc Pools(QueryAllPoolRequest) returns (QueryAllPoolResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/pool/{pagination.key}";
    }
    rpc Position(PositionRequest) returns (PositionResponse) {
        option (google.api.http).get = "/elys-network/elys/leveragelp/position/{address}/{id}";
    }
}
```

#### ParamsRequest

This message requests the parameters of the `leveragelp` module.

```proto
message ParamsRequest {}
```

#### ParamsResponse

This message responds with the parameters of the `leveragelp` module.

```proto
message ParamsResponse {
    Params params = 1 [(gogoproto.nullable) = false];
}
```

#### PositionsRequest

This message requests the list of all positions.

```proto
message PositionsRequest {
    cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### PositionsResponse

This message responds with the list of all positions.

```proto
message PositionsResponse {
    repeated Position positions = 1;
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### PositionsByPoolRequest

This message requests the list of positions for a specific pool.

```proto
message PositionsByPoolRequest {
    uint64 amm_pool_id = 1;
    cosmos.base.query.v1beta1.PageRequest pagination = 2;
}
```

#### PositionsByPoolResponse

This message responds with the list of positions for a specific pool.

```proto
message PositionsByPoolResponse {
    repeated Position positions = 1;
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### StatusRequest

This message requests the status of the module.

```proto
message StatusRequest {}
```

#### StatusResponse

This message responds with the status of the module.

```proto
message StatusResponse {
    uint64 open_position_count = 1;
    uint64 lifetime_position_count = 2;
}
```

#### PositionsForAddressRequest

This message requests the list of positions for a specific address.

```proto
message PositionsForAddressRequest {
    string address = 1;
    cosmos.base.query.v1beta1.PageRequest pagination = 2;
}
```

#### PositionsForAddressResponse

This message responds with the list of positions for a specific address.

```proto
message PositionsForAddressResponse {
    repeated Position positions = 1;
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### WhitelistRequest

This message requests the whitelist.

```proto
message WhitelistRequest {
    cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### WhitelistResponse

This message responds with the whitelist.

```proto
message WhitelistResponse {
    repeated string whitelist = 1;
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### IsWhitelistedRequest

This message checks if an address is whitelisted.

```proto
message IsWhitelistedRequest {
    string address = 1;
}
```

#### IsWhitelistedResponse

This message responds with the whitelisted status of an address.

```proto
message IsWhitelistedResponse {
    string address = 1;
    bool is_whitelisted = 2;
}
```

#### QueryGetPoolRequest

This message requests a specific pool by its ID.

```proto
message QueryGetPoolRequest {
    uint64 index = 1;
}
```

#### QueryGetPoolResponse

This message responds with a specific pool.

```proto
message QueryGetPoolResponse {
    Pool pool = 1 [(gogoproto.nullable) = false];
}
```

#### QueryAllPoolRequest

This message requests all pools.

```proto
message QueryAllPoolRequest {
    cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllPoolResponse

This message responds with all pools.

```proto
message QueryAllPoolResponse {
    repeated Pool pool = 1 [(gogoproto.nullable) = false];
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### PositionRequest

This message requests a specific position by address and ID.

```proto
message PositionRequest {
    string address = 1;
    uint64 id = 2;
}
```

#### PositionResponse

This message responds with a specific position.

```proto
message PositionResponse {
    Position position = 1;
}
```
