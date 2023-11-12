package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgUpdateVestingInfo(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdateVestingInfo) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdateVestingInfo null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgUpdateVestingInfo := types.NewMsgUpdateVestingInfo(
		msg.Authority,
		msg.BaseDenom,
		msg.VestingDenom,
		msg.EpochIdentifier,
		msg.NumEpochs,
		msg.VestNowFactor,
		msg.NumMaxVestings,
	)

	if err := msgUpdateVestingInfo.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgUpdateVestingInfo")
	}

	res, err := msgServer.UpdateVestingInfo(
		sdk.WrapSDKContext(ctx),
		msgUpdateVestingInfo,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdateVestingInfo msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdateVestingInfo response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
