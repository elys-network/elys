package types

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmvm/v2/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg types.CosmosMsg) (events []sdk.Event, data [][]byte, msgResponses [][]*cdctypes.Any, err error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg ElysMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, nil, errorsmod.Wrap(err, "elys msg")
		}

		// Iterate over the module message handlers and dispatch to the appropriate one
		for _, handler := range m.moduleMessengers {
			event, resp, err := handler.HandleMsg(ctx, contractAddr, contractIBCPortID, contractMsg)
			if err != nil {
				if err == ErrCannotHandleMsg {
					// This handler cannot handle the message, try the next one
					continue
				}
				// Some other error occurred, return it
				return nil, nil, nil, err
			}
			// Message was handled successfully, return the result
			// TODO check if sending cdctypes.Any nil is correct
			return event, resp, nil, nil
		}

		// If no handler could handle the message, return an error
		return nil, nil, nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "unknown message type")
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}
