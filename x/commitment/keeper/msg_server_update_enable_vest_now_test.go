package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/elys-network/elys/v4/app"
	commitmentkeeper "github.com/elys-network/elys/v4/x/commitment/keeper"
	"github.com/elys-network/elys/v4/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateEnableVestNow(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	govAddress := sdk.AccAddress(address.Module("gov"))

	// Define the test data
	signer := govAddress.String()

	msg := types.MsgUpdateEnableVestNow{
		Authority:     signer,
		EnableVestNow: false,
	}
	_, err := msgServer.UpdateEnableVestNow(ctx, &msg)
	require.NoError(t, err)
}
