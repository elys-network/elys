package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
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

// AmmKeeper defines the expected interface needed to swap tokens
//
//go:generate mockery --srcpkg . --name AmmKeeper --structname AmmKeeper --filename amm_keeper.go --with-expecter
type AmmKeeper interface {
	SwapByDenom(ctx sdk.Context, msg *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)
}

// TierKeeper defines the expected interface needed to get tier information
//
//go:generate mockery --srcpkg . --name TierKeeper --structname TierKeeper --filename tier_keeper.go --with-expecter
type TierKeeper interface {
	GetMembershipTier(ctx sdk.Context, user sdk.AccAddress) (total_portfolio sdk.Dec, tier string, discount uint64)

	CalculateUSDValue(ctx sdk.Context, denom string, amount sdk.Int) sdk.Dec
}

// PerpetualKeeper defines the expected interface needed to open and close perpetual positions
//
//go:generate mockery --srcpkg . --name PerpetualKeeper --structname PerpetualKeeper --filename perpetual_keeper.go --with-expecter
type PerpetualKeeper interface {
	Open(ctx sdk.Context, msg *perpetualtypes.MsgOpen, isBroker bool) (*perpetualtypes.MsgOpenResponse, error)
	Close(ctx sdk.Context, msg *perpetualtypes.MsgClose) (*perpetualtypes.MsgCloseResponse, error)
	GetMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) (perpetualtypes.MTP, error)
	GetPool(ctx sdk.Context, poolId uint64) (val perpetualtypes.Pool, found bool)
	HandleOpenEstimation(ctx sdk.Context, req *perpetualtypes.QueryOpenEstimationRequest) (*perpetualtypes.QueryOpenEstimationResponse, error)
	HandleCloseEstimation(ctx sdk.Context, req *perpetualtypes.QueryCloseEstimationRequest) (res *perpetualtypes.QueryCloseEstimationResponse, err error)
	GetAssetPrice(ctx sdk.Context, asset string) (sdk.Dec, error)
}
