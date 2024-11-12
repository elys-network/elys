package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

func (m *Messenger) msgClaimRewards(ctx sdk.Context, contractAddr sdk.AccAddress, msgClaimRewards *leveragelptypes.MsgClaimRewards) ([]sdk.Event, [][]byte, error) {
	if msgClaimRewards == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "ClaimRewards null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgClaimRewards.Sender != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "ClaimRewards wrong sender"}
	}

	if err := msgClaimRewards.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgClaimRewards")
	}

	msgServer := leveragelpkeeper.NewMsgServerImpl(*m.keeper)
	res, err := msgServer.ClaimRewards(sdk.WrapSDKContext(ctx), msgClaimRewards)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "leveragelp ClaimRewards msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize ClaimRewards response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
