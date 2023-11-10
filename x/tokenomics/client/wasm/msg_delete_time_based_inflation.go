package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgDeleteTimeBasedInflation(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgDeleteTimeBasedInflation) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "DeleteTimeBasedInflation null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgDeleteTimeBasedInflation := types.NewMsgDeleteTimeBasedInflation(
		msg.Authority,
		msg.StartBlockHeight,
		msg.EndBlockHeight,
	)

	if err := msgDeleteTimeBasedInflation.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgDeleteTimeBasedInflation")
	}

	res, err := msgServer.DeleteTimeBasedInflation(
		sdk.WrapSDKContext(ctx),
		msgDeleteTimeBasedInflation,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "DeleteTimeBasedInflation msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize DeleteTimeBasedInflation response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
