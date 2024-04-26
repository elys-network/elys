package migrations

import (
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
)

type Migrator struct {
	incentiveKeeper  incentivekeeper.Keeper
	estakingKeeper   estakingkeeper.Keeper
	masterchefKeeper masterchefkeeper.Keeper
	distrKeeper      distrkeeper.Keeper
	commitmentKeeper commitmentkeeper.Keeper
}

func NewMigrator(
	incentiveKeeper incentivekeeper.Keeper,
	estakingKeeper estakingkeeper.Keeper,
	masterchefKeeper masterchefkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	commitmentKeeper commitmentkeeper.Keeper,
) Migrator {
	return Migrator{
		incentiveKeeper:  incentiveKeeper,
		estakingKeeper:   estakingKeeper,
		masterchefKeeper: masterchefKeeper,
		distrKeeper:      distrKeeper,
		commitmentKeeper: commitmentKeeper,
	}
}
