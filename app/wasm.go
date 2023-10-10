package app

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtype "github.com/elys-network/elys/x/amm/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
	}
}

type QueryPlugin struct {
	oracleKeeper *oraclekeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	oracle *oraclekeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		oracleKeeper: oracle,
	}
}

func RegisterCustomPlugins(
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(oracle)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(amm),
	)
	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}

// CustomQuerier dispatches custom CosmWasm bindings queries.
func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery ElysQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "elys query")
		}

		switch {
		case contractQuery.PriceAll != nil:
			pagination := contractQuery.PriceAll.Pagination

			// Calling the PriceAll function and handling its response
			priceResponse, err := qp.oracleKeeper.PriceAll(ctx, &oracletypes.QueryAllPriceRequest{Pagination: pagination})
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to get all prices")
			}

			// copy array priceResponse.Price
			price := make([]oracletypes.Price, len(priceResponse.Price))
			copy(price, priceResponse.Price)

			res := AllPriceResponse{
				Price: price,
				Pagination: &query.PageResponse{
					NextKey: priceResponse.Pagination.NextKey,
				},
			}

			// Serializing the response to a JSON byte array
			responseBytes, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to serialize price response")
			}

			return responseBytes, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown elys query variant"}
		}
	}
}

type ElysQuery struct {
	PriceAll *PriceAll `json:"price_all,omitempty"`
}

type PriceAll struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type AllPriceResponse struct {
	Price      []oracletypes.Price `protobuf:"bytes,1,rep,name=price,proto3" json:"price"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func CustomMessageDecorator(amm *ammkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			amm:     amm,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	amm     *ammkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg ElysMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "osmosis msg")
		}
		if contractMsg.MsgSwapExactAmountIn != nil {
			return m.msgSwapExactAmountIn(ctx, contractAddr, contractMsg.MsgSwapExactAmountIn)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) msgSwapExactAmountIn(ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *MsgSwapExactAmountIn) ([]sdk.Event, [][]byte, error) {
	err := PerformMsgSwapExactAmountIn(m.amm, ctx, contractAddr, msgSwapExactAmountIn)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform swap")
	}
	return nil, nil, nil
}

func PerformMsgSwapExactAmountIn(f *ammkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *MsgSwapExactAmountIn) error {
	if msgSwapExactAmountIn == nil {
		return wasmvmtypes.InvalidRequest{Err: "swap null swap"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*f)

	var PoolIds []uint64
	var TokenOutDenoms []string

	for _, route := range msgSwapExactAmountIn.Routes {
		PoolIds = append(PoolIds, route.PoolId)
		TokenOutDenoms = append(TokenOutDenoms, route.TokenOutDenom)
	}

	msgMsgSwapExactAmountIn := ammtype.NewMsgSwapExactAmountIn(msgSwapExactAmountIn.Sender, msgSwapExactAmountIn.TokenIn, msgSwapExactAmountIn.TokenOutMinAmount, PoolIds, TokenOutDenoms)

	if err := msgMsgSwapExactAmountIn.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "failed validating MsgMsgSwapExactAmountIn")
	}

	// Create denom
	_, err := msgServer.SwapExactAmountIn(
		sdk.WrapSDKContext(ctx),
		msgMsgSwapExactAmountIn,
	)
	if err != nil {
		return errorsmod.Wrap(err, "swap msg")
	}
	return nil
}

type ElysMsg struct {
	MsgSwapExactAmountIn *MsgSwapExactAmountIn `json:"msg_swap_exact_amount_in,omitempty"`
}

type MsgSwapExactAmountIn struct {
	Sender            string                      `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Routes            []ammtype.SwapAmountInRoute `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	TokenIn           sdk.Coin                    `protobuf:"bytes,3,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	TokenOutMinAmount cosmos_sdk_math.Int         `protobuf:"bytes,4,opt,name=tokenOutMinAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_min_amount,omitempty"`
}

type MsgSwapExactAmountInResponse struct {
	TokenOutAmount string `protobuf:"bytes,1,opt,name=tokenOutAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_amount,omitempty"`
}
