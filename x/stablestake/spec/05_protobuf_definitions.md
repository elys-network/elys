<!--
order: 5
-->

# Protobuf Definitions

## Types

### Debt

The `Debt` message tracks the borrowed amount, interest paid, and interest stacked for a given address, along with timestamps for borrowing and last interest calculation.

```proto
message Debt {
  string address = 1;
  string borrowed = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string interest_paid = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string interest_stacked = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  uint64 borrow_time = 5;
  uint64 last_interest_calc_time = 6;
}
```

### Params

The `Params` message defines the parameters for the `stablestake` module, including deposit denomination, redemption rate, epoch length, interest rates, and total value.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;

  string deposit_denom = 1;
  string redemption_rate = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  int64 epoch_length = 3;
  string interest_rate = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string interest_rate_max = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string interest_rate_min = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string interest_rate_increase = 7 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string interest_rate_decrease = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string health_gain_factor = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string total_value = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

### GenesisState

The `GenesisState` message defines the initial state of the `stablestake` module at genesis.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

## Messages

### Msg Service

The `Msg` service defines the transactions available in the `stablestake` module.

```proto
service Msg {
  rpc Bond(MsgBond) returns (MsgBondResponse);
  rpc Unbond(MsgUnbond) returns (MsgUnbondResponse);
}
```

#### MsgBond

This message allows a user to bond a specified amount of tokens.

```proto
message MsgBond {
  string creator = 1;
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}

message MsgBondResponse {}
```

#### MsgUnbond

This message allows a user to unbond a specified amount of tokens.

```proto
message MsgUnbond {
  string creator = 1;
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos

/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}

message MsgUnbondResponse {}
```

## Queries

### Query Service

The `Query` service defines the gRPC querier service for the `stablestake` module.

```proto
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/elys-network/elys/stablestake/params";
  }
  rpc BorrowRatio(QueryBorrowRatioRequest) returns (QueryBorrowRatioResponse) {
    option (google.api.http).get = "/elys-network/elys/stablestake/borrow-ratio";
  }
}
```

#### QueryParamsRequest

This message requests the parameters of the `stablestake` module.

```proto
message QueryParamsRequest {}
```

#### QueryParamsResponse

This message responds with the parameters of the `stablestake` module.

```proto
message QueryParamsResponse {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

#### QueryBorrowRatioRequest

This message requests the borrow ratio in the `stablestake` module.

```proto
message QueryBorrowRatioRequest {}
```

#### QueryBorrowRatioResponse

This message responds with the total deposits, total borrowings, and the borrow ratio in the `stablestake` module.

```proto
message QueryBorrowRatioResponse {
  string total_deposit = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string total_borrow = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string borrow_ratio = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```
