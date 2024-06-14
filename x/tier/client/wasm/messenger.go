package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/tier/keeper"
)

// Messenger handles messages for the Tier LP module.
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
	case msg.TierSetPortfolio != nil:
		return m.msgSetPortfolio(ctx, contractAddr, msg.TierSetPortfolio)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
