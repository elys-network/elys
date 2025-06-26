package indexer

import (
	"github.com/cosmos/cosmos-sdk/types/tx"
	"log"

	comethttp "github.com/cometbft/cometbft/rpc/client/http"
	"google.golang.org/grpc"
)

var (
	currentChainHeight = uint64(0)
	GrpcEndpoint       = "127.0.0.1:9090"
	RPCEndPoint        = "http://127.0.0.1:26657"
)

var GrpcClient *grpc.ClientConn
var RPCClient *comethttp.HTTP
var TxClient tx.ServiceClient

func SetCurrentChainHeight(height uint64) {
	currentChainHeight = height
}

func GetCurrentChainHeight() uint64 {
	return currentChainHeight
}

func SetupChainClients() {
	var err error
	//GrpcClient, err = grpc.NewClient(GrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}

	RPCClient, err = comethttp.New(RPCEndPoint, "/websocket")
	if err != nil {
		log.Fatalf("Tendermint RPC client creation failed: %v", err)
	}

	//TxClient = tx.NewServiceClient(GrpcClient)
}
