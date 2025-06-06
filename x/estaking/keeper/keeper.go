package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/v6/x/estaking/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

type (
	Keeper struct {
		*stakingkeeper.Keeper
		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
		parameterKeeper    types.ParameterKeeper
		commKeeper         types.CommitmentKeeper
		distrKeeper        types.DistrKeeper
		tokenomicsKeeper   types.TokenomicsKeeper
		assetProfileKeeper types.AssetProfileKeeper
		authority          string
	}
)

var (
	EdenValPubKey     cryptotypes.PubKey
	EdenBValPubKey    cryptotypes.PubKey
	EdenValPubKeyAny  *codectypes.Any
	EdenBValPubKeyAny *codectypes.Any
)

func init() {
	// validator with duplicated pubkey fails and this is unique pubKey
	EdenValPubKey := ed25519.GenPrivKeyFromSecret([]byte(ptypes.Eden)).PubKey()
	pk1Any, err := codectypes.NewAnyWithValue(EdenValPubKey)
	if err != nil {
		panic(err)
	}
	EdenValPubKeyAny = pk1Any

	EdenBValPubKey := ed25519.GenPrivKeyFromSecret([]byte(ptypes.EdenB)).PubKey()
	pk2Any, err := codectypes.NewAnyWithValue(EdenBValPubKey)
	if err != nil {
		panic(err)
	}
	EdenBValPubKeyAny = pk2Any
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	parameterKeeper types.ParameterKeeper,
	stakingKeeper *stakingkeeper.Keeper,
	commKeeper types.CommitmentKeeper,
	distrKeeper types.DistrKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	tokenomicsKeeper types.TokenomicsKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:             stakingKeeper,
		cdc:                cdc,
		storeService:       storeService,
		parameterKeeper:    parameterKeeper,
		commKeeper:         commKeeper,
		authority:          authority,
		distrKeeper:        distrKeeper,
		assetProfileKeeper: assetProfileKeeper,
		tokenomicsKeeper:   tokenomicsKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) TotalBondedTokens(goCtx context.Context) (math.Int, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bondedTokens, err := k.Keeper.TotalBondedTokens(ctx)
	if err != nil {
		return math.Int{}, err
	}
	edenValidator := k.GetEdenValidator(ctx)
	edenBValidator := k.GetEdenBValidator(ctx)
	bondedTokens = bondedTokens.Add(edenValidator.GetTokens()).Add(edenBValidator.GetTokens())
	return bondedTokens, nil
}

func (k Keeper) TotalBondedElysEdenTokens(goCtx context.Context) (math.Int, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bondedTokens, err := k.Keeper.TotalBondedTokens(ctx)
	if err != nil {
		return math.Int{}, err
	}
	edenValidator := k.GetEdenValidator(ctx)
	bondedTokens = bondedTokens.Add(edenValidator.GetTokens())
	return bondedTokens, nil
}

func (k Keeper) GetEdenValidator(ctx sdk.Context) stakingtypes.ValidatorI {
	params := k.GetParams(ctx)
	commParams := k.commKeeper.GetParams(ctx)
	totalEdenCommit := commParams.TotalCommitted.AmountOf(ptypes.Eden)

	return stakingtypes.Validator{
		OperatorAddress: params.EdenCommitVal,
		ConsensusPubkey: EdenValPubKeyAny,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          totalEdenCommit,
		DelegatorShares: math.LegacyNewDecFromInt(totalEdenCommit),
		Description: stakingtypes.Description{
			Moniker: "EdenValidator",
		},
		UnbondingHeight:         0,
		UnbondingTime:           time.Time{},
		Commission:              stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
		MinSelfDelegation:       math.ZeroInt(),
		UnbondingOnHoldRefCount: 0,
		UnbondingIds:            []uint64{},
	}
}

func (k Keeper) GetEdenBValidator(ctx sdk.Context) stakingtypes.ValidatorI {
	params := k.GetParams(ctx)
	commParams := k.commKeeper.GetParams(ctx)
	totalEdenBCommit := commParams.TotalCommitted.AmountOf(ptypes.EdenB)

	return stakingtypes.Validator{
		OperatorAddress: params.EdenbCommitVal,
		ConsensusPubkey: EdenBValPubKeyAny,
		Jailed:          false,
		Status:          stakingtypes.Unbonded,
		Tokens:          totalEdenBCommit,
		DelegatorShares: math.LegacyNewDecFromInt(totalEdenBCommit),
		Description: stakingtypes.Description{
			Moniker: "EdenBValidator",
		},
		UnbondingHeight:         0,
		UnbondingTime:           time.Time{},
		Commission:              stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
		MinSelfDelegation:       math.ZeroInt(),
		UnbondingOnHoldRefCount: 0,
		UnbondingIds:            []uint64{},
	}
}

