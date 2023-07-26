package app

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/elys-network/elys/x/commitment/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

const Bech32Prefix = "elys"

// Initiate a new ElysApp object - Common function used by the following 2 functions.
func InitiateNewElysApp(opts ...wasm.Option) *ElysApp {
	db := dbm.NewMemDB()
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue
	nodeHome := ""
	appOptions[flags.FlagHome] = nodeHome // ensure unique folder
	appOptions[server.FlagInvCheckPeriod] = 1
	app := NewElysApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		wasmtypes.EnableAllProposals,
		appOptions,
		opts,
	)

	return app
}

// Initializes a new ElysApp without IBC functionality
func InitElysTestApp(initChain bool) *ElysApp {
	app := InitiateNewElysApp()
	if initChain {
		genesisState, valSet, _, _ := GenesisStateWithValSet(app)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simtestutil.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)

		// commit genesis changes
		app.Commit()
		app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
			Height:             app.LastBlockHeight() + 1,
			AppHash:            app.LastCommitID().Hash,
			ValidatorsHash:     valSet.Hash(),
			NextValidatorsHash: valSet.Hash(),
		}})
	}

	return app
}

// Initializes a new ElysApp without IBC functionality and returns genesis account (delegator)
func InitElysTestAppWithGenAccount() (*ElysApp, sdk.AccAddress, sdk.ValAddress) {
	app := InitiateNewElysApp()

	genesisState, valSet, genAcount, valAddress := GenesisStateWithValSet(app)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return app, genAcount, valAddress
}

func GenesisStateWithValSet(app *ElysApp) (GenesisState, *tmtypes.ValidatorSet, sdk.AccAddress, sdk.ValAddress) {
	privVal := mock.NewPV()
	pubKey, _ := privVal.GetPubKey()
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	senderPrivKey.PubKey().Address()
	acc := authtypes.NewBaseAccountWithAddress(senderPrivKey.PubKey().Address().Bytes())
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100000000000000))),
	}

	//////////////////////
	balances := []banktypes.Balance{balance}
	genesisState := NewDefaultGenesisState(app.AppCodec())
	genAccs := []authtypes.GenesisAccount{acc}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, _ := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		pkAny, _ := codectypes.NewAnyWithValue(pk)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 2), sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(10, 2)),
			MinSelfDelegation: sdkmath.OneInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

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
	return genesisState, valSet, genAccs[0].GetAddress(), sdk.ValAddress(validator.Address)
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *ElysApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

func addTestAddrs(app *ElysApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	bondDenom := app.StakingKeeper.BondDenom(ctx)
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
func AddTestCommitment(app *ElysApp, ctx sdk.Context, address sdk.AccAddress, committed []sdk.Coins, uncommitted []sdk.Coins) {
	commitment, found := app.CommitmentKeeper.GetCommitments(ctx, address.String())
	if !found {
		commitment = ctypes.Commitments{
			Creator:           address.String(),
			CommittedTokens:   []*types.CommittedTokens{},
			UncommittedTokens: []*types.UncommittedTokens{},
		}
	}

	// Loop uncommitted tokens
	for _, uc := range uncommitted {
		Denom := uc.GetDenomByIndex(0)

		// Get the uncommitted tokens for the creator
		uncommittedToken, _ := commitment.GetUncommittedTokensForDenom(Denom)
		if !found {
			uncommittedTokens := commitment.GetUncommittedTokens()
			uncommittedToken = &types.UncommittedTokens{
				Denom:  Denom,
				Amount: sdk.ZeroInt(),
			}
			uncommittedTokens = append(uncommittedTokens, uncommittedToken)
			commitment.UncommittedTokens = uncommittedTokens
		}
		// Update the uncommitted tokens amount
		uncommittedToken.Amount = uncommittedToken.Amount.Add(uc.AmountOf(Denom))
	}

	for _, c := range committed {
		Denom := c.GetDenomByIndex(0)

		// Get the uncommitted tokens for the creator
		committedToken, _ := commitment.GetCommittedTokensForDenom(Denom)
		if !found {
			committedTokens := commitment.GetCommittedTokens()
			committedToken = &types.CommittedTokens{
				Denom:  Denom,
				Amount: sdk.ZeroInt(),
			}
			committedTokens = append(committedTokens, committedToken)
			commitment.CommittedTokens = committedTokens
		}
		// Update the uncommitted tokens amount
		committedToken.Amount = committedToken.Amount.Add(c.AmountOf(Denom))
	}

	app.CommitmentKeeper.SetCommitments(ctx, commitment)
}
