package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/launchpad/keeper"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.LaunchpadKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

// TODO:
// func (k Keeper) IsEnabledToken(ctx sdk.Context, spendingToken string) bool {
// func (k Keeper) GenerateOrder(ctx sdk.Context, orderMaker string, spendingToken string, elysAmount math.Int, bonusRate sdk.Dec, price sdk.Dec) types.Purchase {
// func (k Keeper) CalcBuyElysResult(ctx sdk.Context, sender string, spendingToken string, tokenAmount math.Int) (math.Int, []types.Purchase, error) {
// func (k msgServer) BuyElys(goCtx context.Context, msg *types.MsgBuyElys) (*types.MsgBuyElysResponse, error) {
// func (k Keeper) CalcReturnElysResult(ctx sdk.Context, orderId uint64, returnElysAmount math.Int) (math.Int, error) {
// func (k msgServer) ReturnElys(goCtx context.Context, msg *types.MsgReturnElys) (*types.MsgReturnElysResponse, error) {
// func (k msgServer) DepositElysToken(goCtx context.Context, msg *types.MsgDepositElysToken) (*types.MsgDepositElysTokenResponse, error) {
// func (k msgServer) WithdrawRaised(goCtx context.Context, msg *types.MsgWithdrawRaised) (*types.MsgWithdrawRaisedResponse, error) {
