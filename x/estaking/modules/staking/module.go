package staking

import (
	"context"
	"cosmossdk.io/core/appmodule"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/cometbft/cometbft/abci/types"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasServices         = AppModule{}
	_ module.HasInvariants       = AppModule{}
	_ module.HasABCIGenesis      = AppModule{}
	_ module.HasABCIEndBlock     = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/staking AppModuleBasic.
type AppModuleBasic struct {
	staking.AppModuleBasic
}

// AppModule embeds the Cosmos SDK's x/staking AppModule where we only override
// specific methods.
type AppModule struct {
	// embed the Cosmos SDK's x/staking AppModule
	staking.AppModule
	cdc        codec.Codec
	keeper     *keeper.Keeper
	accKeeper  types.AccountKeeper
	bankKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object using the native x/staking module
// AppModule constructor.
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, subspace exported.Subspace) AppModule {
	stakingAppMod := staking.NewAppModule(cdc, keeper, ak, bk, subspace)
	return AppModule{
		AppModule:  stakingAppMod,
		keeper:     keeper,
		accKeeper:  ak,
		bankKeeper: bk,
		cdc:        cdc,
	}
}

// InitGenesis delegates the InitGenesis call to the underlying x/staking module,
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	return am.keeper.InitGenesis(ctx, &genesisState)
}

// EndBlock delegates the EndBlock call to the underlying x/staking module,
func (am AppModule) EndBlock(goCtx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)
	if sdkCtx.BlockHeight() == 11517072 {
		blockTime := sdkCtx.BlockTime()
		blockHeight := sdkCtx.BlockHeight()
		unbondingValIterator, err := am.keeper.ValidatorQueueIterator(sdkCtx, blockTime, blockHeight)
		if err != nil {
			panic(err)
		}
		defer unbondingValIterator.Close()

		validatorAddressCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

		for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
			key := unbondingValIterator.Key()
			keyTime, keyHeight, err := types.ParseValidatorQueueKey(key)
			if err != nil {
				panic(err)
			}

			// All addresses for the given key have the same unbonding height and time.
			// We only unbond if the height and time are less than the current height
			// and time.
			if keyHeight <= blockHeight && (keyTime.Before(blockTime) || keyTime.Equal(blockTime)) {
				addrs := types.ValAddresses{}
				if err = am.cdc.Unmarshal(unbondingValIterator.Value(), &addrs); err != nil {
					panic(err)
				}

				for _, valAddr := range addrs.Addresses {
					addr, err := validatorAddressCodec.StringToBytes(valAddr)
					if err != nil {
						panic(err)
					}
					val, err := am.keeper.GetValidator(sdkCtx, addr)
					if err != nil {
						panic(err)
					}

					if !val.IsUnbonding() {
						val.Status = types.Unbonding
						val.Jailed = true
						err = am.keeper.SetValidator(sdkCtx, val)
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}

	return am.keeper.BlockValidatorUpdates(goCtx)
}

// BeginBlock returns the begin blocker for the staking module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.keeper.BeginBlocker(ctx)
}
