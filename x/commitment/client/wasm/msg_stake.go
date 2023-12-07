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

func (m *Messenger) msgStake(ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *commitmenttypes.MsgStake) ([]sdk.Event, [][]byte, error) {
	if msgStake == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgStake.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "stake wrong sender"}
	}

	entry, found := m.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid usdc denom"}
	}
	baseCurrency := entry.Denom

	var res *commitmenttypes.MsgStakeResponse
	var err error
	// USDC
	if msgStake.Asset == baseCurrency {
		msgServer := stablekeeper.NewMsgServerImpl(*m.stableKeeper)
		msgMsgBond := stabletypes.NewMsgBond(msgStake.Creator, msgStake.Amount)

		if err = msgMsgBond.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgBond")
		}

		_, err = msgServer.Bond(sdk.WrapSDKContext(ctx), msgMsgBond)
		if err != nil { // Discard the response because it's empty
			return nil, nil, errorsmod.Wrap(err, "usdc staking msg")
		}
		res = &commitmenttypes.MsgStakeResponse{
			Code:   ptypes.RES_OK,
			Result: "usdc staking msg succeed",
		}
	} else {
		// Elys, Eden, Eden Boost
		msgServer := commitmentkeeper.NewMsgServerImpl(*m.keeper)
		msgMsgStake := commitmenttypes.NewMsgStake(msgStake.Creator, msgStake.Amount, msgStake.Asset, msgStake.ValidatorAddress)

		if err = msgMsgStake.ValidateBasic(); err != nil {
			return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgStake")
		}

		res, err = msgServer.Stake(sdk.WrapSDKContext(ctx), msgMsgStake)
		if err != nil { // Discard the response because it's empty
			return nil, nil, errorsmod.Wrap(err, "elys staking msg")
		}
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
