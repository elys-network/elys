<!--
order: 3
-->

# Keeper

## Managing Commitments

The `commitment` module's keeper handles the logic for committing, vesting, and uncommitting tokens. It ensures that all token movements and state changes are properly recorded and managed.

### Setting Commitments

The `SetCommitments` function updates or creates a new commitment entry in the store.

```go
func (k Keeper) SetCommitments(ctx sdk.Context, commitments types.Commitments) {
    // Implementation here
}
```

### Getting Commitments

The `GetCommitments` function retrieves a commitment entry from the store based on the creator's address.

```go
func (k Keeper) GetCommitments(ctx sdk.Context, creator string) types.Commitments {
    // Implementation here
}
```

### Removing Commitments

The `RemoveCommitments` function deletes a commitment entry from the store.

```go
func (k Keeper) RemoveCommitments(ctx sdk.Context, creator string) {
    // Implementation here
}
```

### Updating Parameters

The `UpdateParams` function modifies the parameters of the `commitment` module.

```go
func (k Keeper) UpdateParams(ctx sdk.Context, params types.Params) {
    // Implementation here
}
```
