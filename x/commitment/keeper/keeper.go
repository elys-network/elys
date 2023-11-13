package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"

	ptypes "github.com/elys-network/elys/x/parameter/types"
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

	// Withdraw tokens
	// context, creator, denom, amount
	ProcessClaimReward(sdk.Context, string, string, sdk.Int) error

	// Withdraw validator commission
	// context, delegator, validator, denom, amount
	ProcessWithdrawValidatorCommission(sdk.Context, string, string, string, sdk.Int) error

	// Withdraw tokens - only USDC
	// context, creator, denom, amount
	ProcessWithdrawUSDC(ctx sdk.Context, creator string, denom string, amount sdk.Int) error
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
			Creator:          delegator,
			CommittedTokens:  []*types.CommittedTokens{},
			RewardsUnclaimed: sdk.Coins{},
		}
		k.SetCommitments(ctx, commitments)
	}

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
			Creator:          validator,
			CommittedTokens:  []*types.CommittedTokens{},
			RewardsUnclaimed: sdk.Coins{},
		}
		k.SetCommitments(ctx, commitments)
	}

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, delegator, sdk.Coins{sdk.NewCoin(denom, sdk.ZeroInt())})

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

// Withdraw Token
func (k Keeper) ProcessClaimReward(ctx sdk.Context, creator string, denom string, amount sdk.Int) error {
	assetProfile, found := k.apKeeper.GetEntry(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, creator)
	if !found {
		return sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", creator)
	}

	// Subtract the withdrawn amount from the unclaimed balance
	err := commitments.SubRewardsUnclaimed(sdk.NewCoin(denom, amount))
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	err = k.HandleWithdrawFromCommitment(ctx, &commitments, addr, withdrawCoins)
	if err != nil {
		return err
	}
	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}

func (k Keeper) HandleWithdrawFromCommitment(ctx sdk.Context, commitments *types.Commitments, addr sdk.AccAddress, amount sdk.Coins) error {
	fmt.Println("HandleWithdrawFromCommitment", commitments, addr.String(), amount.String())

	edenAmount := amount.AmountOf(ptypes.Eden)
	edenBAmount := amount.AmountOf(ptypes.EdenB)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, edenAmount))
	commitments.AddClaimed(sdk.NewCoin(ptypes.EdenB, edenBAmount))
	k.SetCommitments(ctx, *commitments)
	fmt.Println("HandleWithdrawFromCommitment1", commitments.Claimed.String())

	withdrawCoins := amount.
		Sub(sdk.NewCoin(ptypes.Eden, edenAmount)).
		Sub(sdk.NewCoin(ptypes.EdenB, edenBAmount))

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, addr.String(), withdrawCoins)

	// Send the coins to the user's account
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
	return err
}

// Withdraw validator's commission to self delegator
func (k Keeper) ProcessWithdrawValidatorCommission(ctx sdk.Context, delegator string, creator string, denom string, amount sdk.Int) error {
	assetProfile, found := k.apKeeper.GetEntry(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	commitments, err := k.DeductCommitments(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Withdraw to the delegated wallet
	addr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	commitments, _ = k.GetCommitments(ctx, delegator)
	err = k.HandleWithdrawFromCommitment(ctx, &commitments, addr, withdrawCoins)
	if err != nil {
		return err
	}

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}

// Withdraw Token - USDC
// Only withraw USDC from dexRevenue wallet
func (k Keeper) ProcessWithdrawUSDC(ctx sdk.Context, creator string, denom string, amount sdk.Int) error {
	if denom != ptypes.BaseCurrency {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	assetProfile, found := k.apKeeper.GetEntry(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	commitments, err := k.DeductCommitments(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}

// Vesting token
// Check if vesting entity count is not exceeding the maximum and if it is fine, creates a new vesting entity
// Deduct from unclaimed bucket. If it is insufficent, deduct from committed bucket as well.
func (k Keeper) ProcessTokenVesting(ctx sdk.Context, denom string, amount sdk.Int, creator string) error {
	vestingInfo, _ := k.GetVestingInfo(ctx, denom)

	if vestingInfo == nil {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", denom)
	}

	commitments, found := k.GetCommitments(ctx, creator)
	if !found {
		return sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", creator)
	}

	// Create vesting tokens entry and add to commitments
	vestingTokens := commitments.GetVestingTokens()
	if vestingInfo.NumMaxVestings <= (int64)(len(vestingTokens)) {
		return sdkerrors.Wrapf(types.ErrExceedMaxVestings, "creator: %s", creator)
	}

	commitments, err := k.DeductCommitments(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	vestingTokens = append(vestingTokens, &types.VestingTokens{
		Denom:           vestingInfo.VestingDenom,
		TotalAmount:     amount,
		UnvestedAmount:  amount,
		EpochIdentifier: vestingInfo.EpochIdentifier,
		NumEpochs:       vestingInfo.NumEpochs,
		CurrentEpoch:    0,
	})
	commitments.VestingTokens = vestingTokens

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}
