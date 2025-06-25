package keeper

import (
	"context"
	"fmt"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/v6/x/commitment/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

// Interface declearation
type CommitmentKeeperI interface {
	// Iterate all commitments
	IterateCommitments(sdk.Context, func(types.Commitments) (stop bool))

	// Update commitment
	SetCommitments(ctx sdk.Context, commitments types.Commitments)

	// Get commitment
	GetCommitments(sdk.Context, sdk.AccAddress) types.Commitments
}

var _ CommitmentKeeperI = Keeper{}

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		hooks        types.CommitmentHooks

		accountKeeper      types.AccountKeeper
		bankKeeper         types.BankKeeper
		stakingKeeper      *stakingkeeper.Keeper
		assetProfileKeeper types.AssetProfileKeeper
		authority          string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper *stakingkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,

		accountKeeper:      accountKeeper,
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

func (k Keeper) GetAllBalances(goCtx context.Context, addr sdk.AccAddress) sdk.Coins {
	ctx := sdk.UnwrapSDKContext(goCtx)
	commitments := k.GetCommitments(ctx, addr)
	edenEdenBAmounts := sdk.Coins{}
	edenEdenBAmounts = edenEdenBAmounts.Add(sdk.NewCoin(ptypes.Eden, commitments.Claimed.AmountOf(ptypes.Eden)))
	edenEdenBAmounts = edenEdenBAmounts.Add(sdk.NewCoin(ptypes.EdenB, commitments.Claimed.AmountOf(ptypes.EdenB)))

	balances := k.bankKeeper.GetAllBalances(ctx, addr)
	return balances.Add(edenEdenBAmounts...)
}

func (k Keeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bankKeeper.SpendableCoins(ctx, addr)
}

func (k Keeper) BlockedAddr(addr sdk.AccAddress) bool {
	return k.bankKeeper.BlockedAddr(addr)
}

func (k Keeper) AddEdenEdenBOnAccount(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Coins) {
	commitments := k.GetCommitments(ctx, addr)
	var coinsChanged sdk.Coins
	if amt.AmountOf(ptypes.Eden).IsPositive() {
		coin := sdk.NewCoin(ptypes.Eden, amt.AmountOf(ptypes.Eden))
		coinsChanged.Add(coin)
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}
	if amt.AmountOf(ptypes.EdenB).IsPositive() {
		coin := sdk.NewCoin(ptypes.EdenB, amt.AmountOf(ptypes.EdenB))
		coinsChanged.Add(coin)
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}

	// Save the updated Commitments
	k.SetCommitments(ctx, commitments)

	return amt, coinsChanged
}

func (k Keeper) AddEdenEdenBOnModule(ctx sdk.Context, moduleName string, amt sdk.Coins) (sdk.Coins, sdk.Coins) {
	addr := authtypes.NewModuleAddress(moduleName)
	commitments := k.GetCommitments(ctx, addr)
	var coinsChanged sdk.Coins
	if amt.AmountOf(ptypes.Eden).IsPositive() {
		coin := sdk.NewCoin(ptypes.Eden, amt.AmountOf(ptypes.Eden))
		coinsChanged.Add(coin)
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}
	if amt.AmountOf(ptypes.EdenB).IsPositive() {
		coin := sdk.NewCoin(ptypes.EdenB, amt.AmountOf(ptypes.EdenB))
		coinsChanged.Add(coin)
		amt = amt.Sub(coin)
		commitments.AddClaimed(coin)
	}

	// Save the updated Commitments
	k.SetCommitments(ctx, commitments)

	return amt, coinsChanged
}

func (k Keeper) SubEdenEdenBOnModule(ctx sdk.Context, moduleName string, amt sdk.Coins) (sdk.Coins, sdk.Coins, error) {
	addr := authtypes.NewModuleAddress(moduleName)
	commitments := k.GetCommitments(ctx, addr)
	var coinsChanged sdk.Coins
	if amt.AmountOf(ptypes.Eden).IsPositive() {
		coin := sdk.NewCoin(ptypes.Eden, amt.AmountOf(ptypes.Eden))
		coinsChanged.Add(coin)
		amt = amt.Sub(coin)
		err := commitments.SubClaimed(coin)
		if err != nil {
			return amt, nil, err
		}
	}
	if amt.AmountOf(ptypes.EdenB).IsPositive() {
		coin := sdk.NewCoin(ptypes.EdenB, amt.AmountOf(ptypes.EdenB))
		coinsChanged.Add(coin)
		amt = amt.Sub(coin)
		err := commitments.SubClaimed(coin)
		if err != nil {
			return amt, nil, err
		}
	}

	// Save the updated Commitments
	k.SetCommitments(ctx, commitments)
	return amt, coinsChanged, nil
}

func (k Keeper) MintCoins(goCtx context.Context, moduleName string, amt sdk.Coins) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	amt, coinsChanged := k.AddEdenEdenBOnModule(ctx, moduleName, amt)

	prev := k.GetTotalSupply(ctx)
	prev.TotalEdenSupply = prev.TotalEdenSupply.Add(coinsChanged.AmountOf(ptypes.Eden))
	prev.TotalEdenbSupply = prev.TotalEdenbSupply.Add(coinsChanged.AmountOf(ptypes.EdenB))
	k.SetTotalSupply(ctx, prev)

	// Emit event to track Eden and EdenB mint amount
	if !coinsChanged.Empty() {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMintCoins,
				sdk.NewAttribute("module", moduleName),
				sdk.NewAttribute("coins", coinsChanged.String()),
			),
		)
	}

	if amt.Empty() {
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, moduleName, amt)
}

