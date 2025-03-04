package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	// reset params
	params := m.keeper.GetParams(ctx)
	params.TakerFees = math.LegacyMustNewDecFromStr("0.001")

	takerAddr := authtypes.NewModuleAddress("taker_fee_collection")
	params.TakerFeeCollectionAddress = takerAddr.String()
	m.keeper.SetParams(ctx, params)

	return nil
}
