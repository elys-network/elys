package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/assetprofile/keeper"
	types "github.com/elys-network/elys/x/assetprofile/types"
)

func (m *Messenger) msgDeleteEntry(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgDeleteEntry) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "DeleteEntry null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgDeleteEntry := types.NewMsgDeleteEntry(
		msg.Authority,
		msg.BaseDenom,
	)

	if err := msgDeleteEntry.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgDeleteEntry")
	}

	res, err := msgServer.DeleteEntry(
		sdk.WrapSDKContext(ctx),
		msgDeleteEntry,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "DeleteEntry msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize DeleteEntry response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
