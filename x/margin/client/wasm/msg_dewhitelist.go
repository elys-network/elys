package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/margin/keeper"
	types "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgDewhitelist(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgDewhitelist) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Dewhitelist null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgDewhitelist := types.NewMsgDewhitelist(
		msg.Authority,
		msg.WhitelistedAddress,
	)

	if err := msgDewhitelist.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgDewhitelist")
	}

	res, err := msgServer.Dewhitelist(
		sdk.WrapSDKContext(ctx),
		msgDewhitelist,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "Dewhitelist msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize Dewhitelist response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
