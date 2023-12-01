package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (m *Messenger) msgUnstake(ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *commitmenttypes.MsgUnstake) ([]sdk.Event, [][]byte, error) {
	if msgUnstake == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid unstaking parameter"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgUnstake.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "unstake wrong sender"}
	}

	entry, found := m.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid usdc denom"}
	}
	baseCurrency := entry.Denom

	var res *commitmenttypes.MsgUnstakeResponse
	var err error
	// USDC
	if msgUnstake.Asset == baseCurrency {
		msgServer := stablekeeper.NewMsgServerImpl(*m.stableKeeper)
		msgMsgUnBond := stabletypes.NewMsgUnbond(msgUnstake.Address, msgUnstake.Amount)

		if err = msgMsgUnBond.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgBond")
		}

		_, err = msgServer.Unbond(sdk.WrapSDKContext(ctx), msgMsgUnBond)
		if err != nil { // Discard the response because it's empty
			return nil, nil, errorsmod.Wrap(err, "usdc unstaking msg")
		}
		res = &commitmenttypes.MsgUnstakeResponse{
			Code:   ptypes.RES_OK,
			Result: "usdc unstaking msg succeed",
		}
	} else {
		msgServer := commitmentkeeper.NewMsgServerImpl(*m.keeper)
		msgMsgUnstake := commitmenttypes.NewMsgUnstake(msgUnstake.Address, msgUnstake.Amount, msgUnstake.Asset, msgUnstake.ValidatorAddress)

		if err = msgMsgUnstake.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgUnstake")
		}

		res, err = msgServer.Unstake(sdk.WrapSDKContext(ctx), msgMsgUnstake)
		if err != nil { // Discard the response because it's empty
			return nil, nil, errorsmod.Wrap(err, "elys unstaking msg")
		}
	}
	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize unstake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
