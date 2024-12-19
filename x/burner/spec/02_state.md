<!--
order: 2
-->

# State

## State Objects

The `x/burner` module keeps the following `objects in state`:

| State Object  | Description                            | Key                                   | Value                                       | Store |
| ------------- | -------------------------------------- | ------------------------------------- | ------------------------------------------- | ----- |
| `HistoryList` | History of all the burnt tokens events | `[]History{timestamp, denom, amount}` | `[]History{1629792000, "uelys", "1000000"}` | KV    |

### History

A `History` defines several variables:

1. `timestamp` keeps the time the event occurred
2. `denom` keeps the token denom used
3. `amount` keeps the amount of tokens burnt

```protobuf
message History {
  string timestamp = 1;
  string denom = 2;
  string amount = 3;
}
```

The `burner` module keeps these `History` objects in state, which are initialized at genesis and are modified on begin blockers or end blockers.

### Genesis State

The `x/burner` module's `GenesisState` defines the state necessary for initializing the chain from a previously exported height. It contains a slice containing all the `History` objects kept in state:

```go
// GenesisState defines the burner module's genesis state.
message GenesisState {
           Params  params      = 1 [(gogoproto.nullable) = false];
  repeated History historyList = 2 [(gogoproto.nullable) = false];
}
```
