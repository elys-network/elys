package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgUnstake(ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *commitmenttypes.MsgUnstake) ([]sdk.Event, [][]byte, error) {
	if msgUnstake == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid unstaking parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*m.keeper)
	msgMsgUnstake := commitmenttypes.NewMsgUnstake(msgUnstake.Creator, msgUnstake.Amount, msgUnstake.Asset, msgUnstake.ValidatorAddress)

	if err := msgMsgUnstake.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgUnstake")
	}

	res, err := msgServer.Unstake(ctx, msgMsgUnstake)
	if err != nil { // Discard the response because it's empty
		return nil, nil, errorsmod.Wrap(err, "elys unstaking msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize unstake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
