package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/clock/keeper"
	types "github.com/elys-network/elys/x/clock/types"
)

func (m *Messenger) msgUpdateParams(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdateParams) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdateParams null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	acc, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to get acc")
	}

	msgUpdateParams := types.NewMsgUpdateParams(
		acc,
		msg.Params.ContractAddresses,
	)

	if err := msgUpdateParams.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgUpdateParams")
	}

	res, err := msgServer.UpdateParams(
		sdk.WrapSDKContext(ctx),
		msgUpdateParams,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdateParams msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdateParams response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
