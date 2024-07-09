<!--
order: 1
-->

# Concepts

The `parameter` module in the Elys Network manages key configuration parameters for the network, enabling dynamic updates and retrievals. It facilitates the modification of various parameters related to commissions, voting power, self-delegation, broker address, total blocks per year, and WASM configuration.

## Module Components

1. **Keeper**: The core component that handles the business logic of the `parameter` module.
2. **MsgServer**: Handles the gRPC messages to update parameters.
3. **Querier**: Handles the queries to retrieve parameter values.

## Key Functions

### 1. `UpdateMinCommission`

The `UpdateMinCommission` function updates the minimum commission rate for validators. It performs the following steps:

- **Check Authority**: Ensures the caller has the authority to update the parameter.
- **Validate Commission Rate**: Checks if the new commission rate is valid.
- **Update State**: Saves the updated commission rate to the store.

**Function Signature**:

```go
func (k msgServer) UpdateMinCommission(goCtx context.Context, msg *types.MsgUpdateMinCommission) (*types.MsgUpdateMinCommissionResponse, error)
```

**Parameters**:

- `goCtx`: The context of the current transaction.
- `msg`: The message containing the new commission rate.

**Return**:

- Returns an error if the caller lacks authority or if the commission rate is invalid.

### 2. `UpdateWasmConfig`

The `UpdateWasmConfig` function updates the configuration parameters for WASM. It includes the maximum label size, maximum size, and maximum proposal size.

**Function Signature**:

```go
func (k msgServer) UpdateWasmConfig(goCtx context.Context, msg *types.MsgUpdateWasmConfig) (*types.MsgUpdateWasmConfigResponse, error)
```

**Parameters**:

- `goCtx`: The context of the current transaction.
- `msg`: The message containing the new WASM configuration values.

**Return**:

- Returns an error if the caller lacks authority or if the parameters are invalid.

## Hooks

The `parameter` module does not include hooks as it primarily deals with static configuration parameters rather than dynamic state changes.

## Error Handling

The `parameter` module includes error handling to manage invalid parameters or unauthorized access. Notable errors include:

- `types.ErrInvalidMinCommissionRate`
- `types.ErrInvalidMaxVotingPower`
- `types.ErrInvalidMinSelfDelegation`
- `types.ErrInvalidBrokerAddress`

These errors ensure that the module can gracefully handle invalid updates and provide meaningful feedback for debugging and resolution.

## Integration

The `parameter` module integrates with other modules by providing a centralized mechanism for updating and retrieving configuration parameters. It ensures consistent and controlled parameter management across the network.
