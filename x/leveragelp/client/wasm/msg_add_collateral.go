package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

func (m *Messenger) msgAddCollateral(ctx sdk.Context, contractAddr sdk.AccAddress, msgAddCollateral *leveragelptypes.MsgAddCollateral) ([]sdk.Event, [][]byte, error) {
	if msgAddCollateral == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "add collateral null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgAddCollateral.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "add collateral wrong sender"}
	}

	res, err := PerformMsgAddCollateral(m.keeper, ctx, contractAddr, msgAddCollateral)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform add collateral loss")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize add collateral response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgAddCollateral(f *leveragelpkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgUpdateStopLoss *leveragelptypes.MsgAddCollateral) (*leveragelptypes.MsgAddCollateralResponse, error) {
	if msgUpdateStopLoss == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "leveragelp add collateral null leveragelp add collateral"}
	}
	msgServer := leveragelpkeeper.NewMsgServerImpl(*f)

	msgMsgAddCollateral := leveragelptypes.NewMsgAddCollateral(msgUpdateStopLoss.Creator, msgUpdateStopLoss.Id, msgUpdateStopLoss.Collateral)

	if err := msgMsgAddCollateral.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgAddCollateral")
	}

	_, err := msgServer.AddCollateral(sdk.WrapSDKContext(ctx), msgMsgAddCollateral) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "leveragelp add collateral msg")
	}

	resp := &leveragelptypes.MsgAddCollateralResponse{}
	return resp, nil
}
