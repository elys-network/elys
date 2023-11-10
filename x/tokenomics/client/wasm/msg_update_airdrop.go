package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgUpdateAirdrop(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdateAirdrop) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdateAirdrop null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgUpdateAirdrop := types.NewMsgUpdateAirdrop(
		msg.Authority,
		msg.Intent,
		msg.Amount,
	)

	if err := msgUpdateAirdrop.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgUpdateAirdrop")
	}

	res, err := msgServer.UpdateAirdrop(
		sdk.WrapSDKContext(ctx),
		msgUpdateAirdrop,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdateAirdrop msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdateAirdrop response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
