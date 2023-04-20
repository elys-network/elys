package keeper_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestDepositTokens(t *testing.T) {
	app := app.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	msgServer := commitmentkeeper.NewMsgServerImpl(keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Create a deposit message
	depositMsg := &types.MsgDepositTokens{
		Creator: creator.String(),
		Denom:   "ueden",
		Amount:  sdk.NewInt(100),
	}

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: depositMsg.Denom, CommitEnabled: true})

	// Add initial funds to creator's account
	coins := sdk.NewCoins(sdk.NewCoin("ueden", sdk.NewInt(200)))
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, coins)
	require.NoError(t, err)
	balance := app.BankKeeper.GetBalance(ctx, creator, depositMsg.Denom)
	require.Equal(t, coins.AmountOf("ueden"), balance.Amount, "creator balance did not initialize")

	require.NoError(t, err)

	// Execute the DepositTokens function
	_, err = msgServer.DepositTokens(ctx, depositMsg)
	require.NoError(t, err)

	// Check if the tokens were deposited and uncommitted balance was updated
	commitments, found := keeper.GetCommitments(ctx, depositMsg.Creator)
	require.True(t, found, "commitments not found")

	uncommittedBalance := commitments.GetUncommittedAmountForDenom(depositMsg.Denom)
	require.Equal(t, depositMsg.Amount, uncommittedBalance, "uncommitted balance did not update correctly")

	// Check if the deposited tokens were deducted from creator balance
	remainingCoins := app.BankKeeper.GetBalance(ctx, creator, depositMsg.Denom)
	require.Equal(t, sdk.NewInt(100), remainingCoins.Amount, "tokens were not deducted correctly")

	// Check if the deposited tokens were burned
	remainingCoins = app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName), depositMsg.Denom)
	require.Equal(t, sdk.NewInt(0), remainingCoins.Amount, "tokens were not burned correctly")
}
