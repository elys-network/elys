package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	stakingKeeper *stakingkeeper.Keeper,
	commKeeper types.CommitmentKeeper,
	distrKeeper types.DistrKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:             stakingKeeper,
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		commKeeper:         commKeeper,
		authority:          authority,
		distrKeeper:        distrKeeper,
		assetProfileKeeper: assetProfileKeeper,
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
