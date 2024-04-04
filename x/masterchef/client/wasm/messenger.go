package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/masterchef/keeper"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// Messenger handles messages for the Masterchef module.
type Messenger struct {
	keeper           *keeper.Keeper
	stakingKeeper    *stakingkeeper.Keeper
	commitmentKeeper *commitmentkeeper.Keeper
	parameterKeeper  *parameterkeeper.Keeper
	masterchefKeeper *masterchefkeeper.Keeper
}

func NewMessenger(
	keeper *keeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	parameterKeeper *parameterkeeper.Keeper,
) *Messenger {
	return &Messenger{
		keeper:           keeper,
		stakingKeeper:    stakingKeeper,
		commitmentKeeper: commitmentKeeper,
		parameterKeeper:  parameterKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.MasterchefClaimRewards != nil:
		return m.msgClaimRewards(ctx, contractAddr, msg.MasterchefClaimRewards)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
