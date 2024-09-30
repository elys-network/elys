package keeper_test

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/transferhook"
	"github.com/elys-network/elys/x/transferhook/types"
)

func getAmmPacketMetadata(address, action string, routes []ammtypes.SwapAmountInRoute) string {
	routesText, err := json.Marshal(routes)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`
		{
			"transferhook": {
				"receiver": "%[1]s",
				"amm": {
					"action": "%[2]s",
					"routes": %s
				} 
			}
		}`, address, action, routesText)
}

func (suite *KeeperTestSuite) TestSwapOnRecvPacket() {
	packet := channeltypes.Packet{
		Sequence:           1,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-0",
		Data:               []byte{},
		TimeoutHeight:      clienttypes.Height{},
		TimeoutTimestamp:   0,
	}
	swapRoutes := []ammtypes.SwapAmountInRoute{
		{
			PoolId:        1,
			TokenOutDenom: "uelys",
		},
	}

	usdcHostDenom := "uusdc"
	prefixedDenom := transfertypes.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), usdcHostDenom)
	usdcIbcDenom := transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()
	prefixedDenomChan1000 := transfertypes.GetPrefixedDenom(packet.GetDestPort(), "channel-1000", usdcHostDenom)
	usdcIbcDenomChan1000 := transfertypes.ParseDenomTrace(prefixedDenomChan1000).IBCDenom()

	elysDenom := "uelys"
	prefixedDenom = transfertypes.GetPrefixedDenom(packet.GetSourcePort(), packet.GetSourceChannel(), elysDenom)
	elysIbcDenom := transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()

	receiver := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	testCases := []struct {
		forwardingActive bool
		recvDenom        string
		packetData       transfertypes.FungibleTokenPacketData
		destChannel      string
		expSuccess       bool
		expSwap          bool
	}{
		{ // params not enabled - 0
			forwardingActive: false,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: getAmmPacketMetadata(receiver.String(), "Swap", swapRoutes),
				Memo:     "",
			},
			destChannel: "channel-0",
			recvDenom:   usdcIbcDenom,
			expSuccess:  false,
			expSwap:     false,
		},
		{ // elys denom - 1
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    elysIbcDenom,
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: getAmmPacketMetadata(receiver.String(), "Swap", swapRoutes),
				Memo:     "",
			},
			destChannel: "channel-0",
			recvDenom:   "uelys",
			expSuccess:  false,
			expSwap:     false,
		},
		{ // all okay - 2
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: getAmmPacketMetadata(receiver.String(), "Swap", swapRoutes),
				Memo:     "",
			},
			destChannel: "channel-0",
			recvDenom:   usdcIbcDenom,
			expSuccess:  true,
			expSwap:     true,
		},
		{ // ibc denom uusdc from different channel - 3
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: getAmmPacketMetadata(receiver.String(), "Swap", swapRoutes),
				Memo:     "",
			},
			destChannel: "channel-1000",
			recvDenom:   usdcIbcDenomChan1000,
			expSuccess:  false,
			expSwap:     false,
		},
		{ // all okay with memo Swap since ibc-go v5.1.0 - 4
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: receiver.String(),
				Memo:     getAmmPacketMetadata(receiver.String(), "Swap", swapRoutes),
			},
			destChannel: "channel-0",
			recvDenom:   usdcIbcDenom,
			expSuccess:  true,
			expSwap:     true,
		},
		{ // all okay with no functional part - 5
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: receiver.String(),
				Memo:     "",
			},
			destChannel: "channel-0",
			recvDenom:   usdcIbcDenom,
			expSuccess:  true,
			expSwap:     false,
		},
		{ // invalid elys address (receiver) - 6
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: getAmmPacketMetadata("invalid_address", "Swap", swapRoutes),
				Memo:     "",
			},
			destChannel: "channel-0",
			recvDenom:   usdcIbcDenom,
			expSuccess:  false,
			expSwap:     false,
		},
		{ // invalid elys address (memo) - 7
			forwardingActive: true,
			packetData: transfertypes.FungibleTokenPacketData{
				Denom:    "uusdc",
				Amount:   "1000000",
				Sender:   "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
				Receiver: receiver.String(),
				Memo:     getAmmPacketMetadata("invalid_address", "Swap", swapRoutes),
			},
			destChannel: "channel-0",
			recvDenom:   usdcIbcDenom,
			expSuccess:  false,
			expSwap:     false,
		},
	}

	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %d", i), func() {
			packet.DestinationChannel = tc.destChannel
			packet.Data = transfertypes.ModuleCdc.MustMarshalJSON(&tc.packetData)

			suite.SetupTest() // reset
			ctx := suite.ctx

			suite.app.TransferhookKeeper.SetParams(suite.ctx, types.Params{
				AmmActive: tc.forwardingActive,
			})
			// bootstrap accounts
			poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			poolCoins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 100000000), sdk.NewInt64Coin(usdcIbcDenom, 100000000)}.Sort()

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolCoins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolCoins)
			suite.Require().NoError(err)

			// execute function
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
				Denom:     ptypes.Elys,
				Liquidity: math.NewInt(100000000),
			})
			suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
				Denom:     usdcIbcDenom,
				Liquidity: math.NewInt(100000000),
			})

			pool := ammtypes.Pool{
				PoolId:            1,
				Address:           poolAddr.String(),
				RebalanceTreasury: treasuryAddr.String(),
				PoolParams: ammtypes.PoolParams{
					SwapFee:  math.LegacyZeroDec(),
					FeeDenom: usdcIbcDenom,
				},
				TotalShares: sdk.Coin{},
				PoolAssets: []ammtypes.PoolAsset{
					{
						Token:  poolCoins[0],
						Weight: math.NewInt(10),
					},
					{
						Token:  poolCoins[1],
						Weight: math.NewInt(10),
					},
				},
				TotalWeight: math.ZeroInt(),
			}
			suite.app.AmmKeeper.SetPool(suite.ctx, pool)

			transferIBCModule := transfer.NewIBCModule(suite.app.TransferKeeper)
			routerIBCModule := transferhook.NewIBCModule(suite.app.TransferhookKeeper, transferIBCModule)
			ack := routerIBCModule.OnRecvPacket(
				ctx,
				packet,
				receiver,
			)
			suite.app.AmmKeeper.EndBlocker(suite.ctx)

			if tc.expSuccess {
				suite.Require().True(ack.Success(), "ack should be successful - ack: %+v", string(ack.Acknowledgement()))

				// check minted balance for swap
				allBalance := suite.app.BankKeeper.GetAllBalances(ctx, receiver)
				resultBalance := suite.app.BankKeeper.GetBalance(ctx, receiver, ptypes.Elys)
				if tc.expSwap {
					suite.Require().True(resultBalance.Amount.IsPositive(), "result balance should be positive but was %s", allBalance.String())
				} else {
					suite.Require().True(resultBalance.Amount.IsZero(), "result balance should be zero but was %s", allBalance.String())
				}
			} else {
				suite.Require().False(ack.Success(), "ack should have failed - ack: %+v", string(ack.Acknowledgement()))
			}
		})
	}
}
