package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/tokenomics/keeper"
)

// Messenger handles messages for the Tokenomics module.
type Messenger struct {
	keeper *keeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper) *Messenger {
	return &Messenger{keeper: keeper}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.TokenomicsCreateAirdrop != nil:
		return m.msgCreateAirdrop(ctx, contractAddr, msg.TokenomicsCreateAirdrop)
	case msg.TokenomicsUpdateAirdrop != nil:
		return m.msgUpdateAirdrop(ctx, contractAddr, msg.TokenomicsUpdateAirdrop)
	case msg.TokenomicsDeleteAirdrop != nil:
		return m.msgDeleteAirdrop(ctx, contractAddr, msg.TokenomicsDeleteAirdrop)
	case msg.TokenomicsCreateTimeBasedInflation != nil:
		return m.msgCreateTimeBasedInflation(ctx, contractAddr, msg.TokenomicsCreateTimeBasedInflation)
	case msg.TokenomicsUpdateTimeBasedInflation != nil:
		return m.msgUpdateTimeBasedInflation(ctx, contractAddr, msg.TokenomicsUpdateTimeBasedInflation)
	case msg.TokenomicsDeleteTimeBasedInflation != nil:
		return m.msgDeleteTimeBasedInflation(ctx, contractAddr, msg.TokenomicsDeleteTimeBasedInflation)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
