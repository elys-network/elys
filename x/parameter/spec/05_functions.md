<!--
order: 5
-->

# Functions

## UpdateMinCommission

The `UpdateMinCommission` function updates the minimum commission rate.

```go
func (k msgServer) UpdateMinCommission(goCtx context.Context, msg *types.MsgUpdateMinCommission) (*types.MsgUpdateMinCommissionResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    if k.authority != msg.Creator {
        return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
    }

    minCommission, err := sdk.LegacyNewDecFromStr(msg.MinCommission)
    if err != nil {
        return nil, err
    }

    params := k.GetParams(ctx)
    params.MinCommissionRate = minCommission
    k.SetParams(ctx, params)
    return &types.MsgUpdateMinCommissionResponse{}, nil
}
```