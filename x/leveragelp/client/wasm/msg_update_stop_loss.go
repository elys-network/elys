package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

func (m *Messenger) msgUpdateStopLoss(ctx sdk.Context, contractAddr sdk.AccAddress, msgUpdateStopLoss *leveragelptypes.MsgUpdateStopLoss) ([]sdk.Event, [][]byte, error) {
	if msgUpdateStopLoss == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Open null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgUpdateStopLoss.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "open wrong sender"}
	}

	res, err := PerformMsgUpdateStopLoss(m.keeper, ctx, contractAddr, msgUpdateStopLoss)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform open")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgUpdateStopLoss(f *leveragelpkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgUpdateStopLoss *leveragelptypes.MsgUpdateStopLoss) (*leveragelptypes.MsgUpdateStopLossResponse, error) {
	if msgUpdateStopLoss == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "leveragelp update stop loss null leveragelp stop loss"}
	}
	msgServer := leveragelpkeeper.NewMsgServerImpl(*f)

	msgMsgUpdateStopLoss := leveragelptypes.NewMsgUpdateStopLoss(msgUpdateStopLoss.Creator, msgUpdateStopLoss.Position, msgUpdateStopLoss.Price)

	if err := msgMsgUpdateStopLoss.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgUpdateStopLoss")
	}

	_, err := msgServer.UpdateStopLoss(sdk.WrapSDKContext(ctx), msgMsgUpdateStopLoss) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "leveragelp stop loss msg")
	}

	resp := &leveragelptypes.MsgUpdateStopLossResponse{}
	return resp, nil
}
