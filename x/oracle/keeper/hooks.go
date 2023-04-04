package keeper

import (
	"fmt"
	"time"

	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"

	epochstypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	params := k.GetParams(ctx)
	if epochIdentifier == params.BandEpoch {
		sourcePort := types.PortID
		channelCap, ok := k.ScopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, params.BandChannelSource))
		if !ok {
			return
		}

		assetInfos := k.GetAllAssetInfo(ctx)
		symbols := []string{}
		for _, assetInfo := range assetInfos {
			if assetInfo.BandTicker != "" {
				symbols = append(symbols, assetInfo.BandTicker)
			}
		}
		encodedCalldata := obi.MustEncode(&types.CoinRatesCallData{
			Symbols:    symbols,
			Multiplier: 1,
		})
		packetData := packet.NewOracleRequestPacketData(
			params.ClientID,
			params.OracleScriptID,
			encodedCalldata,
			params.AskCount,
			params.MinCount,
			params.FeeLimit,
			params.PrepareGas,
			params.ExecuteGas,
		)

		sequence, err := k.ChannelKeeper.SendPacket(ctx, channelCap, sourcePort, params.BandChannelSource, clienttypes.NewHeight(0, 0), uint64(ctx.BlockTime().UnixNano()+int64(10*time.Minute)), packetData.GetBytes())
		if err != nil {
			fmt.Println("epoch sequence", sequence, err)
			return
		}
	}
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {}

// Hooks wrapper struct
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
