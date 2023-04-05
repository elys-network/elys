package oracle

import (
	"fmt"

	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// handleOraclePacket handles the result of the received BandChain oracles
// packet and saves the data into the KV database
func (im IBCModule) handleOraclePacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
) (channeltypes.Acknowledgement, error) {
	fmt.Println("handleOraclePacket1")
	var ack channeltypes.Acknowledgement
	var modulePacketData packet.OracleResponsePacketData
	if err := types.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &modulePacketData); err != nil {
		return ack, nil
	}
	fmt.Println("handleOraclePacket2", modulePacketData)

	switch modulePacketData.GetClientID() {

	case types.BandPriceClientIDKey:
		var BandPriceResult types.BandPriceResult
		if err := obi.Decode(modulePacketData.Result, &BandPriceResult); err != nil {
			ack = channeltypes.NewErrorAcknowledgement(err)
			return ack, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "cannot decode the BandPrice received packet")
		}
		reqID := types.OracleRequestID(modulePacketData.RequestID)
		im.keeper.SetBandPriceResult(ctx, reqID, BandPriceResult)
		fmt.Println("handleOraclePacket3", BandPriceResult)

		params := im.keeper.GetParams(ctx)
		request, err := im.keeper.GetBandRequest(ctx, reqID)
		if err != nil {
			return ack, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "historical request does not exist")
		}

		if len(request.Symbols) != len(BandPriceResult.Rates) {
			return ack, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "request and result count does not match")
		}

		for index, symbol := range request.Symbols {
			im.keeper.SetPrice(ctx, types.Price{
				Asset:    symbol,
				Price:    sdk.NewDecWithPrec(int64(BandPriceResult.Rates[index]), int64(params.Multiplier)),
				Source:   "bandchain",
				Provider: "automation",
			})
		}
		// this line is used by starport scaffolding # oracle/module/recv

	default:
		err := sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal,
			"oracle received packet not found: %s", modulePacketData.GetClientID())
		ack = channeltypes.NewErrorAcknowledgement(err)
		fmt.Println("handleOraclePacket4", err)
		return ack, err

	}
	ack = channeltypes.NewResultAcknowledgement(
		types.ModuleCdc.MustMarshalJSON(
			packet.NewOracleRequestPacketAcknowledgement(modulePacketData.RequestID),
		),
	)
	fmt.Println("handleOraclePacket5", modulePacketData.RequestID)
	return ack, nil
}

// handleOracleAcknowledgment handles the acknowledgment result from the BandChain
// request and saves the request-id into the KV database
func (im IBCModule) handleOracleAcknowledgment(
	ctx sdk.Context,
	ack channeltypes.Acknowledgement,
	modulePacket channeltypes.Packet,
) (*sdk.Result, error) {
	fmt.Println("handleOracleAcknowledgment-1", ack.Response)
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		fmt.Println("handleOracleAcknowledgment1", resp)
		var oracleAck packet.OracleRequestPacketAcknowledgement
		err := types.ModuleCdc.UnmarshalJSON(resp.Result, &oracleAck)
		if err != nil {
			return nil, nil
		}
		fmt.Println("handleOracleAcknowledgment2", oracleAck)

		var data packet.OracleRequestPacketData
		if err = types.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &data); err != nil {
			return nil, nil
		}
		requestID := types.OracleRequestID(oracleAck.RequestID)
		fmt.Println("handleOracleAcknowledgment3", requestID)

		switch data.GetClientID() {

		case types.BandPriceClientIDKey:
			fmt.Println("handleOracleAcknowledgment4", requestID)
			var RequestBandPrice types.BandPriceCallData
			if err = obi.Decode(data.GetCalldata(), &RequestBandPrice); err != nil {
				return nil, sdkerrors.Wrap(err,
					"cannot decode the BandPrice oracle acknowledgment packet")
			}
			fmt.Println("handleOracleAcknowledgment5", requestID)
			im.keeper.SetLastBandRequestId(ctx, requestID)
			im.keeper.SetBandRequest(ctx, requestID, RequestBandPrice)
			return &sdk.Result{}, nil
			// this line is used by starport scaffolding # oracle/module/ack

		default:
			fmt.Println("handleOracleAcknowledgment6", data.GetClientID())
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal,
				"oracle acknowledgment packet not found: %s", data.GetClientID())
		}
	}
	return nil, nil
}
