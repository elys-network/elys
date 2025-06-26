package indexer

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/elys-network/elys/v6/indexer/schema/tradeshield"
	tradeshieldmoduletypes "github.com/elys-network/elys/v6/x/tradeshield/types"
	"log"
	"strconv"
)

type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
}

var encodingConfig = EncodingConfig{}

func setEncodingConfig(ir types.InterfaceRegistry, marshaller codec.Codec, config client.TxConfig) {
	encodingConfig = EncodingConfig{
		InterfaceRegistry: ir,
		Marshaler:         marshaller,
		TxConfig:          config,
	}
}

func ProcessTransactions(blockHeight int64) {

	blockResults, err := RPCClient.BlockResults(context.Background(), &blockHeight)
	if err != nil {
		return
	}

	for _, txResult := range blockResults.TxsResults {
		if txResult.Code == 0 {
			for _, event := range txResult.Events {
				switch event.Type {
				case tradeshieldmoduletypes.TypeEvtCreatePerpetualLimitOrder:
					attributes := event.Attributes

					poolId, err := strconv.ParseInt(attributes[0].Value, 10, 64)
					if err != nil {
						panic(err)
					}

					orderId, err := strconv.ParseInt(attributes[2].Value, 10, 64)
					if err != nil {
						panic(err)
					}

					val := tradeshield.PerpetualOrder{
						OwnerAddress:     attributes[1].Value,
						PoolID:           poolId,
						OrderID:          orderId,
						OrderType:        tradeshield.OrderTypeLimitOpen,
						IsLong:           attributes[9].Value == tradeshield.LONG,
						Leverage:         attributes[6].Value,
						CollateralAmount: attributes[4].Value,
						CollateralDenom:  attributes[5].Value,
						Price:            attributes[3].Value,
						TakeProfitPrice:  attributes[7].Value,
						StopLossPrice:    attributes[8].Value,
					}

					err = tradeshield.CreatePerpetualOrder(&val)
					if err != nil {
						log.Fatal(err)
					}

				default:
				}
			}
		}
	}
}
