<!--
order: 6
-->

### Pool Status

#### Overview

Managing pool status is crucial. Two key states are whether a pool is enabled or closed.

#### Pool Status Functions

```go
func (k Keeper) IsPoolEnabled(ctx sdk.Context, poolId uint64) bool {
    pool, found := k.GetPool(ctx, poolId)
    if (!found) {
        pool = types.NewPool(poolId)
        k.SetPool(ctx, pool)
    }
    return pool.Enabled
}

func (k Keeper) IsPoolClosed(ctx sdk.Context, poolId uint64) bool {
    pool, found := k.GetPool(ctx, poolId)
    if (!found) {
        pool = types.NewPool(poolId)
        k.SetPool(ctx, pool)
    }
    return pool.Closed
}
```

#### Functionality

1. **IsPoolEnabled**:

   - Checks if a pool is enabled.
   - If not found, initializes and sets a new pool.
   - An enabled pool is processed by the module.

2. **IsPoolClosed**:
   - Checks if a pool is closed.
   - If not found, initializes and sets a new pool.
   - A closed pool prevents new positions but allows existing ones to be processed.

#### Key Differences

- **Enabled**: The pool is either processed or excluded entirely.
- **Closed**: Only affects the opening of new positions; existing positions continue as usual.
