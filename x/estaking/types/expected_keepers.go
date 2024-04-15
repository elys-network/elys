package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	tokenomictypes "github.com/elys-network/elys/x/tokenomics/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI            // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI // get a particular validator by consensus address

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []stakingtypes.Delegation
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator string) commitmenttypes.Commitments
	GetParams(ctx sdk.Context) commitmenttypes.Params
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
}

type DistrKeeper interface {
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
	IncrementValidatorPeriod(ctx sdk.Context, val stakingtypes.ValidatorI) uint64
	CalculateDelegationRewards(ctx sdk.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins)
}

// TokenomicsKeeper defines the expected tokenomics keeper used for simulations (noalias)
type TokenomicsKeeper interface {
	GetAllTimeBasedInflation(ctx sdk.Context) (list []tokenomictypes.TimeBasedInflation)
}

// AssetProfileKeeper defines the expected interfaces
type AssetProfileKeeper interface {
	// SetEntry set a specific entry in the store from its index
	SetEntry(ctx sdk.Context, entry assetprofiletypes.Entry)
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val assetprofiletypes.Entry, found bool)
	// GetEntryByDenom returns a entry from its denom value
	GetEntryByDenom(ctx sdk.Context, denom string) (val assetprofiletypes.Entry, found bool)
}
