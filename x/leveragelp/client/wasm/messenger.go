package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// Messenger handles messages for the Leverage LP module.
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
	case msg.LeveragelpOpen != nil:
		return m.msgOpen(ctx, contractAddr, msg.LeveragelpOpen)
	case msg.LeveragelpClose != nil:
		return m.msgClose(ctx, contractAddr, msg.LeveragelpClose)
	case msg.LeveragelpUpdateParams != nil:
		return m.msgUpdateParams(ctx, contractAddr, msg.LeveragelpUpdateParams)
	case msg.LeveragelpUpdatePools != nil:
		return m.msgUpdatePools(ctx, contractAddr, msg.LeveragelpUpdatePools)
	case msg.LeveragelpWhitelist != nil:
		return m.msgWhitelist(ctx, contractAddr, msg.LeveragelpWhitelist)
	case msg.LeveragelpDewhitelist != nil:
		return m.msgDewhitelist(ctx, contractAddr, msg.LeveragelpDewhitelist)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
