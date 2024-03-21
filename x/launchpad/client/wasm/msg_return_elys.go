package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	launchpadkeeper "github.com/elys-network/elys/x/launchpad/keeper"
	launchpadtypes "github.com/elys-network/elys/x/launchpad/types"
)

func (m *Messenger) msgReturnElys(ctx sdk.Context, contractAddr sdk.AccAddress, msgReturnElys *launchpadtypes.MsgReturnElys) ([]sdk.Event, [][]byte, error) {
	if msgReturnElys == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "launchpad ReturnElys null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgReturnElys.Sender != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "launchpad ReturnElys wrong sender"}
	}

	res, err := PerformMsgReturn(m.keeper, ctx, contractAddr, msgReturnElys)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform launchpad ReturnElys")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize launchpad ReturnElys response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgReturn(f *launchpadkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgReturnElys *launchpadtypes.MsgReturnElys) (*launchpadtypes.MsgReturnElysResponse, error) {
	if msgReturnElys == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "launchpad ReturnElys"}
	}
	msgServer := launchpadkeeper.NewMsgServerImpl(*f)

	msgMsgReturnElys := launchpadtypes.NewMsgReturnElys(msgReturnElys.Sender, msgReturnElys.OrderId, msgReturnElys.ReturnElysAmount)

	if err := msgMsgReturnElys.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	_, err := msgServer.ReturnElys(sdk.WrapSDKContext(ctx), msgMsgReturnElys) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "launchpad ReturnElys msg")
	}

	resp := &launchpadtypes.MsgReturnElysResponse{}
	return resp, nil
}
