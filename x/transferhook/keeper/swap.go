package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/transferhook/types"
)

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

	amount, ok := sdk.NewIntFromString(data.Amount)
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

func (k Keeper) SwapExactAmountIn(ctx sdk.Context, addr sdk.AccAddress, tokenIn sdk.Coin, routes []ammtypes.SwapAmountInRoute) error {
	msg := &ammtypes.MsgSwapExactAmountIn{
		Sender:            addr.String(),
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: sdk.OneInt(),
	}
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	msgServer := ammkeeper.NewMsgServerImpl(k.ammKeeper)
	_, err := msgServer.SwapExactAmountIn(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}
	return nil
}
