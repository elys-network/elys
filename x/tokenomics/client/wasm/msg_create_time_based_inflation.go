package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgCreateTimeBasedInflation(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgCreateTimeBasedInflation) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "CreateTimeBasedInflation null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgCreateTimeBasedInflation := types.NewMsgCreateTimeBasedInflation(
		msg.Authority,
		msg.StartBlockHeight,
		msg.EndBlockHeight,
		msg.Description,
		msg.Inflation,
	)

	if err := msgCreateTimeBasedInflation.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgCreateTimeBasedInflation")
	}

	res, err := msgServer.CreateTimeBasedInflation(
		sdk.WrapSDKContext(ctx),
		msgCreateTimeBasedInflation,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "CreateTimeBasedInflation msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize CreateTimeBasedInflation response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
