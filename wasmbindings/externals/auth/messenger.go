package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	wasmbindingstypes "github.com/elys-network/elys/v6/wasmbindings/types"
)

// Messenger handles messages for the Auth module.
type Messenger struct {
	keeper *keeper.AccountKeeper
}

func NewMessenger(keeper *keeper.AccountKeeper) *Messenger {
	return &Messenger{keeper: keeper}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
