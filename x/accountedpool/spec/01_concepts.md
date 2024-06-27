<!--
order: 1
-->

# Concepts

The `accountedpool` module in the Elys Network manages accounted pools, ensuring accurate accounting of assets within pools. It integrates with Automated Market Maker (AMM) and Perpetual pools, updating pool states based on their activities.

## Module Components

1. **Keeper**: The core component that handles the business logic of the `AccountedPool` module.
2. **Hooks**: Functions that trigger updates to the accounted pool based on specific events in the AMM and Perpetual modules.

## Key Functions

### 1. `UpdateAccountedPool`

The `UpdateAccountedPool` function is responsible for recalculating and updating the balance of an accounted pool. It performs the following steps:

- **Check Pool Existence**: Ensures the pool exists before attempting to update it.
- **Retrieve Accounted Pool**: Fetches the current state of the accounted pool.
- **Calculate New Balance**: Computes the new balance of the accounted pool by summing the balances from the AMM pool and the Perpetual pool, including liabilities but excluding already deducted custody amounts.
- **Update State**: Saves the updated accounted pool state back to the store.

**Function Signature**:

```go
func (k Keeper) UpdateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) error
```

**Parameters**:

- `ctx`: The context of the current transaction.
- `ammPool`: The AMM pool instance.
- `perpetualPool`: The Perpetual pool instance.

**Return**:

- Returns an error if the pool does not exist or if there is an issue retrieving balances.

### 2. `InitiateAccountedPool`

The `InitiateAccountedPool` function is called after the creation of a new AMM pool to initialize its corresponding accounted pool.

**Function Signature**:

```go
func (k Keeper) InitiateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool)
```

**Parameters**:

- `ctx`: The context of the current transaction.
- `ammPool`: The AMM pool instance.

## Hooks

Hooks are specialized functions that automatically trigger updates to the accounted pool based on specific actions within the AMM and Perpetual modules. These actions include opening, modifying, or closing perpetual positions, as well as creating, joining, exiting, or swapping in AMM pools.

**Hook Functions**:

- **AfterPerpetualPositionOpen**
- **AfterPerpetualPositionModified**
- **AfterPerpetualPositionClosed**
- **AfterAmmPoolCreated**
- **AfterAmmJoinPool**
- **AfterAmmExitPool**
- **AfterAmmSwap**

Each hook function calls `UpdateAccountedPool` or `InitiateAccountedPool` as appropriate, ensuring the accounted pool's state remains consistent with the latest transactions.

**Example**:

```go
func (k Keeper) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

## Error Handling

The `AccountedPool` module includes error handling to manage cases where pools do not exist or balances cannot be retrieved. Notable errors include:

- `types.ErrPoolDoesNotExist`

These errors ensure that the module can gracefully handle unexpected states and provide meaningful feedback for debugging and resolution.

## Integration

The `AccountedPool` module integrates with the AMM and Perpetual modules, leveraging their pool data to maintain an accurate and up-to-date accounted pool state. By responding to various lifecycle events in these modules, the `AccountedPool` module ensures comprehensive and dynamic pool management.
