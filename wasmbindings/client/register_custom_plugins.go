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
	incentiveclientwasm "github.com/elys-network/elys/x/incentive/client/wasm"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	marginclientwasm "github.com/elys-network/elys/x/margin/client/wasm"
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
	margin *marginkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	parameter *parameterkeeper.Keeper,
	stablestake *stablestakekeeper.Keeper,
	staking *stakingkeeper.Keeper,
	tokenomics *tokenomicskeeper.Keeper,
	transferhook *transferhookkeeper.Keeper,
) []wasmkeeper.Option {
	ammQuerier := ammclientwasm.NewQuerier(amm, bank, commitment)
	ammMessenger := ammclientwasm.NewMessenger(amm)

	commitmentQuerier := commitmentclientwasm.NewQuerier(commitment)
	commitmentMessenger := commitmentclientwasm.NewMessenger(commitment, staking)

	incentiveQuerier := incentiveclientwasm.NewQuerier(incentive)
	incentiveMessenger := incentiveclientwasm.NewMessenger(incentive, staking, commitment)

	marginQuerier := marginclientwasm.NewQuerier(margin)
	marginMessenger := marginclientwasm.NewMessenger(margin)

	oracleQuerier := oracleclientwasm.NewQuerier(oracle)
	oracleMessenger := oracleclientwasm.NewMessenger(oracle)

	moduleQueriers := []types.ModuleQuerier{
		ammQuerier,
		commitmentQuerier,
		incentiveQuerier,
		marginQuerier,
		oracleQuerier,
	}

	wasmQueryPlugin := types.NewQueryPlugin(moduleQueriers, amm, oracle, bank, staking, commitment, margin, incentive)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: types.CustomQuerier(wasmQueryPlugin),
	})

	moduleMessengers := []types.ModuleMessenger{
		ammMessenger,
		commitmentMessenger,
		incentiveMessenger,
		marginMessenger,
		oracleMessenger,
	}

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		types.CustomMessageDecorator(moduleMessengers, amm, margin, staking, commitment, incentive),
	)
	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
