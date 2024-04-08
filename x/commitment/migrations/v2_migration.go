package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	if params.VestingInfos != nil {
		for _, vi := range params.VestingInfos {
			if vi != nil && vi.BaseDenom == ptypes.Eden {
				vi.NumBlocks = 4320 // 10s * 4320 = 12 hrs
			}
		}
	}
	m.keeper.SetParams(ctx, params)
	return nil
}
