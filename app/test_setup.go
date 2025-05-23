package app

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcCommitment "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	consumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/v6/x/ccv/types"
	atypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/v5/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/v5/x/stablestake/types"
)

// Initiate a new ElysApp object - Common function used by the following 2 functions.
func InitiateNewElysApp(t *testing.T) *ElysApp {
	db := dbm.NewMemDB()
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = 5
	nodeHome := ""
	appOptions[flags.FlagHome] = nodeHome // ensure unique folder
	appOptions[server.FlagInvCheckPeriod] = 1
	app := NewElysApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		t.TempDir(),
		appOptions,
		[]wasmkeeper.Option{},
	)

	return app
}

// Initializes a new ElysApp without IBC functionality
func InitElysTestApp(initChain bool, t *testing.T) *ElysApp {
	app := InitiateNewElysApp(t)
	if initChain {
		genesisState, valSet, _, _ := GenesisStateWithValSet(app)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		_, err = app.InitChain(&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		})
		if err != nil {
			panic(err)
		}

		// commit genesis changes
		_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
			Height:             app.LastBlockHeight() + 1,
			Hash:               app.LastCommitID().Hash,
			NextValidatorsHash: valSet.Hash(),
		})
		if err != nil {
			panic(err)
		}
		_, err = app.BeginBlocker(app.BaseApp.NewContext(initChain))
		if err != nil {
			panic(err)
		}
	}

	return app
}

// Initializes a new ElysApp without IBC functionality and returns genesis account (delegator)
func InitElysTestAppWithGenAccount(t *testing.T) (*ElysApp, sdk.AccAddress, sdk.ValAddress) {
	app := InitiateNewElysApp(t)

	genesisState, valSet, genAcount, valAddress := GenesisStateWithValSet(app)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = app.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	if err != nil {
		panic(err)
	}

	// commit genesis changes
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             app.LastBlockHeight() + 1,
		Hash:               app.LastCommitID().Hash,
		NextValidatorsHash: valSet.Hash(),
	})
	if err != nil {
		panic(err)
	}
	_, err = app.BeginBlocker(app.BaseApp.NewContext(true))
	if err != nil {
		panic(err)
	}

	return app, genAcount, valAddress
}

