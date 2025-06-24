package indexer

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"log"
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

	blockResponse, err := TxClient.GetBlockWithTxs(context.Background(), &tx.GetBlockWithTxsRequest{Height: blockHeight})
	if err != nil {
		return
	}

	blockResults, err := RPCClient.BlockResults(context.Background(), &blockHeight)
	if err != nil {
		return
	}

	for i, tx := range blockResponse.Txs {
		// Decode the raw transaction bytes into a readable format.
		//decodedTx, err := encodingConfig.TxConfig.TxDecoder()(txBytes)
		//if err != nil {
		//	log.Printf("--- TX %d: Failed to decode --- \nError: %v\n", i, err)
		//	continue
		//}

		fmt.Println("----")
		fmt.Println(tx)

		// Get the corresponding transaction result.
		txResult := blockResults.TxsResults[i]

		// Determine transaction status. Code 0 is success.
		status := "✅ SUCCESSFUL"
		if txResult.Code != 0 {
			status = "❌ FAILED"
		}

		fmt.Printf("--- Transaction %d: %s ---\n", i, status)
		fmt.Printf("  - Result Code: %d\n", txResult.Code)

		fmt.Println(txResult.Data)
		t1, err := encodingConfig.TxConfig.TxDecoder()(txResult.Data)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(t1)
		if txResult.Code != 0 {
			fmt.Printf("  - Error Log: %s\n", txResult.Log)
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
