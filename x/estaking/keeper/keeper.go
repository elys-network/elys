package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		*stakingkeeper.Keeper
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		commKeeper types.CommitmentKeeper
		authority  string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	stakingKeeper *stakingkeeper.Keeper,
	commKeeper types.CommitmentKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		Keeper:     stakingKeeper,
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		commKeeper: commKeeper,
		authority:  authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) TotalBondedTokens(ctx sdk.Context) math.Int {
	bondedTokens := k.Keeper.TotalBondedTokens(ctx)
	edenValidator := k.GetEdenValidator(ctx)
	edenBValidator := k.GetEdenBValidator(ctx)
	bondedTokens = bondedTokens.Add(edenValidator.GetTokens()).Add(edenBValidator.GetTokens())
	return bondedTokens
}

func (k Keeper) GetEdenValidator(ctx sdk.Context) stakingtypes.ValidatorI {
	params := k.GetParams(ctx)
	// TODO: check potential issue
	pk1 := ed25519.GenPrivKeyFromSecret([]byte{1}).PubKey()
	pk1Any, err := codectypes.NewAnyWithValue(pk1)
	if err != nil {
		panic(err)
	}
	return stakingtypes.Validator{
		OperatorAddress: params.EdenCommitVal,
		ConsensusPubkey: pk1Any,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          sdk.ZeroInt(), // TODO: total Eden commitment
		DelegatorShares: sdk.ZeroDec(), // TODO: total eden commitment
		Description: stakingtypes.Description{
			Moniker: "EdenValidator",
		},
		UnbondingHeight:         0,
		UnbondingTime:           time.Time{},
		Commission:              stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		MinSelfDelegation:       sdk.ZeroInt(),
		UnbondingOnHoldRefCount: 0,
		UnbondingIds:            []uint64{},
	}
}

func (k Keeper) GetEdenBValidator(ctx sdk.Context) stakingtypes.ValidatorI {
	params := k.GetParams(ctx)
	// TODO: check potential issue
	pk2 := ed25519.GenPrivKeyFromSecret([]byte{2}).PubKey()
	pk2Any, err := codectypes.NewAnyWithValue(pk2)
	if err != nil {
		panic(err)
	}
	return stakingtypes.Validator{
		OperatorAddress: params.EdenbCommitVal,
		ConsensusPubkey: pk2Any,
		Jailed:          false,
		Status:          stakingtypes.Unbonded,
		Tokens:          sdk.ZeroInt(), // TODO: total EdenB commitment
		DelegatorShares: sdk.ZeroDec(), // TODO: total edenb commitment
		Description: stakingtypes.Description{
			Moniker: "EdenBValidator",
		},
		UnbondingHeight:         0,
		UnbondingTime:           time.Time{},
		Commission:              stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		MinSelfDelegation:       sdk.ZeroInt(),
		UnbondingOnHoldRefCount: 0,
		UnbondingIds:            []uint64{},
	}
}

// extended staking keeper functionalities
func (k Keeper) Validator(ctx sdk.Context, address sdk.ValAddress) stakingtypes.ValidatorI {
	params := k.GetParams(ctx)
	if address.String() == params.EdenCommitVal {
		return k.GetEdenValidator(ctx)
	}

	if address.String() == params.EdenbCommitVal {
		return k.GetEdenBValidator(ctx)
	}

	val, found := k.GetValidator(ctx, address)
	if !found {
		return nil
	}

	return val
}

func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) stakingtypes.DelegationI {
	params := k.GetParams(ctx)
	if addrVal.String() == params.EdenCommitVal {
		commitments := k.commKeeper.GetCommitments(ctx, addrDel.String())
		edenCommit := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		return stakingtypes.Delegation{
			DelegatorAddress: addrDel.String(),
			ValidatorAddress: addrVal.String(),
			Shares:           sdk.NewDecFromInt(edenCommit),
		}
	}

	if addrVal.String() == params.EdenbCommitVal {
		commitments := k.commKeeper.GetCommitments(ctx, addrDel.String())
		edenBCommit := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
		return stakingtypes.Delegation{
			DelegatorAddress: addrDel.String(),
			ValidatorAddress: addrVal.String(),
			Shares:           sdk.NewDecFromInt(edenBCommit),
		}
	}

	bond, ok := k.GetDelegation(ctx, addrDel, addrVal)
	if !ok {
		return nil
	}

	return bond
}

func (k Keeper) IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
	fn func(index int64, delegation stakingtypes.DelegationI) (stop bool)) {

	params := k.GetParams(ctx)
	commitments := k.commKeeper.GetCommitments(ctx, delegator.String())
	edenCommit := commitments.GetCommittedAmountForDenom(ptypes.Eden)
	edenBCommit := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	if edenCommit.IsPositive() {
		edenDel := stakingtypes.Delegation{
			DelegatorAddress: delegator.String(),
			ValidatorAddress: params.EdenCommitVal,
			Shares:           sdk.NewDecFromInt(edenCommit),
		}
		if stop := fn(0, edenDel); stop {
			return
		}
	}
	if edenBCommit.IsPositive() {
		edenBDel := stakingtypes.Delegation{
			DelegatorAddress: delegator.String(),
			ValidatorAddress: params.EdenbCommitVal,
			Shares:           sdk.NewDecFromInt(edenBCommit),
		}
		if stop := fn(0, edenBDel); stop {
			return
		}
	}
	k.Keeper.IterateDelegations(ctx, delegator, fn)
}

// iterate through the bonded validator set and perform the provided function
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) {
	if false { // TODO: if eden commitment positive
		edenValidator := k.GetEdenValidator(ctx)
		if stop := fn(0, edenValidator); stop {
			return
		}
	}
	if false { // TODO: if edenB commmitment positive
		edenBValidator := k.GetEdenBValidator(ctx)
		if stop := fn(0, edenBValidator); stop {
			return
		}
	}
	k.Keeper.IterateBondedValidatorsByPower(ctx, fn)
}
