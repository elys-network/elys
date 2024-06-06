<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### AssetInfo

The `AssetInfo` message contains information about an asset, including its denomination, display name, Band ticker, Elys ticker, and decimal precision.

```proto
message AssetInfo {
  string denom = 1;
  string display = 2;
  string band_ticker = 3;
  string elys_ticker = 4;
  uint64 decimal = 5;
}
```

### BandPriceCallData

The `BandPriceCallData` message contains the symbols to query and the multiplier used in Band price queries.

```proto
message BandPriceCallData {
  repeated string symbols = 1;
  uint64 multiplier = 2;
}
```

### BandPriceResult

The `BandPriceResult` message contains the results of a Band price query.

```proto
message BandPriceResult {
  repeated uint64 rates = 1;
}
```

### Params

The `Params` message defines the parameters for the `oracle` module, including settings for Band channel source, Oracle script ID, multiplier, ask count, minimum count, fee limit, gas parameters, client ID, band epoch, price expiry time, and lifetime in blocks.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;

  string band_channel_source = 1;
  uint64 oracle_script_id = 2 [
    (gogoproto.customname) = "OracleScriptID"
  ];
  uint64 multiplier = 3;
  uint64 ask_count = 4;
  uint64 min_count = 5;
  repeated cosmos.base.v1beta1.Coin fee_limit = 6 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
  uint64 prepare_gas = 7;
  uint64 execute_gas = 8;
  string client_id = 9 [
    (gogoproto.customname) = "ClientID"
  ];
  string band_epoch = 10;
  uint64 price_expiry_time = 11;
  uint64 life_time_in_blocks = 12;
}
```

### Price

The `Price` message contains price data for an asset, including the asset identifier, price value, source, provider, timestamp, and block height.

```proto
message Price {
  string asset = 1;
  string price = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string source = 3;
  string provider = 4;
  uint64 timestamp = 5;
  uint64 block_height = 6;
}
```

### PriceFeeder

The `PriceFeeder` message contains information about a price feeder, including the feeder identifier and whether it is active.

```proto
message PriceFeeder {
  string feeder = 1;
  bool is_active = 2;
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `oracle` module at genesis.

```proto
message GenesisState {
  Params params = 1 [
    (gogoproto.nullable) = false
  ];
  string port_id = 2;
  repeated AssetInfo asset_infos = 3 [
    (gogoproto.nullable) = false
  ];
  repeated Price prices = 4 [
    (gogoproto.nullable) = false
  ];
  repeated PriceFeeder price_feeders = 5 [
    (gogoproto.nullable) = false
  ];
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `oracle` module.

```proto
service Msg {
  rpc FeedPrice(MsgFeedPrice) returns (MsgFeedPriceResponse);
  rpc FeedMultiplePrices(MsgFeedMultiplePrices) returns (MsgFeedMultiplePricesResponse);
  rpc SetPriceFeeder(MsgSetPriceFeeder) returns (MsgSetPriceFeederResponse);
  rpc DeletePriceFeeder(MsgDeletePriceFeeder) returns (MsgDeletePriceFeederResponse);
  // proposals
  rpc AddAssetInfo(MsgAddAssetInfo) returns (MsgAddAssetInfoResponse);
  rpc RemoveAssetInfo(MsgRemoveAssetInfo) returns (MsgRemoveAssetInfoResponse);
  rpc AddPriceFeeders(MsgAddPriceFeeders) returns (MsgAddPriceFeedersResponse);
  rpc RemovePriceFeeders(MsgRemovePriceFeeders) returns (MsgRemovePriceFeedersResponse);
  rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
}
```

#### MsgFeedPrice

This message allows a price feeder to submit a price for an asset.

```proto
message MsgFeedPrice {
  string asset = 1;
  string price = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string source = 3;
  string provider = 4;
}

message MsgFeedPriceResponse {}
```

#### MsgFeedMultiplePrices

This message allows a price feeder to submit multiple prices at once.

```proto
message MsgFeedMultiplePrices {
  string creator = 1;
  repeated Price prices = 2 [
    (gogoproto.nullable) = false
  ];
}

message MsgFeedMultiplePricesResponse {}
```

#### MsgSetPriceFeeder

This message allows the authority to set the status of a price feeder.

```proto
message MsgSetPriceFeeder {
  string feeder = 1;
  bool is_active = 2;
}

message MsgSetPriceFeederResponse {}
```

#### MsgDeletePriceFeeder

This message allows the authority to delete a price feeder.

```proto
message MsgDeletePriceFeeder {
  string feeder = 1;
}

message MsgDeletePriceFeederResponse {}
```

#### MsgAddAssetInfo

This message allows the authority to add information about an asset.

```proto
message MsgAddAssetInfo {
  string authority = 1;
  string denom = 2;
  string display = 3;
  string band_ticker = 4;
  string elys_ticker = 5;
  uint64 decimal = 6;
}

message MsgAddAssetInfoResponse {}
```

#### MsgRemoveAssetInfo

This message allows the authority to remove information about an asset.

```proto
message MsgRemoveAssetInfo {
  string authority = 1;
  string denom = 2;
}

message MsgRemoveAssetInfoResponse {}
```

#### MsgAddPriceFeeders

This message allows the authority to add multiple price feeders.

```proto
message MsgAddPriceFeeders {
  string authority = 1;
  repeated string feeders = 2;
}

message MsgAddPriceFeedersResponse {}
```

#### MsgRemovePriceFeeders

This message allows the authority to remove multiple price feeders.

```proto
message MsgRemovePriceFeeders {
  string authority = 1;
  repeated string feeders = 2;
}

message MsgRemovePriceFeedersResponse {}
```

#### MsgUpdateParams

This message allows the authority to update the parameters of the `oracle` module.

```proto
message MsgUpdateParams {
  string authority = 1;
  Params params = 2 [
    (gogoproto.nullable) = false
  ];
}

message MsgUpdateParamsResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `oracle` module.

```proto
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/params";
  }
  rpc BandPriceResult(QueryBandPriceRequest) returns (QueryBandPriceResponse) {
    option (google.api.http).get = "/elys/oracle/band_price_result/{request_id}";
  }
  rpc LastBandRequestId(QueryLastBandRequestIdRequest) returns (QueryLastBandRequestIdResponse) {
    option (google.api.http).get = "/elys/oracle/last_band_price_request_id";
  }
  rpc AssetInfo(QueryGetAssetInfoRequest) returns (QueryGetAssetInfoResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/asset_info/{denom}";
  }
  rpc AssetInfoAll(QueryAllAssetInfoRequest) returns (QueryAllAssetInfoResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/asset_info";
  }
  rpc Price(QueryGetPriceRequest) returns (QueryGetPriceResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/price";
  }
  rpc PriceAll(QueryAllPriceRequest) returns (QueryAllPriceResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/price";
  }
  rpc PriceFeeder(QueryGetPriceFeederRequest) returns (QueryGetPriceFeederResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/price_feeder/{feeder}";
  }
  rpc PriceFeederAll(QueryAllPriceFeederRequest) returns (QueryAllPriceFeederResponse) {
    option (google.api.http).get = "/elys-network/elys/oracle/price_feeder";
  }
}
```

#### QueryParamsRequest

This message requests the parameters of the `oracle` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `oracle` module.

```proto
message QueryParamsResponse {
  Params params = 1 [
    (gogoproto.nullable) = false
  ];
}
```

#### QueryBandPriceRequest

This message requests the Band price result for a given request ID.

```proto
message QueryBandPriceRequest {
  int64 request_id = 1;
}
```

#### QueryBandPriceResponse

This message responds with the Band price result for a given request ID.

```proto
message QueryBandPriceResponse {
  BandPriceResult result = 1;
}
```

#### QueryLastBandRequestIdRequest

This message requests the last Band price request ID.

```proto
message QueryLastBandRequestIdRequest {}
```

#### QueryLastBandRequestIdResponse

This message responds with the last Band price request ID.

```proto
message QueryLastBandRequestIdResponse {
  int64 request_id = 1;
}
```

#### QueryGetAssetInfoRequest

This message requests the information of an asset by its denomination.

```proto
message QueryGetAssetInfoRequest {
  string denom = 1;
}
```

#### QueryGetAssetInfoResponse

This message responds with the information of an asset by its denomination.

```proto
message QueryGetAssetInfoResponse {
  AssetInfo asset_info = 1 [
    (gogoproto.nullable) = false
  ];
}
```

#### QueryAllAssetInfoRequest

This message requests the list of all asset information with pagination support.

```proto
message QueryAllAssetInfoRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllAssetInfoResponse

This message responds with the list of all asset information and pagination details.

```proto
message QueryAllAssetInfoResponse {
  repeated AssetInfo asset_info = 1 [
    (gogoproto.nullable) = false
  ];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### QueryGetPriceRequest

This message requests the price of an asset, optionally specifying the source and timestamp.

```proto
message QueryGetPriceRequest {
  string asset = 1;
  string source = 2;
  uint64 timestamp = 3;
}
```

#### QueryGetPriceResponse

This message responds with the price of an asset.

```proto
message QueryGetPriceResponse {
  Price price = 1 [
    (gogoproto.nullable) = false
  ];
}
```

#### QueryAllPriceRequest

This message requests the list of all prices with pagination support.

```proto
message QueryAllPriceRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllPriceResponse

This message responds with the list of all prices and pagination details.

```proto
message QueryAllPriceResponse {
  repeated Price price = 1 [
    (gogoproto.nullable) = false
  ];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

#### QueryGetPriceFeederRequest

This message requests the information of a price feeder by its identifier.

```proto
message QueryGetPriceFeederRequest {
  string feeder = 1;
}
```

#### QueryGetPriceFeederResponse

This message responds with the information of a price feeder by its identifier.

```proto
message QueryGetPriceFeederResponse {
  PriceFeeder price_feeder = 1 [
    (gogoproto.nullable) = false
  ];
}
```

#### QueryAllPriceFeederRequest

This message requests the list of all price feeders with pagination support.

```proto
message QueryAllPriceFeederRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllPriceFeederResponse

This message responds with the list of all price feeders and pagination details.

```proto
message QueryAllPriceFeederResponse {
  repeated PriceFeeder price_feeder = 1 [
    (gogoproto.nullable) = false
  ];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```
