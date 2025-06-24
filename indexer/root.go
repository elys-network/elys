package indexer

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/elys-network/elys/v6/indexer/db"
	"log"
	"time"
)

func Start(ir types.InterfaceRegistry, marshaller codec.Codec, config client.TxConfig) {

	setEncodingConfig(ir, marshaller, config)
	db.Connect()
	//defer db.Close()

	// Check if the connection is actually alive.
	err := db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	SetupChainClients()

	for {
		ProcessBlock()
		time.Sleep(200 * time.Millisecond) // Sleep for 2 seconds
	}
	//defer GrpcClient.Close()
}
