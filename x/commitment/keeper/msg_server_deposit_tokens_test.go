package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
)

func TestDepositTokens(t *testing.T) {
	t.Skip() // TODO test keeper needs bank keeper
	// Create a test context and keeper
	keeper, ctx := keepertest.CommitmentKeeper(t)
	msgServer := commitmentkeeper.NewMsgServerImpl(*keeper)

	// Create a new account
	creator, _ := sdk.AccAddressFromBech32("cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5")

	// Create a deposit message
	depositMsg := &types.MsgDepositTokens{
		Creator: creator.String(),
		Denom:   "stake",
		Amount:  sdk.NewInt(100),
	}

	// Add initial funds to creator's account
	// coins := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(200)))
	// err := keeper.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	// require.NoError(t, err)
	// // err = keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, coins)
	// require.NoError(t, err)

	// Execute the DepositTokens function
	_, err := msgServer.DepositTokens(ctx, depositMsg)
	require.NoError(t, err)

	// Check if the tokens were deposited and uncommitted balance was updated
	commitments, found := keeper.GetCommitments(ctx, depositMsg.Creator)
	require.True(t, found, "commitments not found")

	uncommittedBalance := commitments.GetUncommittedAmountForDenom(depositMsg.Denom)
	require.Equal(t, depositMsg.Amount, uncommittedBalance, "uncommitted balance did not update correctly")

	// Check if the deposited tokens were burned
	// remainingCoins := k.bankKeeper.GetBalance(ctx, creator, depositMsg.Denom)
	// require.Equal(t, sdk.NewInt(100), remainingCoins.Amount, "tokens were not burned correctly")
}
