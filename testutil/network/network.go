package network

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	pruningtypes "cosmossdk.io/store/pruning/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	cometbftrand "github.com/cometbft/cometbft/libs/rand"
	cosmosdb "github.com/cosmos/cosmos-db"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/app"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type (
	Network = network.Network
	Config  = network.Config
)

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t *testing.T, configs ...network.Config) *network.Network {
	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = DefaultConfig(t.TempDir())
	} else {
		cfg = configs[0]
	}
	net, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err)
	_, err = net.WaitForHeight(1)
	require.NoError(t, err)
	t.Cleanup(net.Cleanup)
	return net
}

// DefaultConfig will initialize config for the network with custom application,
// genesis and single validator. All other parameters are inherited from cosmos-sdk/testutil/network.DefaultConfig
func DefaultConfig(tempDirectory string) network.Config {

	fmt.Println("EPA")

	var (
		encoding = app.MakeEncodingConfig()
		chainId  = "elys-" + cometbftrand.NewRand().Str(6)
	)

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = tempDirectory

	tempApplication := app.NewElysApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		tempDirectory,
		appOptions,
		[]wasmkeeper.Option{},
	)

	return network.Config{
		Codec:             tempApplication.AppCodec(),
		TxConfig:          tempApplication.TxConfig(),
		LegacyAmino:       tempApplication.LegacyAmino(),
		InterfaceRegistry: tempApplication.InterfaceRegistry(),
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.ValidatorI) servertypes.Application {

			tempDirectory := tempDirectory + uuid.New().String()

			return app.NewElysApp(
				val.GetCtx().Logger, cosmosdb.NewMemDB(), nil, true,
				map[int64]bool{},
				tempDirectory,
				appOptions,
				[]wasm.Option{},
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
				baseapp.SetChainID(chainId),
			)
		},
		GenesisState:    app.NewDefaultGenesisState(tempApplication, encoding.Marshaler),
		TimeoutCommit:   2 * time.Second,
		ChainID:         chainId,
		NumValidators:   1,
		BondDenom:       sdk.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: pruningtypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}
}
