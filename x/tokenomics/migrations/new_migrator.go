package migrations

import (
	"github.com/elys-network/elys/v4/x/tokenomics/keeper"
)

type Migrator struct {
	keeper keeper.Keeper
}

func NewMigrator(keeper keeper.Keeper) Migrator {
	return Migrator{keeper: keeper}
}
