package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/elys-network/elys/x/commitment/types"
)

// Interface declearation
type CommitmentKeeperI interface {
	// Iterate all commitments
	IterateCommitments(sdk.Context, func(types.Commitments) (stop bool))

	// Update commitment
	SetCommitments(ctx sdk.Context, commitments types.Commitments)

	// Get commitment
	GetCommitments(sdk.Context, string) types.Commitments

	// Update commitments for claim reward operation
	RecordClaimReward(sdk.Context, string, string, sdk.Int, types.EarnType) error

	// Update commitments for validator's commission withdrawal to self delegator
	RecordWithdrawValidatorCommission(sdk.Context, string, string, string, sdk.Int) error

	// Update commitments for Withdraw tokens - only USDC
	RecordWithdrawUSDC(ctx sdk.Context, creator string, denom string, amount sdk.Int) error
}

var _ CommitmentKeeperI = Keeper{}

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		hooks      types.CommitmentHooks

		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		apKeeper      types.AssetProfileKeeper
		authority     string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	apKeeper types.AssetProfileKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		apKeeper:      apKeeper,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) BankKeeper() types.BankKeeper {
	return k.bankKeeper
}

// SetHooks set the epoch hooks
func (k *Keeper) SetHooks(eh types.CommitmentHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set epochs hooks twice")
	}

	k.hooks = eh

	return k
}
