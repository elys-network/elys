<!--
order: 2
-->

# State

## Params

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;

  repeated IncentiveInfo lp_incentives = 1 [(gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"lp_incentives\""];
  repeated IncentiveInfo stake_incentives = 2 [(gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"stake_incentives\""];
  string community_tax = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  bool withdraw_addr_enabled = 4;
}
```

## IncentiveInfo

```proto
message IncentiveInfo {
    // reward amount
    string amount = 1
        [ (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"amount\""];
    // epoch identifier
    string epoch_identifier = 2 [(gogoproto.moretags) = "yaml:\"epoch_identifier\""];
    // start_time of the distribution
    google.protobuf.Timestamp start_time = 3
        [(gogoproto.stdtime) = true, (gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"start_time\""];
    // distribution duration
    int64 num_epochs = 4 [(gogoproto.moretags) = "yaml:\"num_epochs\""];
    int64 current_epoch = 5 [(gogoproto.moretags) = "yaml:\"current_epoch\""];
    int64 eden_boost_apr = 6 [(gogoproto.moretags) = "yaml:\"eden_boost_apr\""];
}
```

`IncentiveInfo` describes the incentive information for specific pool.

## FeePool

```proto
message FeePool {
  repeated cosmos.base.v1beta1.DecCoin community_pool = 1 [
    (gogoproto.nullable)     = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.DecCoins"
  ];
}
```

`FeePool` describes the amount of tokens on community pool.
