package client

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/elys-network/elys/wasmbindings/types"
	accountedpoolkeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	ammclientwasm "github.com/elys-network/elys/x/amm/client/wasm"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	assetprofilekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	burnerkeeper "github.com/elys-network/elys/x/burner/keeper"
	clockkeeper "github.com/elys-network/elys/x/clock/keeper"
	commitmentclientwasm "github.com/elys-network/elys/x/commitment/client/wasm"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	epochskeeper "github.com/elys-network/elys/x/epochs/keeper"
	incentiveclientwasmmessenger "github.com/elys-network/elys/x/incentive/client/wasm/messenger"
	incentiveclientwasmquerier "github.com/elys-network/elys/x/incentive/client/wasm/querier"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	liquidityproviderkeeper "github.com/elys-network/elys/x/liquidityprovider/keeper"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	oracleclientwasm "github.com/elys-network/elys/x/oracle/client/wasm"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
)

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
	accountedPool *accountedpoolkeeper.Keeper,
	amm *ammkeeper.Keeper,
	assetprofile *assetprofilekeeper.Keeper,
	burner *burnerkeeper.Keeper,
	clock *clockkeeper.Keeper,
	commitment *commitmentkeeper.Keeper,
	epochs *epochskeeper.Keeper,
	incentive *incentivekeeper.Keeper,
	leveragelp *leveragelpkeeper.Keeper,
	liquidityprovider *liquidityproviderkeeper.Keeper,
	margin *marginkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	parameter *parameterkeeper.Keeper,
	stablestake *stablestakekeeper.Keeper,
	staking *stakingkeeper.Keeper,
	tokenomics *tokenomicskeeper.Keeper,
	transferhook *transferhookkeeper.Keeper,
) []wasmkeeper.Option {
	oracleQuerier := oracleclientwasm.NewQuerier(oracle)
	oracleMessenger := oracleclientwasm.NewMessenger(oracle)

	ammQuerier := ammclientwasm.NewQuerier(amm, bank, commitment)
	ammMessenger := ammclientwasm.NewMessenger(amm)

	commitmentQuerier := commitmentclientwasm.NewQuerier(commitment)
	commitmentMessenger := commitmentclientwasm.NewMessenger(commitment, staking)

	incentiveQuerier := incentiveclientwasmquerier.NewQuerier(incentive)
	incentiveMessenger := incentiveclientwasmmessenger.NewMessenger(incentive, staking, commitment, incentive)

	moduleQueriers := []types.ModuleQuerier{
		oracleQuerier,
		ammQuerier,
		commitmentQuerier,
		incentiveQuerier,
	}

	wasmQueryPlugin := types.NewQueryPlugin(moduleQueriers, amm, oracle, bank, staking, commitment, margin, incentive)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: types.CustomQuerier(wasmQueryPlugin),
	})

	moduleMessengers := []types.ModuleMessenger{
		oracleMessenger,
		ammMessenger,
		commitmentMessenger,
		incentiveMessenger,
	}

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		types.CustomMessageDecorator(moduleMessengers, amm, margin, staking, commitment, incentive),
	)
	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
