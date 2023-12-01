package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/margin/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// Messenger handles messages for the Margin module.
type Messenger struct {
	keeper          *keeper.Keeper
	parameterKeeper *parameterkeeper.Keeper
}

func NewMessenger(
	keeper *keeper.Keeper,
	parameterKeeper *parameterkeeper.Keeper,
) *Messenger {
	return &Messenger{
		keeper:          keeper,
		parameterKeeper: parameterKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.MarginOpen != nil:
		return m.msgOpen(ctx, contractAddr, msg.MarginOpen)
	case msg.MarginClose != nil:
		return m.msgClose(ctx, contractAddr, msg.MarginClose)
	case msg.MarginBrokerOpen != nil:
		return m.msgBrokerOpen(ctx, contractAddr, msg.MarginBrokerOpen)
	case msg.MarginBrokerClose != nil:
		return m.msgBrokerClose(ctx, contractAddr, msg.MarginBrokerClose)
	case msg.MarginUpdateParams != nil:
		return m.msgUpdateParams(ctx, contractAddr, msg.MarginUpdateParams)
	case msg.MarginUpdatePools != nil:
		return m.msgUpdatePools(ctx, contractAddr, msg.MarginUpdatePools)
	case msg.MarginWhitelist != nil:
		return m.msgWhitelist(ctx, contractAddr, msg.MarginWhitelist)
	case msg.MarginDewhitelist != nil:
		return m.msgDewhitelist(ctx, contractAddr, msg.MarginDewhitelist)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
