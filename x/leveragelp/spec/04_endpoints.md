<!--
order: 4
-->

# Endpoints

## Gov Proposals

```proto
service Msg {
    rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
    rpc UpdatePools(MsgUpdatePools) returns (MsgUpdatePoolsResponse);
    rpc Whitelist(MsgWhitelist) returns (MsgWhitelistResponse);
    rpc Dewhitelist(MsgDewhitelist) returns (MsgDewhitelistResponse);
}
```

### UpdateParams

`MsgUpdateParams` is used by governance to update leveragelp module params.

```proto
message MsgUpdateParams {
  string authority = 1;
  Params params = 2 [(gogoproto.nullable) = false];
}
```

### MsgUpdatePools

`MsgUpdatePools` is used by governance to update leveragelp pool infos.

```proto
message MsgUpdatePools {
  string authority = 1;
  repeated Pool pools = 2 [(gogoproto.nullable) = false];
}
```

### MsgWhitelist

`MsgWhitelist` is used by governance to approve an address to the whitelist.

```proto
message MsgWhitelist {
  string authority = 1;
  string whitelisted_address = 2;
}
```

### MsgDewhitelist

`MsgDewhitelist` is used by governance to remove an address from the whitelist.

```proto
message MsgDewhitelist {
  string authority = 1;
  string whitelisted_address = 2;
}
```

## Msgs

The `Msg` service defines the transactions available in the `leveragelp` module.

```proto
service Msg {
    rpc Open(MsgOpen) returns (MsgOpenResponse);
    rpc Close(MsgClose) returns (MsgCloseResponse);
}
```

### MsgOpen

`MsgOpen` is used to open a leveragelp position.

```proto
message MsgOpen {
  string creator = 1;
  string collateral_asset = 2;
  string collateral_amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  uint64 amm_pool_id = 4;
  string leverage = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string stop_loss_price = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

### MsgClose

`MsgClose` is used to partially close a leveragelp position.

```proto
message MsgClose {
  string creator = 1;
  uint64 id = 2;
  string lp_amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

## Query endpoints

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

### Params

Query the parameters of the `leveragelp` module.

### QueryPositions

Query all positions through pagination.

### QueryPositionsByPool

Query open positions by pool id.

### GetStatus

Queries status of leveragelp module, it returns number of active open positions, and positions opened from the beginning.

### QueryPositionsForAddress

Queries open positions for an address

### GetWhitelist

Queries whitelisted addresses for position open, in case `params.whitelisting_enabled` is set to true.

### IsWhitelisted

Queries if an address is in the list of whitelisted array.

### Pool

Queries leveragelp pool info by pool id.

### Pools

Queries all leveragelp pool infos.

### Position

Queries a position by address and id.

### OpenEst

Queries position open estimation result.

### CloseEst

Queries position close estimation result.

## Wasmbindings

### Messages

#### LeveragelpOpen

Connect wasmbinding to `MsgOpen`

#### LeveragelpClose

Connect wasmbinding to `MsgClose`

### Queries

#### LeveragelpParams

Connect wasmbinding to `Params` query.

#### LeveragelpQueryPositions

Connect wasmbinding to `QueryPositions` query.

#### LeveragelpQueryPositionsByPool

Connect wasmbinding to `QueryPositionsByPool` query.

#### LeveragelpGetStatus

Connect wasmbinding to `GetStatus` query.

#### LeveragelpQueryPositionsForAddress

Connect wasmbinding to `QueryPositionsForAddress` query.

#### LeveragelpGetWhitelist

Connect wasmbinding to `GetWhitelist` query.

#### LeveragelpIsWhitelisted

Connect wasmbinding to `IsWhitelisted` query.

#### LeveragelpPool

Connect wasmbinding to `Pool` query.

#### LeveragelpPools

Connect wasmbinding to `Pools` query.

#### LeveragelpPosition

Connect wasmbinding to `Position` query.

#### LeveragelpOpenEst

Connect wasmbinding to `OpenEst` query.

#### LeveragelpCloseEst

Connect wasmbinding to `CloseEst` query.
