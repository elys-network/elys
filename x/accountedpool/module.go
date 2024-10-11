package tvl

import (
	"context"
	"encoding/json"
	"fmt"

	// this line is used by starport scaffolding # 1

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/elys-network/elys/x/accountedpool/client/cli"
	"github.com/elys-network/elys/x/accountedpool/keeper"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface that defines the independent methods a Cosmos SDK module needs to implement.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the name of the module as a string
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns a default GenesisState for the module, marshalled to json.RawMessage. The default GenesisState need to be defined by the module developer and is primarily used for testing
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root Tx command for the module. The subcommands of this root command are used by end-users to generate new transactions containing messages defined in the module
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for the module. The subcommands of this root command are used by end-users to generate new queries to the subset of the state defined by the module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the module. It should be incremented on each consensus-breaking change introduced by the module. To avoid wrong/empty versions, the initial version should be set to 1
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	currentHeight := ctx.BlockHeight()

	if currentHeight == 10314702 {
		ctx.Logger().Info("reached currentHeight == 10314702")

		// delete pools
		pools := am.keeper.GetAllAccountedPool(ctx)
		for _, pool := range pools {
			am.keeper.RemoveAccountedPool(ctx, pool.PoolId)
		}

		// - address: elys1t7z4shh8tzvjc2u9exu2fs8rmewlm6hza494x3dna0n7aumm05aq209wy9
		// pool_assets:
		// - token:
		// 	amount: "2135879098040"
		// 	denom: ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65
		//   weight: "10737418240"
		// - token:
		// 	amount: "174601908231"
		// 	denom: ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4
		//   weight: "10737418240"
		// pool_id: "2"
		// pool_params:
		//   exit_fee: "0.000000000000000000"
		//   external_liquidity_ratio: "1.801707284442911481"
		//   fee_denom: ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65
		//   swap_fee: "0.002000000000000000"
		//   threshold_weight_difference: "0.300000000000000000"
		//   use_oracle: true
		//   weight_breaking_fee_exponent: "2.500000000000000000"
		//   weight_breaking_fee_multiplier: "0.000500000000000000"
		//   weight_recovery_fee_portion: "0.100000000000000000"
		// rebalance_treasury: elys1zfz2hcvzgcg2kgw0xyc9l27nmda6qnxahjrnjrateuejc84zf2kspdd9cl
		// total_shares:
		//   amount: "4041987385104785752062188"
		//   denom: amm/pool/2
		// total_weight: "21474836480"

		totalShares, ok := sdk.NewIntFromString("4041987385104785752062188")
		if !ok {
			panic("failed to parse totalShares")
		}

		err := am.keeper.InitiateAccountedPool(ctx, ammtypes.Pool{
			PoolId:  2,
			Address: "elys1t7z4shh8tzvjc2u9exu2fs8rmewlm6hza494x3dna0n7aumm05aq209wy9",
			PoolParams: ammtypes.PoolParams{
				SwapFee:                     sdk.MustNewDecFromStr("0.002000000000000000"),
				ExitFee:                     sdk.MustNewDecFromStr("0.000000000000000000"),
				ThresholdWeightDifference:   sdk.MustNewDecFromStr("0.300000000000000000"),
				WeightBreakingFeeExponent:   sdk.MustNewDecFromStr("2.500000000000000000"),
				WeightBreakingFeeMultiplier: sdk.MustNewDecFromStr("0.000500000000000000"),
				WeightRecoveryFeePortion:    sdk.MustNewDecFromStr("0.100000000000000000"),
				UseOracle:                   true,
				ExternalLiquidityRatio:      sdk.MustNewDecFromStr("1.801707284442911481"),
				FeeDenom:                    "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			},
			TotalShares: sdk.NewCoin("amm/pool/2", totalShares),
			TotalWeight: sdk.NewInt(21474836480),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token: sdk.NewCoin(
						"ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
						sdk.NewIntFromUint64(2135879098040),
					),
					Weight: sdk.NewInt(10737418240),
				},
				{
					Token: sdk.NewCoin(
						"ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4",
						sdk.NewIntFromUint64(174601908231),
					),
					Weight: sdk.NewInt(10737418240),
				},
			},
			RebalanceTreasury: "elys1zfz2hcvzgcg2kgw0xyc9l27nmda6qnxahjrnjrateuejc84zf2kspdd9cl",
		})
		if err != nil {
			panic(err)
		}
	}
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
