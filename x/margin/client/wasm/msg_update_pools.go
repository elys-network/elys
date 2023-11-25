package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/margin/keeper"
	types "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgUpdatePools(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdatePools) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdatePools null msg"}
	}

	if msg.Authority != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "update pools wrong sender"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.UpdatePools(
		sdk.WrapSDKContext(ctx),
		msg,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdatePools msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdatePools response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
