package main

import (
	"encoding/json"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctypes "github.com/cosmos/ibc-go/v7/modules/core/types"
	accountedpooltypes "github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	burnertypes "github.com/elys-network/elys/x/burner/types"
	clocktypes "github.com/elys-network/elys/x/clock/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parametertypes "github.com/elys-network/elys/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	transferhooktypes "github.com/elys-network/elys/x/transferhook/types"
	// genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

type Genesis struct {
	GenesisTime     time.Time       `json:"genesis_time"`
	ChainID         string          `json:"chain_id"`
	InitialHeight   string          `json:"initial_height"`
	ConsensusParams ConsensusParams `json:"consensus_params"`
	AppHash         string          `json:"app_hash"`
	AppState        AppState        `json:"app_state"`
	// Include other top-level fields as needed
}

type ConsensusParams struct {
	Version   Version   `json:"version"`
	Block     Block     `json:"block"`
	Evidence  Evidence  `json:"evidence"`
	Validator Validator `json:"validator"`
}

type Version struct {
	App string `json:"app"`
}

type Validator struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

type Evidence struct {
	MaxAgeNumBlocks string `json:"max_age_num_blocks"`
	MaxAgeDuration  string `json:"max_age_duration"`
	MaxBytes        string `json:"max_bytes,omitempty"`
}

type Block struct {
	MaxBytes string `json:"max_bytes"`
	MaxGas   string `json:"max_gas"`
}

type AppState struct {
	Amm           Amm                             `json:"amm"`
	AssetProfile  AssetProfile                    `json:"assetprofile"`
	Auth          Auth                            `json:"auth"`
	AuthZ         authz.GenesisState              `json:"authz"`
	Bank          banktypes.GenesisState          `json:"bank"`
	Burner        burnertypes.GenesisState        `json:"burner"`
	Capability    Capability                      `json:"capability"`
	Clock         Clock                           `json:"clock"`
	Commitment    Commitment                      `json:"commitment"`
	Crisis        crisistypes.GenesisState        `json:"crisis"`
	Distribution  Distribution                    `json:"distribution"`
	Epochs        Epochs                          `json:"epochs"`
	Evidence      EvidenceState                   `json:"evidence"`
	Genutil       Genutil                         `json:"genutil"`
	Gov           Gov                             `json:"gov"`
	Ibc           Ibc                             `json:"ibc"`
	Incentive     Incentive                       `json:"incentive"`
	LeverageLP    LeverageLP                      `json:"leveragelp"`
	Perpetual     Perpetual                       `json:"perpetual"`
	Mint          Mint                            `json:"mint"`
	Oracle        Oracle                          `json:"oracle"`
	Parameter     parametertypes.GenesisState     `json:"parameter"`
	Params        interface{}                     `json:"params"`
	PoolAccounted accountedpooltypes.GenesisState `json:"poolaccounted"`
	Slashing      Slashing                        `json:"slashing"`
	StakeStake    StableStake                     `json:"stablestake"`
	Staking       Staking                         `json:"staking"`
	Tokenomics    Tokenomics                      `json:"tokenomics"`
	Transfer      transfertypes.GenesisState      `json:"transfer"`
	TransferHook  transferhooktypes.GenesisState  `json:"transferhook"`
	Upgrade       struct{}                        `json:"upgrade"`
	Wasm          wasmtypes.GenesisState          `json:"wasm"`
	// Include other fields as needed
}

type Tokenomics struct {
	tokenomicstypes.GenesisState

	AirdropList            []interface{}              `json:"airdrop_list"`
	GenesisInflation       TokenomicsGenesisInflation `json:"genesis_inflation"`
	TimeBasedInflationList []interface{}              `json:"time_based_inflation_list"`
}

type TokenomicsGenesisInflation struct {
	tokenomicstypes.GenesisInflation

	Inflation             TokenomicsInflationEntry `json:"inflation"`
	SeedVesting           json.Number              `json:"seed_vesting"`
	StrategicSalesVesting json.Number              `json:"strategic_sales_vesting"`
}

type TokenomicsInflationEntry struct {
	tokenomicstypes.InflationEntry

	LmRewards         json.Number `json:"lm_rewards"`
	IcsStakingRewards json.Number `json:"ics_staking_rewards"`
	CommunityFund     json.Number `json:"community_fund"`
	StrategicReserve  json.Number `json:"strategic_reserve"`
	TeamTokensVested  json.Number `json:"team_tokens_vested"`
}

type StableStake struct {
	stablestaketypes.GenesisState

	Params StableStakeParams `json:"params"`
}

type StableStakeParams struct {
	stablestaketypes.Params

	EpochLength json.Number `json:"epoch_length"`
}

type Incentive struct {
	incentivetypes.GenesisState

	Params IncentiveParams `json:"params"`
}

type IncentiveParams struct {
	incentivetypes.Params

	PoolInfos             []interface{} `json:"pool_infos"`
	ElysStakeSnapInterval json.Number   `json:"elys_stake_snap_interval"`
	DistributionInterval  json.Number   `json:"distribution_interval"`
}

type Epochs struct {
	epochstypes.GenesisState

	Epochs []interface{} `json:"epochs"`
}

type Commitment struct {
	commitmenttypes.GenesisState

	Params      CommitmentParams `json:"params"`
	Commitments []interface{}    `json:"commitments"`
}

type CommitmentParams struct {
	commitmenttypes.Params

	VestingInfos []CommitmentVestingInfo `json:"vesting_infos"`
}

type CommitmentVestingInfo struct {
	commitmenttypes.VestingInfo

	NumEpochs      json.Number `json:"num_epochs"`
	NumMaxVestings json.Number `json:"num_max_vestings"`
}

type Clock struct {
	clocktypes.GenesisState

	Params ClockParams `json:"params"`
}

type ClockParams struct {
	clocktypes.Params

	ContractGasLimit json.Number `json:"contract_gas_limit"`
}

type AssetProfile struct {
	assetprofiletypes.GenesisState

	EntryList []interface{} `json:"entry_list"`
}

type Amm struct {
	ammtypes.GenesisState

	Params         AmmParams     `json:"params"`
	PoolList       []interface{} `json:"pool_list"`
	SlippageTracks []interface{} `json:"slippage_tracks"`
}

type AmmParams struct {
	ammtypes.Params

	PoolCreationFee json.Number `json:"pool_creation_fee"`
}

type Genutil struct {
	// genutiltypes.GenesisState

	GenTxs []interface{} `json:"gen_txs"`
}

type EvidenceState struct {
	evidencetypes.GenesisState

	Evidence []interface{} `json:"evidence"`
}

type Oracle struct {
	oracletypes.GenesisState

	Params     OracleParams  `json:"params"`
	AssetInfos []interface{} `json:"asset_infos"`
	Prices     []interface{} `json:"prices"`
}

type OracleParams struct {
	oracletypes.Params

	OracleScriptID  json.Number `json:"oracle_script_id"`
	Multiplier      json.Number `json:"multiplier"`
	AskCount        json.Number `json:"ask_count"`
	MinCount        json.Number `json:"min_count"`
	PrepareGas      json.Number `json:"prepare_gas"`
	ExecuteGas      json.Number `json:"execute_gas"`
	PriceExpiryTime json.Number `json:"price_expiry_time"`
}

type Capability struct {
	capabilitytypes.GenesisState

	Index  json.Number   `json:"index"`
	Owners []interface{} `json:"owners"`
}

type Slashing struct {
	slashingtypes.GenesisState

	Params       SlashingParams `json:"params"`
	SigningInfos []interface{}  `json:"signing_infos"`
	MissedBlocks []interface{}  `json:"missed_blocks"`
}

type SlashingParams struct {
	slashingtypes.Params

	SignedBlocksWindow   json.Number `json:"signed_blocks_window"`
	DowntimeJailDuration string      `json:"downtime_jail_duration"`
}

type Mint struct {
	minttypes.GenesisState

	Params MintParams `json:"params"`
}

type MintParams struct {
	minttypes.Params

	BlocksPerYear json.Number `json:"blocks_per_year"`
}

type Gov struct {
	govtypes.GenesisState

	StartingProposalId json.Number      `json:"starting_proposal_id"`
	Deposits           []interface{}    `json:"deposits"`
	Votes              []interface{}    `json:"votes"`
	Proposals          []interface{}    `json:"proposals"`
	DepositParams      GovDepositParams `json:"deposit_params"`
	VotingParams       GovVotingParams  `json:"voting_params"`
	Params             GovParams        `json:"params"`
}

type GovParams struct {
	govtypes.Params

	MaxDepositPeriod string `json:"max_deposit_period"`
	VotingPeriod     string `json:"voting_period"`
}

type GovDepositParams struct {
	govtypes.DepositParams

	MaxDepositPeriod string `json:"max_deposit_period"`
}

type GovVotingParams struct {
	govtypes.VotingParams

	VotingPeriod string `json:"voting_period"`
}

type Staking struct {
	stakingtypes.GenesisState

	Params               StakingParams `json:"params"`
	LastValidatorPowers  []interface{} `json:"last_validator_powers"`
	Validators           []interface{} `json:"validators"`
	Delegations          []interface{} `json:"delegations"`
	UnbondingDelegations []interface{} `json:"unbonding_delegations"`
	Redelegations        []interface{} `json:"redelegations"`
}

type StakingParams struct {
	stakingtypes.Params

	UnbondingTime     string      `json:"unbonding_time"`
	MaxValidators     json.Number `json:"max_validators"`
	MaxEntries        json.Number `json:"max_entries"`
	HistoricalEntries json.Number `json:"historical_entries"`
}

type Distribution struct {
	distributiontypes.GenesisState

	DelegatorWithdrawInfos          []interface{} `json:"delegator_withdraw_infos"`
	OutstandingRewards              []interface{} `json:"outstanding_rewards"`
	ValidatorAccumulatedCommissions []interface{} `json:"validator_accumulated_commissions"`
	ValidatorHistoricalRewards      []interface{} `json:"validator_historical_rewards"`
	ValidatorCurrentRewards         []interface{} `json:"validator_current_rewards"`
	DelegatorStartingInfos          []interface{} `json:"delegator_starting_infos"`
	ValidatorSlashEvents            []interface{} `json:"validator_slash_events"`
}

type Ibc struct {
	ibctypes.GenesisState

	ClientGenesis     ClientGenesis     `json:"client_genesis"`
	ConnectionGenesis ConnectionGenesis `json:"connection_genesis"`
	ChannelGenesis    ChannelGenesis    `json:"channel_genesis"`
}

type ClientGenesis struct {
	ibcclienttypes.GenesisState

	Clients            []interface{}         `json:"clients"`
	ClientsConsensus   []interface{}         `json:"clients_consensus"`
	ClientsMetadata    []interface{}         `json:"clients_metadata"`
	Params             ibcclienttypes.Params `json:"params"`
	NextClientSequence json.Number           `json:"next_client_sequence"`
}

type ConnectionGenesis struct {
	ibcconnectiontypes.GenesisState

	Connections            []interface{}           `json:"connections"`
	ClientConnectionPaths  []interface{}           `json:"client_connection_paths"`
	NextConnectionSequence json.Number             `json:"next_connection_sequence"`
	Params                 ConnectionGenesisParams `json:"params"`
}

type ConnectionGenesisParams struct {
	ibcconnectiontypes.Params

	MaxExpectedTimePerBlock json.Number `json:"max_expected_time_per_block"`
}

type ChannelGenesis struct {
	ibcchanneltypes.GenesisState

	Channels            []interface{} `json:"channels"`
	Acknowledgements    []interface{} `json:"acknowledgements"`
	Commitments         []interface{} `json:"commitments"`
	Receipts            []interface{} `json:"receipts"`
	SendSequences       []interface{} `json:"send_sequences"`
	RecvSequences       []interface{} `json:"recv_sequences"`
	AckSequences        []interface{} `json:"ack_sequences"`
	NextChannelSequence json.Number   `json:"next_channel_sequence"`
}

type LeverageLP struct {
	leveragelptypes.GenesisState

	Params       LeverageLPParams `json:"params"`
	PoolList     []interface{}    `json:"pool_list"`
	PositionList []interface{}    `json:"position_list"`
}

type LeverageLPParams struct {
	leveragelptypes.Params

	EpochLength      json.Number `json:"epoch_length"`
	MaxOpenPositions json.Number `json:"max_open_positions"`
}

type Perpetual struct {
	perpetualtypes.GenesisState

	Params   PerpetualParams `json:"params"`
	PoolList []interface{}   `json:"pool_list"`
	MtpList  []interface{}   `json:"mtp_list"`
}

type PerpetualParams struct {
	perpetualtypes.Params

	EpochLength      json.Number `json:"epoch_length"`
	MaxOpenPositions json.Number `json:"max_open_positions"`
}

type AuthParams struct {
	authtypes.Params

	MaxMemoCharacters      json.Number `json:"max_memo_characters"`
	TxSigLimit             json.Number `json:"tx_sig_limit"`
	TxSizeCostPerByte      json.Number `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd25519   json.Number `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256K1 json.Number `json:"sig_verify_cost_secp256k1"`
}

type BaseAccount struct {
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber json.Number `json:"account_number"`
	Sequence      json.Number `json:"sequence"`
}

type ModuleAccount struct {
	BaseAccount BaseAccount `json:"base_account"`
	Name        string      `json:"name"`
	Permissions []string    `json:"permissions"`
}

type Account struct {
	*BaseAccount
	*ModuleAccount

	Type string `json:"@type"`
}

type Auth struct {
	authtypes.GenesisState

	Params   AuthParams `json:"params"`
	Accounts []Account  `json:"accounts"`
}

// KeyOutput represents the JSON structure of the output from the add key command
type KeyOutput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

// StatusOutput represents the JSON structure of the output from the status command
type StatusOutput struct {
	SyncInfo struct {
		LatestBlockHeight string `json:"latest_block_height"`
	} `json:"SyncInfo"`
}

// ProposalsOutput represents the JSON structure of the output from the query proposals command
type ProposalsOutput struct {
	Proposals []struct {
		Id string `json:"id"`
	} `json:"proposals"`
}

// Colors
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)
