package keeper

import (
	"context"
	"time"

	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/elys-network/elys/x/oracle/types"
)

// RequestBandPrice creates the BandPrice packet
// data with obi encoded and send it to the channel
func (k msgServer) RequestBandPrice(goCtx context.Context, msg *types.MsgRequestBandPrice) (*types.MsgRequestBandPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sourcePort := types.PortID
	channelCap, ok := k.ScopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, msg.SourceChannel))
	if !ok {
		return nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound,
			"module does not own channel capability")
	}

	encodedCalldata := obi.MustEncode(*msg.Calldata)
	packetData := packet.NewOracleRequestPacketData(
		msg.ClientID,
		msg.OracleScriptID,
		encodedCalldata,
		msg.AskCount,
		msg.MinCount,
		msg.FeeLimit,
		msg.PrepareGas,
		msg.ExecuteGas,
	)

	_, err := k.ChannelKeeper.SendPacket(ctx, channelCap, sourcePort, msg.SourceChannel, clienttypes.NewHeight(0, 0), uint64(ctx.BlockTime().UnixNano()+int64(10*time.Minute)), packetData.GetBytes())
	if err != nil {
		return nil, err
	}

	return &types.MsgRequestBandPriceResponse{}, nil
}
