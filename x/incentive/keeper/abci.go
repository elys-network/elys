package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/incentive/types"
)

// EndBlocker of incentive module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	// Burn EdenB tokens if staking changed
	k.BurnEdenBIfElysStakingReduced(ctx)
}

func (k Keeper) TakeDelegationSnapshot(ctx sdk.Context, addr string) {
	// Calculate delegated amount per delegator
	delAmount := k.CalcDelegationAmount(ctx, addr)

	elysStaked := types.ElysStaked{
		Address: addr,
		Amount:  delAmount,
	}

	// Set Elys staked amount
	k.SetElysStaked(ctx, elysStaked)
}

func (k Keeper) BurnEdenBIfElysStakingReduced(ctx sdk.Context) {
	addrs := k.GetAllElysStakeChange(ctx)

	// Handle addresses recorded on AfterDelegationModified
	// This hook is exposed for genesis delegations as well
	for _, delAddr := range addrs {
		k.BurnEdenBFromElysUnstaking(ctx, delAddr)
		k.TakeDelegationSnapshot(ctx, delAddr.String())
		k.RemoveElysStakeChange(ctx, delAddr)
	}
}
