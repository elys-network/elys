package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
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
	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (sdk.Coin, sdk.Dec, error)
	CalcInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, slippage sdk.Dec, err error)

	CalcSwapEstimationByDenom(
		ctx sdk.Context,
		amount sdk.Coin,
		denomIn string,
		denomOut string,
		baseCurrency string,
		discount sdk.Dec,
		overrideSwapFee sdk.Dec,
		decimals uint64,
	) (
		inRoute []*ammtypes.SwapAmountInRoute,
		outRoute []*ammtypes.SwapAmountOutRoute,
		outAmount sdk.Coin,
		spotPrice sdk.Dec,
		swapFee sdk.Dec,
		discountOut sdk.Dec,
		availableLiquidity sdk.Coin,
		slippage sdk.Dec,
		weightBonus sdk.Dec,
		priceImpact sdk.Dec,
		err error,
	)

	SwapByDenom(ctx sdk.Context, msg *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)
}
