package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/clock/keeper"
	"github.com/elys-network/elys/x/clock/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            *app.ElysApp
	bankKeeper     bankkeeper.Keeper
	queryClient    types.QueryClient
	clockMsgServer types.MsgServer
}

func (s *IntegrationTestSuite) SetupTest() {
	isCheckTx := false
	s.app = app.InitElysTestApp(true)
	s.ctx = s.app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: "testing",
		Height:  1,
		Time:    time.Now().UTC(),
	})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(s.app.ClockKeeper))

	s.queryClient = types.NewQueryClient(queryHelper)
	s.bankKeeper = s.app.BankKeeper
	s.clockMsgServer = keeper.NewMsgServerImpl(s.app.ClockKeeper)
}

func (s *IntegrationTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := s.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
