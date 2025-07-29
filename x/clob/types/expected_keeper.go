package types

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	oracletypes "github.com/elys-network/elys/v7/x/oracle/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

type BankKeeper interface {
	GetBalance(goCtx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(goCtx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(goCtx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(goCtx context.Context, name string, amt sdk.Coins) error
	SetDenomMetaData(goCtx context.Context, denomMetaData banktypes.Metadata)
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SendCoins(goCtx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (math.LegacyDec, bool)
	GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
