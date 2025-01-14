package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

type LeverageLpKeeper interface {
	GetPool(ctx sdk.Context, poolId uint64) (leveragelpmoduletypes.Pool, bool)
	// Methods imported from account should be defined here
}

type AmmKeeper interface {
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool

	GetAccountedPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	SwapOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee math.LegacyDec, weightBreakingFeePerpetualFactor math.LegacyDec) (tokenOut sdk.Coin, slippage math.LegacyDec, slippageAmount math.LegacyDec, weightBalanceBonus, oracleOut math.LegacyDec, swapFeeFinal math.LegacyDec, err error)
	SwapInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee math.LegacyDec, weightBreakingFeePerpetualFactor math.LegacyDec) (tokenIn sdk.Coin, slippage, slippageAmount math.LegacyDec, weightBalanceBonus, oracleIn math.LegacyDec, swapFeeFinal math.LegacyDec, err error)

	AddToPoolBalanceAndUpdateLiquidity(ctx sdk.Context, pool *ammtypes.Pool, addShares math.Int, coins sdk.Coins) error
	RemoveFromPoolBalanceAndUpdateLiquidity(ctx sdk.Context, pool *ammtypes.Pool, removeShares math.Int, coins sdk.Coins) error
}

type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

type AssetProfileKeeper interface {
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val atypes.Entry, found bool)
	// GetEntryByDenom returns a entry from its denom value
	GetEntryByDenom(ctx sdk.Context, denom string) (val atypes.Entry, found bool)
}

type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) math.LegacyDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
