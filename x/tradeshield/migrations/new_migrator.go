package migrations

import (
	"github.com/elys-network/elys/v6/x/tradeshield/keeper"
)

type Migrator struct {
	keeper keeper.Keeper
}

func NewMigrator(keeper keeper.Keeper) Migrator {
	return Migrator{keeper: keeper}
}
