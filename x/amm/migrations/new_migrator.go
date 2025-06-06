package migrations

import (
	"github.com/elys-network/elys/v6/x/amm/keeper"
)

type Migrator struct {
	keeper keeper.Keeper
}

func NewMigrator(keeper keeper.Keeper) Migrator {
	return Migrator{keeper: keeper}
}
