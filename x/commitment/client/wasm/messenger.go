package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	assetProfileKeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	"github.com/elys-network/elys/x/commitment/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	stableKeeper "github.com/elys-network/elys/x/stablestake/keeper"
)

// Messenger handles messages for the Commitment module.
type Messenger struct {
	keeper             *keeper.Keeper
	stakingKeeper      *stakingkeeper.Keeper
	assetProfileKeeper *assetProfileKeeper.Keeper
	stableKeeper       *stableKeeper.Keeper
	parameterKeeper    *parameterkeeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper, stakingKeeper *stakingkeeper.Keeper, assetProfileKeeper *assetProfileKeeper.Keeper, stableKeeper *stableKeeper.Keeper, parameterKeeper *parameterkeeper.Keeper) *Messenger {
	return &Messenger{
		keeper:             keeper,
		stakingKeeper:      stakingKeeper,
		assetProfileKeeper: assetProfileKeeper,
		stableKeeper:       stableKeeper,
		parameterKeeper:    parameterKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.CommitmentCommitLiquidTokens != nil:
		return m.msgCommitLiquidTokens(ctx, contractAddr, msg.CommitmentCommitLiquidTokens)
	case msg.CommitmentCommitUnclaimedRewards != nil:
		return m.msgCommitClaimedRewards(ctx, contractAddr, msg.CommitmentCommitUnclaimedRewards)
	case msg.CommitmentUncommitTokens != nil:
		return m.msgUncommitTokens(ctx, contractAddr, msg.CommitmentUncommitTokens)
	case msg.CommitmentVest != nil:
		return m.msgVest(ctx, contractAddr, msg.CommitmentVest)
	case msg.CommitmentVestNow != nil:
		return m.msgVestNow(ctx, contractAddr, msg.CommitmentVestNow)
	case msg.CommitmentCancelVest != nil:
		return m.msgCancelVest(ctx, contractAddr, msg.CommitmentCancelVest)
	case msg.CommitmentUpdateVestingInfo != nil:
		return m.msgUpdateVestingInfo(ctx, contractAddr, msg.CommitmentUpdateVestingInfo)
	case msg.CommitmentStake != nil:
		return m.msgStake(ctx, contractAddr, msg.CommitmentStake)
	case msg.CommitmentUnstake != nil:
		return m.msgUnstake(ctx, contractAddr, msg.CommitmentUnstake)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
