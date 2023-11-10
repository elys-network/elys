package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/assetprofile/keeper"
	types "github.com/elys-network/elys/x/assetprofile/types"
)

func (m *Messenger) msgUpdateEntry(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdateEntry) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdateEntry null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgUpdateEntry := types.NewMsgUpdateEntry(
		msg.Authority,
		msg.BaseDenom,
		msg.Decimals,
		msg.Denom,
		msg.Path,
		msg.IbcChannelId,
		msg.IbcCounterpartyChannelId,
		msg.DisplayName,
		msg.DisplaySymbol,
		msg.Network,
		msg.Address,
		msg.ExternalSymbol,
		msg.TransferLimit,
		msg.Permissions,
		msg.UnitDenom,
		msg.IbcCounterpartyDenom,
		msg.IbcCounterpartyChainId,
	)

	if err := msgUpdateEntry.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgUpdateEntry")
	}

	res, err := msgServer.UpdateEntry(
		sdk.WrapSDKContext(ctx),
		msgUpdateEntry,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdateEntry msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdateEntry response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
