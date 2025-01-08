package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	elystypes "github.com/elys-network/elys/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
//
//go:generate mockery --srcpkg . --name BankKeeper --structname BankKeeper --filename bank_keeper.go --with-expecter
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins

	SendCoins(goCtx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(goCtx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// AmmKeeper defines the expected interface needed to swap tokens
//
//go:generate mockery --srcpkg . --name AmmKeeper --structname AmmKeeper --filename amm_keeper.go --with-expecter
type AmmKeeper interface {
	SwapByDenom(ctx sdk.Context, msg *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)
	CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) elystypes.Dec34
	CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) elystypes.Dec34
}

// PerpetualKeeper defines the expected interface needed to open and close perpetual positions
//
//go:generate mockery --srcpkg . --name PerpetualKeeper --structname PerpetualKeeper --filename perpetual_keeper.go --with-expecter
type PerpetualKeeper interface {
	Open(ctx sdk.Context, msg *perpetualtypes.MsgOpen) (*perpetualtypes.MsgOpenResponse, error)
	Close(ctx sdk.Context, msg *perpetualtypes.MsgClose) (*perpetualtypes.MsgCloseResponse, error)
	GetMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) (perpetualtypes.MTP, error)
	GetPool(ctx sdk.Context, poolId uint64) (val perpetualtypes.Pool, found bool)
	GetParams(ctx sdk.Context) perpetualtypes.Params
	HandleOpenEstimation(ctx sdk.Context, req *perpetualtypes.QueryOpenEstimationRequest) (*perpetualtypes.QueryOpenEstimationResponse, error)
	HandleCloseEstimation(ctx sdk.Context, req *perpetualtypes.QueryCloseEstimationRequest) (res *perpetualtypes.QueryCloseEstimationResponse, err error)
	GetAssetPrice(ctx sdk.Context, asset string) (sdkmath.LegacyDec, error)
	GetMTPsForAddressWithPagination(ctx sdk.Context, mtpAddress sdk.AccAddress, pagination *query.PageRequest) ([]*perpetualtypes.MtpAndPrice, *query.PageResponse, error)
}
