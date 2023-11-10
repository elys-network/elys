package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tokenomics/keeper"
	types "github.com/elys-network/elys/x/tokenomics/types"
)

func (m *Messenger) msgUpdateGenesisInflation(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUpdateGenesisInflation) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "UpdateGenesisInflation null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgUpdateGenesisInflation := types.NewMsgUpdateGenesisInflation(
		msg.Authority,
		*msg.Inflation,
		msg.SeedVesting,
		msg.StrategicSalesVesting,
	)

	if err := msgUpdateGenesisInflation.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgUpdateGenesisInflation")
	}

	res, err := msgServer.UpdateGenesisInflation(
		sdk.WrapSDKContext(ctx),
		msgUpdateGenesisInflation,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "UpdateGenesisInflation msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize UpdateGenesisInflation response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
