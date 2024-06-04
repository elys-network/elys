<!--
order: 2
-->

# State

```proto
// GenesisState defines the leveragelp module's genesis state.
message GenesisState {
  Params params = 1 [ (gogoproto.nullable) = false ];
  repeated Pool pool_list = 2 [ (gogoproto.nullable) = false ];
  repeated Position position_list = 3 [ (gogoproto.nullable) = false ];
  repeated string address_whitelist = 4;
}
```

## Params

`Params` holds module parameters like protocol revenue address, supported reward denoms, max eden apr, reward portion for lps and stakers and yearly incentive data.

```proto
// Params defines the parameters for the module.
message Params {
  option (gogoproto.goproto_stringer) = false;
  string leverage_max = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  int64 max_open_positions = 2;
  string pool_open_threshold = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string safety_factor = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  bool whitelisting_enabled = 5;
  int64 epoch_length = 6;
}
```

## Pool

`IncentiveInfo` holds eden reward distribution info for current year.

```proto
message Pool {
  uint64 amm_pool_id = 1;
  string health = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  bool enabled = 3;
  bool closed = 4;
  string leveraged_lp_amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string leverage_max = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

## Position

`Position` holds position information including health, `stopLossPrice`, `liabilities`, `leverage`, `leveragedLpAmount`.

```proto
message Position {
  string address = 1;
  cosmos.base.v1beta1.Coin collateral = 2 [(gogoproto.nullable) = false];
  string liabilities = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ]; // For recording
  string interest_paid = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ]; // For recording
  string leverage = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ]; // For recording
  string leveraged_lp_amount = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string position_health = 7 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  uint64 id = 8;
  uint64 amm_pool_id = 9;
  string stop_loss_price = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```
