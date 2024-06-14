package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tierkeeper "github.com/elys-network/elys/x/tier/keeper"
	tiertypes "github.com/elys-network/elys/x/tier/types"
)

func (m *Messenger) msgSetPortfolio(ctx sdk.Context, contractAddr sdk.AccAddress, msgSetPortfolio *tiertypes.MsgSetPortfolio) ([]sdk.Event, [][]byte, error) {
	if msgSetPortfolio == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Set portfolio null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgSetPortfolio.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "open wrong sender"}
	}

	res, err := PerformMsgSetPortfolio(m.keeper, ctx, contractAddr, msgSetPortfolio)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform set portfolio")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgSetPortfolio(f *tierkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgSetPortfolio *tiertypes.MsgSetPortfolio) (*tiertypes.MsgSetPortfolioResponse, error) {
	if msgSetPortfolio == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "tier set portfolio null"}
	}
	msgServer := tierkeeper.NewMsgServerImpl(*f)

	msgMsgSetPortfolio := tiertypes.NewMsgSetPortfolio(msgSetPortfolio.Creator, msgSetPortfolio.User)

	if err := msgMsgSetPortfolio.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgSetPortfolio")
	}

	_, err := msgServer.SetPortfolio(sdk.WrapSDKContext(ctx), msgMsgSetPortfolio) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "tier set portfolio msg")
	}

	resp := &tiertypes.MsgSetPortfolioResponse{}
	return resp, nil
}