func GenesisStateWithValSet(app *ElysApp) (GenesisState, *cmttypes.ValidatorSet, sdk.AccAddress, sdk.ValAddress) {
	privVal := mock.NewPV()
	pubKey, _ := privVal.GetPubKey()
	validator := cmttypes.NewValidator(pubKey, 1)
	valSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	senderPrivKey.PubKey().Address()
	acc := authtypes.NewBaseAccountWithAddress(senderPrivKey.PubKey().Address().Bytes())
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(ptypes.Elys, math.NewInt(100000000000000))),
	}

	//////////////////////
	balances := []banktypes.Balance{balance}
	genesisState := NewDefaultGenesisState(app, app.AppCodec())
	genAP := atypes.DefaultGenesis()
	genAP.EntryList = []atypes.Entry{
		{
			BaseDenom: ptypes.BaseCurrency,
			Denom:     ptypes.BaseCurrency,
		},
	}
	genesisState[atypes.ModuleName] = app.AppCodec().MustMarshalJSON(genAP)
	genAccs := []authtypes.GenesisAccount{acc}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction
	initValPowers := []abci.ValidatorUpdate{}

	for _, val := range valSet.Validators {
		pk, _ := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		pkAny, _ := codectypes.NewAnyWithValue(pk)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyNewDecWithPrec(5, 2), math.LegacyNewDecWithPrec(10, 2), math.LegacyNewDecWithPrec(10, 2)),
			MinSelfDelegation: math.OneInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), math.LegacyOneDec()))

		// add initial validator powers so consumer InitGenesis runs correctly
		pub, _ := val.ToProto()
		initValPowers = append(initValPowers, abci.ValidatorUpdate{
			Power:  val.VotingPower,
			PubKey: pub.PubKey,
		})
	}

	// set validators and delegations
	params := stakingtypes.DefaultParams()
	params.BondDenom = ptypes.Elys

	stakingGenesis := stakingtypes.NewGenesisState(params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(ptypes.Elys, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(ptypes.Elys, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	vals, err := cmttypes.PB2TM.ValidatorUpdates(initValPowers)
	if err != nil {
		panic("failed to get vals")
	}

	consumerGenesisState := CreateMinimalConsumerTestGenesis()
	consumerGenesisState.Provider.InitialValSet = initValPowers
	consumerGenesisState.Provider.ConsensusState.NextValidatorsHash = cmttypes.NewValidatorSet(vals).Hash()
	consumerGenesisState.Params.Enabled = true
	genesisState[consumertypes.ModuleName] = app.AppCodec().MustMarshalJSON(consumerGenesisState)

	return genesisState, valSet, genAccs[0].GetAddress(), sdk.ValAddress(validator.Address)
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// CreateRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func CreateRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *ElysApp, ctx sdk.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, CreateRandomAccounts)
}

func SetStakingParam(app *ElysApp, ctx sdk.Context) error {
	return app.StakingKeeper.SetParams(ctx, stakingtypes.Params{
		UnbondingTime:     1209600,
		MaxValidators:     60,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		BondDenom:         "uelys",
		MinCommissionRate: math.LegacyNewDec(0),
	})
}

func SetPerpetualParams(app *ElysApp, ctx sdk.Context) {
	app.PerpetualKeeper.SetParams(ctx, &perpetualtypes.DefaultGenesis().Params)
}

func SetMasterChefParams(app *ElysApp, ctx sdk.Context) {
	app.MasterchefKeeper.SetParams(ctx, types.DefaultGenesis().Params)
}

func SetStableStake(app *ElysApp, ctx sdk.Context) {
	app.StablestakeKeeper.SetParams(ctx, stablestaketypes.DefaultGenesis().Params)
}

func SetParameters(app *ElysApp, ctx sdk.Context) {
	app.ParameterKeeper.SetParams(ctx, ptypes.DefaultGenesis().Params)
}

func SetupAssetProfile(app *ElysApp, ctx sdk.Context) {

	app.AssetprofileKeeper.SetEntry(ctx, atypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "uusdc",
		Path:                     "transfer/channel-12",
		IbcChannelId:             "channel-12",
		IbcCounterpartyChannelId: "channel-19",
		DisplayName:              "USDC",
		DisplaySymbol:            "uUSDC",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "uUSDC",
		TransferLimit:            "",
		Permissions:              []string{},
		UnitDenom:                "uusdc",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	})
}

func addTestAddrs(app *ElysApp, ctx sdk.Context, accNum int, accAmt math.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	bondDenom, _ := app.StakingKeeper.BondDenom(ctx)
	initCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func initAccountWithCoins(app *ElysApp, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

// Add testing commitments
func AddTestCommitment(app *ElysApp, ctx sdk.Context, address sdk.AccAddress, committed sdk.Coins) {
	commitment := app.CommitmentKeeper.GetCommitments(ctx, address)

	for _, c := range committed {
		commitment.AddCommittedTokens(c.Denom, c.Amount, uint64(ctx.BlockTime().Unix()))
	}

	app.CommitmentKeeper.SetCommitments(ctx, commitment)
}

func CreateMinimalConsumerTestGenesis() *ccvtypes.ConsumerGenesisState {
	genesisState := ccvtypes.DefaultConsumerGenesisState()
	genesisState.Params.Enabled = true
	genesisState.NewChain = true
	genesisState.Provider.ClientState = ccvprovidertypes.DefaultParams().TemplateClient
	genesisState.Provider.ClientState.ChainId = "stride"
	genesisState.Provider.ClientState.LatestHeight = ibctypes.Height{RevisionNumber: 0, RevisionHeight: 1}
	trustPeriod, err := ccvtypes.CalculateTrustPeriod(genesisState.Params.UnbondingPeriod, ccvprovidertypes.DefaultTrustingPeriodFraction)
	if err != nil {
		panic("provider client trusting period error")
	}
	genesisState.Provider.ClientState.TrustingPeriod = trustPeriod
	genesisState.Provider.ClientState.UnbondingPeriod = genesisState.Params.UnbondingPeriod
	genesisState.Provider.ClientState.MaxClockDrift = ccvprovidertypes.DefaultMaxClockDrift
	genesisState.Provider.ConsensusState = &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
		Root:      ibcCommitment.MerkleRoot{Hash: []byte("dummy")},
	}

	return genesisState
}
