package indexer

import (
	"context"
	"fmt"
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
		// Decode the raw transaction bytes into a readable format.
		//decodedTx, err := encodingConfig.TxConfig.TxDecoder()(txBytes)
		//if err != nil {
		//	log.Printf("--- TX %d: Failed to decode --- \nError: %v\n", i, err)
		//	continue
		//}

		fmt.Println("TX RESULTS:")
		fmt.Println(txResult)

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
						OrderType:        0,
						IsLong:           true,
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

		//for _, msg := range t1.GetMsgs() {
		//	fmt.Println(msg)
		//}

		// Print the transaction body (which contains the messages).
		// We can cast the decoded tx to the `Tx` interface which has `GetMsgs`
		//txWithBody, ok := decodedTx.(tx.Tx)
		//if !ok {
		//	log.Println("  - Could not cast to tx.Tx to get body")
		//	continue
		//}
		//
		//fmt.Println("  - Transaction Body Messages:")
		//for msgIndex, msg := range txWithBody.GetMsgs() {
		//	// Print the message content. Proto-based messages have a good default String() representation.
		//	fmt.Printf("    - Msg %d: %s\n", msgIndex, msg.String())
		//}
		//fmt.Println()
	}
}
