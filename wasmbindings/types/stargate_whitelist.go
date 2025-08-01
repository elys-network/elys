package types

import (
	"fmt"
	"sync"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/cosmos/gogoproto/proto"

	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
)

// stargateResponsePools keeps whitelist and its deterministic
// response binding for stargate queries.
// CONTRACT: since results of queries go into blocks, queries being added here should always be
// deterministic or can cause non-determinism in the state machine.
//
// The query is multi-threaded so we're using a sync.Pool
// to manage the allocation and de-allocation of newly created
// pb objects.
var stargateResponsePools = make(map[string]*sync.Pool)

// Note: When adding a migration here, we should also add it to the Async ICQ params in the upgrade.
// In the future we may want to find a better way to keep these in sync

//nolint:staticcheck
func init() {
	// ibc queries
	setWhitelistedQuery("/ibc.applications.transfer.v1.Query/DenomTrace", &ibctransfertypes.QueryDenomTraceResponse{})

	// cosmos-sdk queries

	// auth
	setWhitelistedQuery("/cosmos.auth.v1beta1.Query/Account", &authtypes.QueryAccountResponse{})
	setWhitelistedQuery("/cosmos.auth.v1beta1.Query/Params", &authtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.auth.v1beta1.Query/ModuleAccounts", &authtypes.QueryModuleAccountsResponse{})

	// bank
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/Balance", &banktypes.QueryBalanceResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/DenomMetadata", &banktypes.QueryDenomMetadataResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/DenomsMetadata", &banktypes.QueryDenomsMetadataResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/Params", &banktypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/SupplyOf", &banktypes.QuerySupplyOfResponse{})

	// distribution
	setWhitelistedQuery("/cosmos.distribution.v1beta1.Query/Params", &distributiontypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.distribution.v1beta1.Query/DelegatorWithdrawAddress", &distributiontypes.QueryDelegatorWithdrawAddressResponse{})
	setWhitelistedQuery("/cosmos.distribution.v1beta1.Query/ValidatorCommission", &distributiontypes.QueryValidatorCommissionResponse{})

	// gov
	setWhitelistedQuery("/cosmos.gov.v1beta1.Query/Deposit", &govtypesv1.QueryDepositResponse{})
	setWhitelistedQuery("/cosmos.gov.v1beta1.Query/Params", &govtypesv1.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.gov.v1beta1.Query/Vote", &govtypesv1.QueryVoteResponse{})

	// slashing
	setWhitelistedQuery("/cosmos.slashing.v1beta1.Query/Params", &slashingtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.slashing.v1beta1.Query/SigningInfo", &slashingtypes.QuerySigningInfoResponse{})

	// staking
	setWhitelistedQuery("/cosmos.staking.v1beta1.Query/Delegation", &stakingtypes.QueryDelegationResponse{})
	setWhitelistedQuery("/cosmos.staking.v1beta1.Query/Params", &stakingtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.staking.v1beta1.Query/Validator", &stakingtypes.QueryValidatorResponse{})

	// elys queries

	// amm
	setWhitelistedQuery("/elys.amm.Query/Params", &ammtypes.QueryParamsResponse{})
	setWhitelistedQuery("/elys.amm.Query/Pool", &ammtypes.QueryGetPoolResponse{})
	setWhitelistedQuery("/elys.amm.Query/PoolAll", &ammtypes.QueryAllPoolResponse{})
	setWhitelistedQuery("/elys.amm.Query/DenomLiquidity", &ammtypes.QueryGetDenomLiquidityResponse{})
	setWhitelistedQuery("/elys.amm.Query/DenomLiquidityAll", &ammtypes.QueryAllDenomLiquidityResponse{})
	setWhitelistedQuery("/elys.amm.Query/SwapEstimation", &ammtypes.QuerySwapEstimationResponse{})
	setWhitelistedQuery("/elys.amm.Query/SwapEstimationExactAmountOut", &ammtypes.QuerySwapEstimationExactAmountOutResponse{})
	setWhitelistedQuery("/elys.amm.Query/JoinPoolEstimation", &ammtypes.QueryJoinPoolEstimationResponse{})
	setWhitelistedQuery("/elys.amm.Query/ExitPoolEstimation", &ammtypes.QueryExitPoolEstimationResponse{})
	setWhitelistedQuery("/elys.amm.Query/SlippageTrack", &ammtypes.QuerySlippageTrackResponse{})
	setWhitelistedQuery("/elys.amm.Query/SlippageTrackAll", &ammtypes.QuerySlippageTrackAllResponse{})
	setWhitelistedQuery("/elys.amm.Query/Balance", &ammtypes.QueryBalanceResponse{})
	setWhitelistedQuery("/elys.amm.Query/InRouteByDenom", &ammtypes.QueryInRouteByDenomResponse{})
	setWhitelistedQuery("/elys.amm.Query/OutRouteByDenom", &ammtypes.QueryOutRouteByDenomResponse{})
	setWhitelistedQuery("/elys.amm.Query/SwapEstimationByDenom", &ammtypes.QuerySwapEstimationByDenomResponse{})
	setWhitelistedQuery("/elys.amm.Query/WeightAndSlippageFee", &ammtypes.QueryWeightAndSlippageFeeResponse{})
}

// IsWhitelistedQuery returns if the query is not whitelisted.
func IsWhitelistedQuery(queryPath string) error {
	_, isWhitelisted := stargateResponsePools[queryPath]
	if !isWhitelisted {
		return wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	return nil
}

// getWhitelistedQuery returns the whitelisted query at the provided path.
// If the query does not exist, or it was setup wrong by the chain, this returns an error.
// CONTRACT: must call returnStargateResponseToPool in order to avoid pointless allocs.
func getWhitelistedQuery(queryPath string) (proto.Message, error) {
	protoResponseAny, isWhitelisted := stargateResponsePools[queryPath]
	if !isWhitelisted {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	protoMarshaler, ok := protoResponseAny.Get().(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to assert type to proto.Messager")
	}
	return protoMarshaler, nil
}

type protoTypeG[T any] interface {
	*T
	proto.Message
}

// setWhitelistedQuery sets the whitelisted query at the provided path.
// This method also creates a sync.Pool for the provided protoMarshaler.
// We use generics so we can properly instantiate an object that the
// queryPath expects as a response.
func setWhitelistedQuery[T any, PT protoTypeG[T]](queryPath string, _ PT) {
	stargateResponsePools[queryPath] = &sync.Pool{
		New: func() any {
			return PT(new(T))
		},
	}
}

// returnStargateResponseToPool returns the provided protoMarshaler to the appropriate pool based on it's query path.
func returnStargateResponseToPool(queryPath string, pb proto.Message) {
	stargateResponsePools[queryPath].Put(pb)
}

func GetStargateWhitelistedPaths() (keys []string) {
	// Iterate over the map and collect the keys
	keys = make([]string, 0, len(stargateResponsePools))
	for k := range stargateResponsePools {
		keys = append(keys, k)
	}
	return keys
}
