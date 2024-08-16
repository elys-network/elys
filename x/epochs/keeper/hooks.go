package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/epochs/types"
)

var _ types.EpochHooks = MultiEpochHooks{}

// combine multiple epoch hooks, all hook functions are run in array sequence
type MultiEpochHooks []types.EpochHooks

func NewMultiEpochHooks(hooks ...types.EpochHooks) MultiEpochHooks {
	return hooks
}

// AfterEpochEnd is called when epoch is going to be ended, epochNumber is the
// number of epoch that is ending
func (mh MultiEpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	for i := range mh {
		err := mh[i].AfterEpochEnd(ctx, epochIdentifier, epochNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

// BeforeEpochStart is called when epoch is going to be started, epochNumber is
// the number of epoch that is starting
func (mh MultiEpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error{
	for i := range mh {
		err := mh[i].BeforeEpochStart(ctx, epochIdentifier, epochNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

// AfterEpochEnd executes the indicated hook after epochs ends
func (k Keeper) AfterEpochEnd(ctx sdk.Context, identifier string, epochNumber int64) error{
	if k.hooks == nil {
		return nil
	}
	return k.hooks.AfterEpochEnd(ctx, identifier, epochNumber)
}

// BeforeEpochStart executes the indicated hook before the epochs
func (k Keeper) BeforeEpochStart(ctx sdk.Context, identifier string, epochNumber int64) error{
	if k.hooks == nil {
		return nil 
	}

	return k.hooks.BeforeEpochStart(ctx, identifier, epochNumber)
}
