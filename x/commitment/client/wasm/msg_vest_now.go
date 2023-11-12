package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgVestNow(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgVestNow) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "VestNow null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgVestNow := types.NewMsgVestNow(
		msg.Creator,
		msg.Amount,
		msg.Denom,
	)

	if err := msgVestNow.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgVestNow")
	}

	res, err := msgServer.VestNow(
		sdk.WrapSDKContext(ctx),
		msgVestNow,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "VestNow msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize VestNow response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
