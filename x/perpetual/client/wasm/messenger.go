package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/perpetual/keeper"
)

// Messenger handles messages for the Perpetual module.
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
	case msg.PerpetualOpen != nil:
		return m.msgOpen(ctx, contractAddr, msg.PerpetualOpen)
	case msg.PerpetualClose != nil:
		return m.msgClose(ctx, contractAddr, msg.PerpetualClose)
	case msg.PerpetualAddCollateral != nil:
		return m.msgAddCollateral(ctx, contractAddr, msg.PerpetualAddCollateral)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
