package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Interface declearation
type CommitmentKeeperI interface {
	// Iterate all commitments
	IterateCommitments(sdk.Context, func(types.Commitments) (stop bool))

	// Update commitment
	SetCommitments(ctx sdk.Context, commitments types.Commitments)

	// Get commitment
	GetCommitments(sdk.Context, string) types.Commitments
}

var _ CommitmentKeeperI = Keeper{}

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		hooks      types.CommitmentHooks

		bankKeeper         types.BankKeeper
		stakingKeeper      types.StakingKeeper
		assetProfileKeeper types.AssetProfileKeeper
		authority          string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		bankKeeper:         bankKeeper,
		stakingKeeper:      stakingKeeper,
		assetProfileKeeper: assetProfileKeeper,
		authority:          authority,
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

func (k Keeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	commitments := k.GetCommitments(ctx, addr.String())
	edenEdenBAmounts := sdk.Coins{}
	edenEdenBAmounts = edenEdenBAmounts.Add(sdk.NewCoin(ptypes.Eden, commitments.Claimed.AmountOf(ptypes.Eden)))
	edenEdenBAmounts = edenEdenBAmounts.Add(sdk.NewCoin(ptypes.EdenB, commitments.Claimed.AmountOf(ptypes.EdenB)))

	balances := k.bankKeeper.GetAllBalances(ctx, addr)
	return balances.Add(edenEdenBAmounts...)
}

func (k Keeper) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bankKeeper.SpendableCoins(ctx, addr)
}

func (k Keeper) BlockedAddr(addr sdk.AccAddress) bool {
	return k.bankKeeper.BlockedAddr(addr)
}

func (k Keeper) AddEdenEdenBOnAccount(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Coins {
	commitments := k.GetCommitments(ctx, addr.String())
	if amt.AmountOf(ptypes.Eden).IsPositive() {
		coin := sdk.NewCoin(ptypes.Eden, amt.AmountOf(ptypes.Eden))
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}
	if amt.AmountOf(ptypes.EdenB).IsPositive() {
		coin := sdk.NewCoin(ptypes.EdenB, amt.AmountOf(ptypes.EdenB))
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}

	// Save the updated Commitments
	k.SetCommitments(ctx, commitments)
	return amt
}

func (k Keeper) AddEdenEdenBOnModule(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Coins {
	addr := authtypes.NewModuleAddress(moduleName)
	commitments := k.GetCommitments(ctx, addr.String())
	if amt.AmountOf(ptypes.Eden).IsPositive() {
		coin := sdk.NewCoin(ptypes.Eden, amt.AmountOf(ptypes.Eden))
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}
	if amt.AmountOf(ptypes.EdenB).IsPositive() {
		coin := sdk.NewCoin(ptypes.EdenB, amt.AmountOf(ptypes.EdenB))
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}

	// Save the updated Commitments
	k.SetCommitments(ctx, commitments)
	return amt
}

func (k Keeper) SubEdenEdenBOnModule(ctx sdk.Context, moduleName string, amt sdk.Coins) (sdk.Coins, error) {
	addr := authtypes.NewModuleAddress(moduleName)
	commitments := k.GetCommitments(ctx, addr.String())
	if amt.AmountOf(ptypes.Eden).IsPositive() {
		coin := sdk.NewCoin(ptypes.Eden, amt.AmountOf(ptypes.Eden))
		amt = amt.Sub(coin)
		err := commitments.SubClaimed(coin)
		if err != nil {
			return amt, err
		}
	}
	if amt.AmountOf(ptypes.EdenB).IsPositive() {
		coin := sdk.NewCoin(ptypes.EdenB, amt.AmountOf(ptypes.EdenB))
		amt = amt.Sub(coin)
		err := commitments.SubClaimed(coin)
		if err != nil {
			return amt, err
		}
	}

	// Save the updated Commitments
	k.SetCommitments(ctx, commitments)
	return amt, nil
}

func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	amt = k.AddEdenEdenBOnModule(ctx, moduleName, amt)
	if amt.Empty() {
		return nil
	}
	return k.bankKeeper.MintCoins(ctx, moduleName, amt)
}

func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	amt, err := k.SubEdenEdenBOnModule(ctx, moduleName, amt)
	if err != nil {
		return err
	}
	if amt.Empty() {
		return nil
	}

	return k.bankKeeper.BurnCoins(ctx, moduleName, amt)
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error {
	_, err := k.SubEdenEdenBOnModule(ctx, senderModule, amt)
	if err != nil {
		return err
	}
	amt = k.AddEdenEdenBOnModule(ctx, recipientModule, amt)
	if amt.Empty() {
		return nil
	}
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
}

func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	_, err := k.SubEdenEdenBOnModule(ctx, senderModule, amt)
	if err != nil {
		return err
	}

	amt = k.AddEdenEdenBOnAccount(ctx, recipientAddr, amt)
	if amt.Empty() {
		return nil
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}

func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}
