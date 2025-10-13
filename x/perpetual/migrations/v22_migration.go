package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Need to do this for testnet upgrade testing as some features is being removed from testnet
func (m Migrator) V22Migration(_ sdk.Context) error {
	return nil
}
