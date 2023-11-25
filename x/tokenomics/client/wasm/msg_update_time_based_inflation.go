package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgUpdateTimeBasedInflation(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdateTimeBasedInflation) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdateTimeBasedInflation null msg"}
	}

	if msg.Authority != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "update time based inflation wrong sender"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.UpdateTimeBasedInflation(
		sdk.WrapSDKContext(ctx),
		msg,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdateTimeBasedInflation msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdateTimeBasedInflation response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
