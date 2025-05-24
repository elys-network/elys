package migrations

import (
	"github.com/elys-network/elys/v5/x/tokenomics/keeper"
)

type Migrator struct {
	keeper keeper.Keeper
}

func NewMigrator(keeper keeper.Keeper) Migrator {
	return Migrator{keeper: keeper}
}