func (k Keeper) BurnCoins(goCtx context.Context, moduleName string, amt sdk.Coins) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	amt, coinsChanged, err := k.SubEdenEdenBOnModule(ctx, moduleName, amt)
	if err != nil {
		return err
	}

	prev := k.GetTotalSupply(ctx)
	prev.TotalEdenSupply = prev.TotalEdenSupply.Sub(coinsChanged.AmountOf(ptypes.Eden))
	prev.TotalEdenbSupply = prev.TotalEdenbSupply.Sub(coinsChanged.AmountOf(ptypes.EdenB))
	k.SetTotalSupply(ctx, prev)

	// Emit event to track Eden and EdenB burn amount
	if !coinsChanged.Empty() {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBurnCoins,
				sdk.NewAttribute("module", moduleName),
				sdk.NewAttribute("coins", coinsChanged.String()),
			),
		)
	}

	if amt.Empty() {
		return nil
	}

	return k.bankKeeper.BurnCoins(ctx, moduleName, amt)
}

func (k Keeper) SendCoinsFromModuleToModule(goCtx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, _, err := k.SubEdenEdenBOnModule(ctx, senderModule, amt)
	if err != nil {
		return err
	}
	amt, coinsChanged := k.AddEdenEdenBOnModule(ctx, recipientModule, amt)
	// Emit event to track Eden and EdenB send amount
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSendCoins,
			sdk.NewAttribute("sender_module", senderModule),
			sdk.NewAttribute("recipient_module", recipientModule),
			sdk.NewAttribute("coins", coinsChanged.String()),
		),
	)
	if amt.Empty() {
		return nil
	}

	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
}

func (k Keeper) SendCoinsFromModuleToAccount(goCtx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, _, err := k.SubEdenEdenBOnModule(ctx, senderModule, amt)
	if err != nil {
		return err
	}

	amt, coinsChanged := k.AddEdenEdenBOnAccount(ctx, recipientAddr, amt)
	// Emit event to track Eden and EdenB send amount
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSendCoins,
			sdk.NewAttribute("sender_module", senderModule),
			sdk.NewAttribute("recipient_address", recipientAddr.String()),
			sdk.NewAttribute("coins", coinsChanged.String()),
		),
	)
	if amt.Empty() {
		return nil
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}

func (k Keeper) SendCoinsFromAccountToModule(goCtx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(goCtx, senderAddr, recipientModule, amt)
}
