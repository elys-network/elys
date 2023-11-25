package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgDeleteAirdrop(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgDeleteAirdrop) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "DeleteAirdrop null msg"}
	}

	if msg.Authority != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "create airdrop wrong sender"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.DeleteAirdrop(
		sdk.WrapSDKContext(ctx),
		msg,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "DeleteAirdrop msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize DeleteAirdrop response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
