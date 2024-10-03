package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

func (m *Messenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *leveragelptypes.MsgOpen) ([]sdk.Event, [][]byte, error) {
	if msgOpen == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Open null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgOpen.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "open wrong sender"}
	}

	res, err := PerformMsgOpen(m.keeper, ctx, contractAddr, msgOpen)
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

func PerformMsgOpen(f *leveragelpkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *leveragelptypes.MsgOpen) (*leveragelptypes.MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "leveragelp open null leveragelp open"}
	}
	msgServer := leveragelpkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := leveragelptypes.NewMsgOpen(msgOpen.Creator, msgOpen.CollateralAsset, msgOpen.CollateralAmount, msgOpen.AmmPoolId, msgOpen.Leverage, msgOpen.StopLossPrice)

	if err := msgMsgOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	_, err := msgServer.Open(sdk.WrapSDKContext(ctx), msgMsgOpen) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "leveragelp open msg")
	}

	resp := &leveragelptypes.MsgOpenResponse{}
	return resp, nil
}
