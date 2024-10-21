package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	m.keeper.SetParams(ctx, types.Params{
		PoolCreationFee: math.NewInt(10_000_000),
		SlippageTrackDuration: 86400*7,
		EnableBaseCurrencyPairedPoolOnly: false,
	})
	return nil
}
