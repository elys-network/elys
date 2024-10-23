package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

type LeverageLpKeeper interface {
	GetPool(ctx sdk.Context, poolId uint64) (leveragelpmoduletypes.Pool, bool)
	// Methods imported from account should be defined here
}

type AmmKeeper interface {
	// Get first pool id that contains all denoms in pool assets
	GetBestPoolWithDenoms(ctx sdk.Context, denoms []string, usesOracle bool) (pool ammtypes.Pool, found bool)
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	GetPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (sdk.Coin, sdk.Dec, error)
	CalcInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, slippage sdk.Dec, err error)

	AddToPoolBalance(ctx sdk.Context, pool *ammtypes.Pool, addShares math.Int, coins sdk.Coins) error
	RemoveFromPoolBalance(ctx sdk.Context, pool *ammtypes.Pool, removeShares math.Int, coins sdk.Coins) error
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

type AssetProfileKeeper interface {
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val atypes.Entry, found bool)
	// GetEntryByDenom returns a entry from its denom value
	GetEntryByDenom(ctx sdk.Context, denom string) (val atypes.Entry, found bool)
}

type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdk.Dec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
