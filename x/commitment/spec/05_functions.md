<!--
order: 5
-->

# Functions

## EndBlocker

The `EndBlocker` function is called at the end of each block to perform necessary updates and maintenance for the `commitment` module.

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
    // Implementation here
}
```

### CommitClaimedRewards

The `CommitClaimedRewards` function commits the tokens from the claimed store to the committed store.

```go
func (k Keeper) CommitClaimedRewards(ctx sdk.Context, msg *types.MsgCommitClaimedRewards) (*types.MsgCommitClaimedRewardsResponse, error) {
    // Implementation here
}
```

### UncommitTokens

The `UncommitTokens` function uncommits

tokens from the committed store and makes them liquid.

```go
func (k Keeper) UncommitTokens(ctx sdk.Context, msg *types.MsgUncommitTokens) (*types.MsgUncommitTokensResponse, error) {
    // Implementation here
}
```

### Vest

The `Vest` function converts a user's committed tokens to vesting tokens.

```go
func (k Keeper) Vest(ctx sdk.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
    // Implementation here
}
```

### ClaimVesting

The `ClaimVesting` function claims the tokens that have already vested.

```go
func (k Keeper) ClaimVesting(ctx sdk.Context, msg *types.MsgClaimVesting) (*types.MsgClaimVestingResponse, error) {
    // Implementation here
}
```
