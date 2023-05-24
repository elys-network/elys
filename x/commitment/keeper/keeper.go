package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/elys-network/elys/x/commitment/types"
)

// Interface declearation
type CommitmentKeeperI interface {
	// Initiate commitment according to standard staking
	StandardStakingToken(sdk.Context, string, string, string) error

	// Iterate all commitments
	IterateCommitments(sdk.Context, func(types.Commitments) (stop bool))

	// Update commitment
	SetCommitments(ctx sdk.Context, commitments types.Commitments)

	// Get commitment
	GetCommitments(sdk.Context, string) (types.Commitments, bool)
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

// Process standard staking elys token
// Create a commitment entity
func (k Keeper) StandardStakingToken(ctx sdk.Context, delegator string, validator string, denom string) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	_, err = sdk.ValAddressFromBech32(validator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert validator address from bech32")
	}

	/***********************************************************/
	////////////////// Delegator entity //////////////////////////
	/***********************************************************/
	// Get the Commitments for the delegator
	commitments, found := k.GetCommitments(ctx, delegator)
	if !found {
		commitments = types.Commitments{
			Creator:           delegator,
			CommittedTokens:   []*types.CommittedTokens{},
			UncommittedTokens: []*types.UncommittedTokens{},
		}
	}
	// Get the uncommitted tokens for the delegator
	uncommittedToken, _ := commitments.GetUncommittedTokensForDenom(denom)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedToken = &types.UncommittedTokens{
			Denom:  denom,
			Amount: sdk.ZeroInt(),
		}
		uncommittedTokens = append(uncommittedTokens, uncommittedToken)
		commitments.UncommittedTokens = uncommittedTokens
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, delegator),
			sdk.NewAttribute(types.AttributeAmount, sdk.ZeroInt().String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	/***************************************************************/
	////////////////////// Validator entity /////////////////////////
	// Get the Commitments for the validator
	commitments, found = k.GetCommitments(ctx, validator)
	if !found {
		commitments = types.Commitments{
			Creator:           validator,
			CommittedTokens:   []*types.CommittedTokens{},
			UncommittedTokens: []*types.UncommittedTokens{},
		}
	}
	// Get the uncommitted tokens for the validator
	uncommittedToken, _ = commitments.GetUncommittedTokensForDenom(denom)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedToken = &types.UncommittedTokens{
			Denom:  denom,
			Amount: sdk.ZeroInt(),
		}
		uncommittedTokens = append(uncommittedTokens, uncommittedToken)
		commitments.UncommittedTokens = uncommittedTokens
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, delegator, sdk.NewCoin(denom, sdk.ZeroInt()))

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, validator),
			sdk.NewAttribute(types.AttributeAmount, sdk.ZeroInt().String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}
