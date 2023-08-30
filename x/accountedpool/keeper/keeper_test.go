package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func SetupStableCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.USDC,
		Display: "USDC",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.USDT,
		Display: "USDT",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     sdk.NewDec(100),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdk.NewDec(100),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
}
