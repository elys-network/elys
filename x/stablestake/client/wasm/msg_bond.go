package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/stablestake/keeper"
	types "github.com/elys-network/elys/x/stablestake/types"
)

func (m *Messenger) msgBond(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgBond) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Bond null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgBond := types.NewMsgBond(
		msg.Creator,
		msg.Amount,
	)

	if err := msgBond.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgBond")
	}

	res, err := msgServer.Bond(
		sdk.WrapSDKContext(ctx),
		msgBond,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "Bond msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize Bond response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
