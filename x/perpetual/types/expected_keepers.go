package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	atypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	leveragelpmoduletypes "github.com/elys-network/elys/v6/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

	GetPoolWithAccountedBalance(ctx sdk.Context, poolId uint64) (val ammtypes.SnapshotPool)

	SwapOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot ammtypes.SnapshotPool, tokensIn sdk.Coins, tokenOutDenom string, swapFee osmomath.BigDec, weightBreakingFeePerpetualFactor osmomath.BigDec, takersFee osmomath.BigDec) (tokenOut sdk.Coin, slippage osmomath.BigDec, slippageAmount osmomath.BigDec, weightBalanceBonus, oracleOut osmomath.BigDec, swapFeeFinal osmomath.BigDec, err error)
	SwapInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot ammtypes.SnapshotPool, tokensOut sdk.Coins, tokenInDenom string, swapFee osmomath.BigDec, weightBreakingFeePerpetualFactor osmomath.BigDec, takersFee osmomath.BigDec) (tokenIn sdk.Coin, slippage, slippageAmount osmomath.BigDec, weightBalanceBonus, oracleIn osmomath.BigDec, swapFeeFinal osmomath.BigDec, err error)

	AddToPoolBalanceAndUpdateLiquidity(ctx sdk.Context, pool *ammtypes.Pool, addShares math.Int, coins sdk.Coins) error
	RemoveFromPoolBalanceAndUpdateLiquidity(ctx sdk.Context, pool *ammtypes.Pool, removeShares math.Int, coins sdk.Coins) error
	CalculateCoinsUSDValue(ctx sdk.Context, coins sdk.Coins) osmomath.BigDec
	CalculateUSDValue(ctx sdk.Context, denom string, amount math.Int) osmomath.BigDec
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
	GetAssetPrice(ctx sdk.Context, asset string) (math.LegacyDec, bool)
	GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
