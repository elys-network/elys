package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgCreateAirdrop(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgCreateAirdrop) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "CreateAirdrop null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgCreateAirdrop := types.NewMsgCreateAirdrop(
		msg.Authority,
		msg.Intent,
		msg.Amount,
	)

	if err := msgCreateAirdrop.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgCreateAirdrop")
	}

	res, err := msgServer.CreateAirdrop(
		sdk.WrapSDKContext(ctx),
		msgCreateAirdrop,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "CreateAirdrop msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize CreateAirdrop response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
