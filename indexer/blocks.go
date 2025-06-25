package indexer

import (
	"fmt"
	"github.com/elys-network/elys/v6/indexer/schema/chain"
	"log"
)

func ProcessBlock() {
	fmt.Println("-----------ProcessBlock---------")
	indexerLastHeight, err := chain.GetLatestBlockHeight()
	if err != nil {
		log.Printf("Error getting latest block height: %v", err)
		return
	}

	chainHeight := int64(currentChainHeight)

	//catchingUp := (currentChainHeight - uint64(indexerLastHeight)) >= 2
	if chainHeight > indexerLastHeight {
		for blockHeight := indexerLastHeight + 1; blockHeight <= chainHeight; blockHeight++ {
			ProcessTransactions(blockHeight)
		}
		_, err = chain.UpsertBlockHeight(chainHeight)
		if err != nil {
			log.Printf("Error upserting block %d: %v", chainHeight, err)
		}
	}
}
