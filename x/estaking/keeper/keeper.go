package keeper

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		*stakingkeeper.Keeper
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		memKey             storetypes.StoreKey
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
	storeKey,
	memKey storetypes.StoreKey,
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
		storeKey:           storeKey,
		memKey:             memKey,
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

func (k Keeper) TotalBondedTokens(ctx sdk.Context) math.Int {
	bondedTokens := k.Keeper.TotalBondedTokens(ctx)
	edenValidator := k.GetEdenValidator(ctx)
	edenBValidator := k.GetEdenBValidator(ctx)
	bondedTokens = bondedTokens.Add(edenValidator.GetTokens()).Add(edenBValidator.GetTokens())
	return bondedTokens
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
		DelegatorShares: sdk.NewDecFromInt(totalEdenCommit),
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
	commParams := k.commKeeper.GetParams(ctx)
	totalEdenBCommit := commParams.TotalCommitted.AmountOf(ptypes.EdenB)

	return stakingtypes.Validator{
		OperatorAddress: params.EdenbCommitVal,
		ConsensusPubkey: EdenBValPubKeyAny,
		Jailed:          false,
		Status:          stakingtypes.Unbonded,
		Tokens:          totalEdenBCommit,
		DelegatorShares: sdk.NewDecFromInt(totalEdenBCommit),
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

func (k Keeper) IterateValidators(ctx sdk.Context,
	fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) {
	params := k.GetParams(ctx)
	if params.EdenCommitVal != "" {
		edenVal := k.GetEdenValidator(ctx)
		fn(0, edenVal)
	}

	if params.EdenbCommitVal != "" {
		edenBVal := k.GetEdenBValidator(ctx)
		fn(0, edenBVal)
	}
	k.Keeper.IterateValidators(ctx, fn)
}

func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) stakingtypes.DelegationI {
	params := k.GetParams(ctx)
	if addrVal.String() == params.EdenCommitVal {
		commitments := k.commKeeper.GetCommitments(ctx, addrDel)
		edenCommit := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		if edenCommit.IsZero() {
			return nil
		}
		return stakingtypes.Delegation{
			DelegatorAddress: addrDel.String(),
			ValidatorAddress: addrVal.String(),
			Shares:           sdk.NewDecFromInt(edenCommit),
		}
	}

	if addrVal.String() == params.EdenbCommitVal {
		commitments := k.commKeeper.GetCommitments(ctx, addrDel)
		edenBCommit := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
		if edenBCommit.IsZero() {
			return nil
		}
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
	commitments := k.commKeeper.GetCommitments(ctx, delegator)
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
	commParams := k.commKeeper.GetParams(ctx)

	if commParams.TotalCommitted.AmountOf(ptypes.Eden).IsPositive() {
		edenValidator := k.GetEdenValidator(ctx)
		if stop := fn(0, edenValidator); stop {
			return
		}
	}

	if commParams.TotalCommitted.AmountOf(ptypes.EdenB).IsPositive() {
		edenBValidator := k.GetEdenBValidator(ctx)
		if stop := fn(0, edenBValidator); stop {
			return
		}
	}
	k.Keeper.IterateBondedValidatorsByPower(ctx, fn)
}

func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) math.Int {
	return k.Keeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
}

func (k Keeper) SlashWithInfractionReason(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec, infraction stakingtypes.Infraction) math.Int {
	return k.Keeper.SlashWithInfractionReason(ctx, consAddr, infractionHeight, power, slashFactor, infraction)
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

	val := k.Validator(ctx, valAddr)
	if val == nil {
		return nil, errorsmod.Wrap(distrtypes.ErrNoValidatorExists, validatorAddress)
	}

	delAddr, err := sdk.AccAddressFromBech32(delegatorAddress)
	if err != nil {
		return nil, err
	}

	del := k.Delegation(ctx, delAddr, valAddr)
	if del == nil {
		return nil, distrtypes.ErrNoDelegationExists
	}

	endingPeriod := k.distrKeeper.IncrementValidatorPeriod(ctx, val)
	rewards := k.distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	return rewards, nil
}
