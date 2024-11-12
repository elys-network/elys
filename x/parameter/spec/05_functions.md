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

## UpdateWasmConfig

The `UpdateWasmConfig` function updates the WASM configuration parameters.

```go
func (k msgServer) UpdateWasmConfig(goCtx context.Context, msg *types.MsgUpdateWasmConfig) (*types.MsgUpdateWasmConfigResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    if k.authority != msg.Creator {
        return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
    }

    wasmMaxLabelSize, ok := math.NewIntFromString(msg.WasmMaxLabelSize)
    if !ok {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wasm max label size must be a positive integer")
    }

    wasmMaxProposalWasmSize, ok := math.NewIntFromString(msg.WasmMaxProposalWasmSize)
    if !ok {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wasm max proposal wasm size must be a positive integer")
    }

    wasmMaxSize, ok := math.NewIntFromString(msg.WasmMaxSize)


    if !ok {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wasm max size must be a positive integer")
    }

    params := k.GetParams(ctx)
    params.WasmMaxLabelSize = wasmMaxLabelSize
    params.WasmMaxProposalWasmSize = wasmMaxProposalWasmSize
    params.WasmMaxSize = wasmMaxSize
    k.SetParams(ctx, params)
    return &types.MsgUpdateWasmConfigResponse{}, nil
}
```
