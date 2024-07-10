<!--
order: 1
-->

# Concepts

The `clock` module in the Elys Network is designed to manage contract execution at the end of each block. It ensures that specific smart contracts are triggered automatically, allowing for scheduled tasks and automated processes within the blockchain.

## Module Components

1. **Keeper**: The core component that handles the business logic of the `Clock` module.
2. **End Blocker**: Executes specified contracts at the end of each block.
3. **WASM Integration**: Facilitates interaction with smart contracts through WASM.

## Key Functions

### 1. `EndBlocker`

The `EndBlocker` function is responsible for executing contracts at the end of each block. It performs the following steps:

- **Retrieve Parameters**: Fetches the current parameters, including contract addresses and gas limits.
- **Execute Contracts**: Iterates over the contract addresses and executes each one with a predefined message.
- **Log Errors**: Logs any errors that occur during contract execution.

**Function Signature**:

```go
func EndBlocker(ctx sdk.Context, k keeper.Keeper)
```

**Parameters**:

- `ctx`: The context of the current transaction.
- `k`: The Keeper instance.

**Return**:

- This function does not return a value but logs execution errors if any occur.

### 2. `SetParams`

The `SetParams` function sets the parameters for the `Clock` module.

**Function Signature**:

```go
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error
```

**Parameters**:

- `ctx`: The context of the current transaction.
- `p`: The parameters to set.

**Return**:

- Returns an error if parameter validation fails.

### 3. `GetParams`

The `GetParams` function retrieves the current parameters for the `Clock` module.

**Function Signature**:

```go
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params)
```

**Parameters**:

- `ctx`: The context of the current transaction.

**Return**:

- Returns the current parameters.

## Error Handling

The `Clock` module includes error handling to manage cases where contract execution fails. Errors are logged with details about the failed executions.

## Integration

The `Clock` module integrates with the WASM module to execute smart contracts. It uses the `Keeper` component to manage contract addresses and execution parameters.
