package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	launchpadkeeper "github.com/elys-network/elys/x/launchpad/keeper"
	launchpadtypes "github.com/elys-network/elys/x/launchpad/types"
)

func (m *Messenger) msgBuyElys(ctx sdk.Context, contractAddr sdk.AccAddress, msgBuyElys *launchpadtypes.MsgBuyElys) ([]sdk.Event, [][]byte, error) {
	if msgBuyElys == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "BuyElys null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgBuyElys.Sender != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "BuyElys wrong sender"}
	}

	res, err := PerformMsgBuy(m.keeper, ctx, contractAddr, msgBuyElys)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform BuyElys")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize BuyElys response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgBuy(f *launchpadkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgBuyElys *launchpadtypes.MsgBuyElys) (*launchpadtypes.MsgBuyElysResponse, error) {
	if msgBuyElys == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "launchpad buyElys null"}
	}
	msgServer := launchpadkeeper.NewMsgServerImpl(*f)

	msgMsgBuyElys := launchpadtypes.NewMsgBuyElys(msgBuyElys.Sender, msgBuyElys.SpendingToken, msgBuyElys.TokenAmount)

	if err := msgMsgBuyElys.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating MsgBuyElys")
	}

	_, err := msgServer.BuyElys(sdk.WrapSDKContext(ctx), msgMsgBuyElys) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "launchpad BuyElys msg")
	}

	resp := &launchpadtypes.MsgBuyElysResponse{}
	return resp, nil
}