// extended staking keeper functionalities
func (k Keeper) Validator(goCtx context.Context, address sdk.ValAddress) (stakingtypes.ValidatorI, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if address.String() == params.EdenCommitVal {
		return k.GetEdenValidator(ctx), nil
	}

	if address.String() == params.EdenbCommitVal {
		return k.GetEdenBValidator(ctx), nil
	}

	val, err := k.GetValidator(ctx, address)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (k Keeper) IterateValidators(goCtx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if params.EdenCommitVal != "" {
		edenVal := k.GetEdenValidator(ctx)
		fn(0, edenVal)
	}

	if params.EdenbCommitVal != "" {
		edenBVal := k.GetEdenBValidator(ctx)
		fn(0, edenBVal)
	}
	return k.Keeper.IterateValidators(ctx, fn)
}

func (k Keeper) Delegation(goCtx context.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) (stakingtypes.DelegationI, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if addrVal.String() == params.EdenCommitVal {
		commitments := k.commKeeper.GetCommitments(ctx, addrDel)
		edenCommit := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		if edenCommit.IsZero() {
			return nil, nil
		}
		return stakingtypes.Delegation{
			DelegatorAddress: addrDel.String(),
			ValidatorAddress: addrVal.String(),
			Shares:           math.LegacyNewDecFromInt(edenCommit),
		}, nil
	}

	if addrVal.String() == params.EdenbCommitVal {
		commitments := k.commKeeper.GetCommitments(ctx, addrDel)
		edenBCommit := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
		if edenBCommit.IsZero() {
			return nil, nil
		}
		return stakingtypes.Delegation{
			DelegatorAddress: addrDel.String(),
			ValidatorAddress: addrVal.String(),
			Shares:           math.LegacyNewDecFromInt(edenBCommit),
		}, nil
	}

	bond, err := k.GetDelegation(ctx, addrDel, addrVal)
	if err != nil {
		return nil, nil
	}

	return bond, nil
}

func (k Keeper) IterateDelegations(goCtx context.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool)) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	commitments := k.commKeeper.GetCommitments(ctx, delegator)
	edenCommit := commitments.GetCommittedAmountForDenom(ptypes.Eden)
	edenBCommit := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	if edenCommit.IsPositive() {
		edenDel := stakingtypes.Delegation{
			DelegatorAddress: delegator.String(),
			ValidatorAddress: params.EdenCommitVal,
			Shares:           math.LegacyNewDecFromInt(edenCommit),
		}
		if stop := fn(0, edenDel); stop {
			return nil
		}
	}
	if edenBCommit.IsPositive() {
		edenBDel := stakingtypes.Delegation{
			DelegatorAddress: delegator.String(),
			ValidatorAddress: params.EdenbCommitVal,
			Shares:           math.LegacyNewDecFromInt(edenBCommit),
		}
		if stop := fn(0, edenBDel); stop {
			return nil
		}
	}
	return k.Keeper.IterateDelegations(ctx, delegator, fn)
}

// iterate through the bonded validator set and perform the provided function
func (k Keeper) IterateBondedValidatorsByPower(goCtx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	commParams := k.commKeeper.GetParams(ctx)

	if commParams.TotalCommitted.AmountOf(ptypes.Eden).IsPositive() {
		edenValidator := k.GetEdenValidator(ctx)
		if stop := fn(0, edenValidator); stop {
			return nil
		}
	}

	if commParams.TotalCommitted.AmountOf(ptypes.EdenB).IsPositive() {
		edenBValidator := k.GetEdenBValidator(ctx)
		if stop := fn(0, edenBValidator); stop {
			return nil
		}
	}
	return k.Keeper.IterateBondedValidatorsByPower(ctx, fn)
}

func (k Keeper) Slash(goCtx context.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor math.LegacyDec) (math.Int, error) {
	return k.Keeper.Slash(goCtx, consAddr, infractionHeight, power, slashFactor)
}

func (k Keeper) SlashWithInfractionReason(goCtx context.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor math.LegacyDec, infraction stakingtypes.Infraction) (math.Int, error) {
	return k.Keeper.SlashWithInfractionReason(goCtx, consAddr, infractionHeight, power, slashFactor, infraction)
}

func (k Keeper) WithdrawEdenBReward(ctx sdk.Context, addr sdk.AccAddress) error {
	params := k.GetParams(ctx)
	valAddr, err := sdk.ValAddressFromBech32(params.EdenbCommitVal)
	if err != nil {
		return err
	}
	_, err = k.distrKeeper.WithdrawDelegationRewards(ctx, addr, valAddr)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) WithdrawEdenReward(ctx sdk.Context, addr sdk.AccAddress) error {
	params := k.GetParams(ctx)
	valAddr, err := sdk.ValAddressFromBech32(params.EdenCommitVal)
	if err != nil {
		return err
	}
	_, err = k.distrKeeper.WithdrawDelegationRewards(ctx, addr, valAddr)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DelegationRewards(ctx sdk.Context, delegatorAddress string, validatorAddress string) (sdk.DecCoins, error) {
	valAddr, err := sdk.ValAddressFromBech32(validatorAddress)
	if err != nil {
		return nil, err
	}

	val, err := k.Validator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	delAddr, err := sdk.AccAddressFromBech32(delegatorAddress)
	if err != nil {
		return nil, err
	}

	del, err := k.Delegation(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	endingPeriod, err := k.distrKeeper.IncrementValidatorPeriod(ctx, val)
	if err != nil {
		return nil, err
	}
	rewards, err := k.distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	if err != nil {
		return nil, err
	}
	return rewards, nil
}
