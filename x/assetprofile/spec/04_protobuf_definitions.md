<!--
order: 4
-->

# Types, Messages, Queries, and States

## Types

### Entry

The `Entry` message defines the properties of an asset, including its base denomination, IBC details, display attributes, and permissions.

```proto
message Entry {
  string base_denom = 1;
  uint64 decimals = 2;
  string denom = 3;
  string path = 4;
  string ibc_channel_id = 5;
  string ibc_counterparty_channel_id = 6;
  string display_name = 7;
  string display_symbol = 8;
  string network = 9;
  string address = 10;
  string external_symbol = 11;
  string transfer_limit = 12;
  repeated string permissions = 13;
  string unit_denom = 14;
  string ibc_counterparty_denom = 15;
  string ibc_counterparty_chain_id = 16;
  string authority = 17;
  bool commit_enabled = 18;
  bool withdraw_enabled = 19;
}
```

### Params

The `Params` message defines the parameters for the `assetprofile` module.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `assetprofile` module at genesis.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated Entry entry_list = 2 [(gogoproto.nullable) = false];
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `assetprofile` module.

```proto
service Msg {
  rpc CreateEntry (MsgCreateEntry) returns (MsgCreateEntryResponse);
  rpc UpdateEntry (MsgUpdateEntry) returns (MsgUpdateEntryResponse);
  rpc DeleteEntry (MsgDeleteEntry) returns (MsgDeleteEntryResponse);
}
```

#### MsgCreateEntry

This message creates a new asset entry.

```proto
message MsgCreateEntry {
  string authority = 1;
  string base_denom = 2;
  uint64 decimals = 3;
  string denom = 4;
  string path = 5;
  string ibc_channel_id = 6;
  string ibc_counterparty_channel_id = 7;
  string display_name = 8;
  string display_symbol = 9;
  string network = 10;
  string address = 11;
  string external_symbol = 12;
  string transfer_limit = 13;
  repeated string permissions = 14;
  string unit_denom = 15;
  string ibc_counterparty_denom = 16;
  string ibc_counterparty_chain_id = 17;
  bool commit_enabled = 18;
  bool withdraw_enabled = 19;
}

message MsgCreateEntryResponse {}
```

#### MsgUpdateEntry

This message updates an existing asset entry.

```proto
message MsgUpdateEntry {
  string authority = 1;
  string base_denom = 2;
  uint64 decimals = 3;
  string denom = 4;
  string path = 5;
  string ibc_channel_id = 6;
  string ibc_counterparty_channel_id = 7;
  string display_name = 8;
  string display_symbol = 9;
  string network = 10;
  string address = 11;
  string external_symbol = 12;
  string transfer_limit = 13;
  repeated string permissions = 14;
  string unit_denom = 15;
  string ibc_counterparty_denom = 16;
  string ibc_counterparty_chain_id = 17;
  bool commit_enabled = 18;
  bool withdraw_enabled = 19;
}

message MsgUpdateEntryResponse {}
```

#### MsgDeleteEntry

This message deletes an existing asset entry.

```proto
message MsgDeleteEntry {
  string authority = 1;
  string base_denom = 2;
}

message

 MsgDeleteEntryResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `assetprofile` module.

```proto
service Query {
  rpc Params (QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/assetprofile/params";
  }
  rpc Entry (QueryGetEntryRequest) returns (QueryGetEntryResponse) {
    option (google.api.http).get = "/elys-network/elys/assetprofile/entry/{base_denom}";
  }
  rpc EntryByDenom (QueryGetEntryByDenomRequest) returns (QueryGetEntryByDenomResponse) {
    option (google.api.http).get = "/elys-network/elys/assetprofile/entry/{denom}";
  }
  rpc EntryAll (QueryAllEntryRequest) returns (QueryAllEntryResponse) {
    option (google.api.http).get = "/elys-network/elys/assetprofile/entry";
  }
}
```

#### QueryParamsRequest

This message requests the parameters of the `assetprofile` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `assetprofile` module.

```proto
message QueryParamsResponse {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

#### QueryGetEntryRequest

This message requests an asset entry by its base denomination.

```proto
message QueryGetEntryRequest {
  string base_denom = 1;
}
```

#### QueryGetEntryResponse

This message responds with the asset entry details.

```proto
message QueryGetEntryResponse {
  Entry entry = 1 [(gogoproto.nullable) = false];
}
```

#### QueryGetEntryByDenomRequest

This message requests an asset entry by its denomination.

```proto
message QueryGetEntryByDenomRequest {
  string denom = 1;
}
```

#### QueryGetEntryByDenomResponse

This message responds with the asset entry details by denomination.

```proto
message QueryGetEntryByDenomResponse {
  Entry entry = 1 [(gogoproto.nullable) = false];
}
```

#### QueryAllEntryRequest

This message requests a list of all asset entries.

```proto
message QueryAllEntryRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}
```

#### QueryAllEntryResponse

This message responds with a list of all asset entries.

```proto
message QueryAllEntryResponse {
  repeated Entry entry = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```
