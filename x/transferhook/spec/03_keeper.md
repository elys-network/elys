<!--
order: 3
-->

# Keeper

## Parameter Management

The `transferhook` module's keeper handles parameter management, allowing the activation or deactivation of the AMM functionality. It ensures that parameters are properly set and retrieved, enabling flexible module configuration.

### GetParams

The `GetParams` function retrieves the current parameters of the `transferhook` module.

```go
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
    k.paramstore.GetParamSet(ctx, &params)
    return params
}
```

### SetParams

The `SetParams` function sets the parameters of the `transferhook` module.

```go
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
    k.paramstore.SetParamSet(ctx, &params)
}
```

## Swap Functionality

The `transferhook` module also includes functionality to perform swaps during token transfers. This is integrated with the AMM module to provide seamless token exchange.

### Swap

The `Swap` function executes a token swap during an IBC transfer.

```go
func (k Keeper) Swap(
    ctx sdk.Context,
    packet channeltypes.Packet,
    data transfertypes.FungibleTokenPacketData,
    packetMetadata types.AmmPacketMetadata,
) error {
    params := k.GetParams(ctx)
    if !params.AmmActive {
        return errorsmod.Wrapf(types.ErrPacketForwardingInactive, "transferhook amm routing is inactive")
    }

    amount, ok := math.NewIntFromString(data.Amount)
    if !ok {
        return errors.New("not a parsable amount field")
    }

    recvDenom := ""
    if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
        voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
        unprefixedDenom := data.Denom[len(voucherPrefix):]
        recvDenom = unprefixedDenom
    } else {
        sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
        prefixedDenom := sourcePrefix + data.Denom
        recvDenom = transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()
    }

    receiverAddress, err := sdk.AccAddressFromBech32(data.Receiver)
    if err != nil {
        return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid elys_address (%s) in transferhook memo", receiverAddress)
    }

    return k.SwapExactAmountIn(ctx, receiverAddress, sdk.NewCoin(recvDenom, amount), packetMetadata.Routes)
}
```

### SwapExactAmountIn

The `SwapExactAmountIn` function performs the exact amount swap using the AMM module.

```go
func (k Keeper) SwapExactAmountIn(ctx sdk.Context, addr sdk.AccAddress, tokenIn sdk.Coin, routes []ammtypes.SwapAmountInRoute) error {
    msg := &ammtypes.MsgSwapExactAmountIn{
        Sender:            addr.String(),
        Routes:            routes,
        TokenIn:           tokenIn,
        TokenOutMinAmount: sdkmath.OneInt(),
        Discount:          sdkmath.LegacyZeroDec(),
    }
    if err := msg.ValidateBasic(); err != nil {
        return err
    }

    msgServer := ammkeeper.NewMsgServerImpl(k.ammKeeper)
    _, err := msgServer.SwapExactAmountIn(ctx, msg)
    if err != nil {
        return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
    }
    return nil
}
```
