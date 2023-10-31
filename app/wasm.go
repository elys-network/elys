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
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
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
	ammKeeper    *ammkeeper.Keeper
	oracleKeeper *oraclekeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		ammKeeper:    amm,
		oracleKeeper: oracle,
	}
}

func RegisterCustomPlugins(
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	margin *marginkeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(amm, oracle)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(amm, margin),
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
		case contractQuery.QuerySwapEstimation != nil:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "QuerySwapEstimation, not implemented yet"}

		case contractQuery.AssetInfo != nil:
			denom := contractQuery.AssetInfo.Denom

			AssetInfoResp, err := qp.oracleKeeper.AssetInfo(ctx, &oracletypes.QueryGetAssetInfoRequest{Denom: denom})
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to query asset info")
			}

			res := AssetInfoResponse{
				AssetInfo: &AssetInfoType{
					Denom:      AssetInfoResp.AssetInfo.Denom,
					Display:    AssetInfoResp.AssetInfo.Display,
					BandTicker: AssetInfoResp.AssetInfo.BandTicker,
					ElysTicker: AssetInfoResp.AssetInfo.ElysTicker,
					Decimal:    AssetInfoResp.AssetInfo.Decimal,
				},
			}

			responseBytes, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to serialize asset info response")
			}
			return responseBytes, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown elys query variant"}
		}
	}
}

type ElysQuery struct {
	PriceAll            *PriceAll                   `json:"price_all,omitempty"`
	QuerySwapEstimation *QuerySwapEstimationRequest `json:"query_swap_estimation,omitempty"`
	AssetInfo           *AssetInfo                  `json:"asset_info,omitempty"`
}

type PriceAll struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type AllPriceResponse struct {
	Price      []oracletypes.Price `protobuf:"bytes,1,rep,name=price,proto3" json:"price"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type QuerySwapEstimationRequest struct {
	TokenIn sdk.Coin                    `protobuf:"bytes,2,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	Routes  []ammtype.SwapAmountInRoute `protobuf:"bytes,1,rep,name=routes,proto3" json:"routes,omitempty"`
}

type QuerySwapEstimationResponse struct {
	SpotPrice sdk.Dec  `protobuf:"bytes,1,opt,name=SpotPrice,proto3" json:"spot_price,omitempty"`
	TokenOut  sdk.Coin `protobuf:"bytes,2,opt,name=tokenOut,proto3" json:"token_out,omitempty"`
}

type AssetInfo struct {
	Denom string `protobuf:"bytes,1,opt,name=Denom,proto3" json:"Denom,omitempty"`
}

type AssetInfoResponse struct {
	AssetInfo *AssetInfoType `protobuf:"bytes,1,opt,name=AssetInfo,proto3" json:"asset_info,omitempty"`
}

type AssetInfoType struct {
	Denom      string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Display    string `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	BandTicker string `protobuf:"bytes,3,opt,name=bandTicker,proto3" json:"band_ticker,omitempty"`
	ElysTicker string `protobuf:"bytes,4,opt,name=elysTicker,proto3" json:"elys_ticker,omitempty"`
	Decimal    uint64 `protobuf:"varint,5,opt,name=decimal,proto3" json:"decimal,omitempty"`
}

func CustomMessageDecorator(amm *ammkeeper.Keeper, margin *marginkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			amm:     amm,
			margin:  margin,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	amm     *ammkeeper.Keeper
	margin  *marginkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg ElysMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "elys msg")
		}
		if contractMsg.MsgSwapExactAmountIn != nil {
			return m.msgSwapExactAmountIn(ctx, contractAddr, contractMsg.MsgSwapExactAmountIn)
		}
		if contractMsg.MsgClose != nil {
			return m.msgClose(ctx, contractAddr, contractMsg.MsgClose)
		}
		if contractMsg.MsgOpen != nil {
			return m.msgOpen(ctx, contractAddr, contractMsg.MsgOpen)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) msgSwapExactAmountIn(ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *MsgSwapExactAmountIn) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgSwapExactAmountIn(m.amm, ctx, contractAddr, msgSwapExactAmountIn)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform swap")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize swap response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgSwapExactAmountIn(f *ammkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *MsgSwapExactAmountIn) (*MsgSwapExactAmountInResponse, error) {
	if msgSwapExactAmountIn == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "swap null swap"}
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
		return nil, errorsmod.Wrap(err, "failed validating MsgMsgSwapExactAmountIn")
	}

	// Swap
	swapResp, err := msgServer.SwapExactAmountIn(
		sdk.WrapSDKContext(ctx),
		msgMsgSwapExactAmountIn,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "swap msg")
	}

	var resp = &MsgSwapExactAmountInResponse{
		TokenOutAmount: swapResp.TokenOutAmount,
		MetaData:       msgSwapExactAmountIn.MetaData,
	}
	return resp, nil
}

func (m *CustomMessenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *MsgOpen) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgOpen(m.margin, ctx, contractAddr, msgOpen)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform open")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgOpen(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *MsgOpen) (*MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin open null margin open"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := margintypes.NewMsgOpen(msgOpen.Creator, msgOpen.CollateralAsset, cosmos_sdk_math.Int(msgOpen.CollateralAmount), msgOpen.BorrowAsset, msgOpen.Position, msgOpen.Leverage, msgOpen.TakeProfitPrice)

	if err := msgMsgOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	_, err := msgServer.Open(ctx, msgMsgOpen) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin open msg")
	}

	var resp = &MsgOpenResponse{
		MetaData: msgOpen.MetaData,
	}
	return resp, nil
}

func (m *CustomMessenger) msgClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *MsgClose) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgClose(m.margin, ctx, contractAddr, msgClose)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform close")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize close response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgClose(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *MsgClose) (*MsgCloseResponse, error) {
	if msgClose == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin close null margin close"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgClose := margintypes.NewMsgClose(msgClose.Creator, uint64(msgClose.Id))

	if err := msgMsgClose.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	_, err := msgServer.Close(ctx, msgMsgClose) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin close msg")
	}

	var resp = &MsgCloseResponse{
		MetaData: msgClose.MetaData,
	}
	return resp, nil
}

type ElysMsg struct {
	MsgSwapExactAmountIn *MsgSwapExactAmountIn `json:"msg_swap_exact_amount_in,omitempty"`
	MsgOpen              *MsgOpen              `json:"msg_open,omitempty"`
	MsgClose             *MsgClose             `json:"msg_close,omitempty"`
}

type MsgSwapExactAmountIn struct {
	Sender            string                      `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Routes            []ammtype.SwapAmountInRoute `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	TokenIn           sdk.Coin                    `protobuf:"bytes,3,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	TokenOutMinAmount cosmos_sdk_math.Int         `protobuf:"bytes,4,opt,name=tokenOutMinAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_min_amount,omitempty"`
	MetaData          *[]byte                     `protobuf:"bytes,5,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgSwapExactAmountInResponse struct {
	TokenOutAmount cosmos_sdk_math.Int `protobuf:"bytes,1,opt,name=tokenOutAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_amount,omitempty"`
	MetaData       *[]byte             `protobuf:"bytes,2,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgOpen struct {
	Creator          string               `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	CollateralAsset  string               `protobuf:"bytes,2,opt,name=collateralAsset,proto3" json:"collateral_asset,omitempty"`
	CollateralAmount sdk.Uint             `protobuf:"bytes,3,opt,name=collateralAmount,proto3" json:"collateral_amount,omitempty"`
	BorrowAsset      string               `protobuf:"bytes,4,opt,name=borrowAsset,proto3" json:"borrow_asset,omitempty"`
	Position         margintypes.Position `protobuf:"bytes,5,opt,name=position,proto3" json:"position,omitempty"`
	Leverage         sdk.Dec              `protobuf:"bytes,6,opt,name=leverage,proto3" json:"leverage,omitempty"`
	TakeProfitPrice  sdk.Dec              `protobuf:"bytes,7,opt,name=takeProfitPrice,proto3" json:"take_profit_price,omitempty"`
	MetaData         *[]byte              `protobuf:"bytes,8,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgClose struct {
	Creator  string  `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Id       int64   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	MetaData *[]byte `protobuf:"bytes,3,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgOpenResponse struct {
	MetaData *[]byte `protobuf:"bytes,1,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}
type MsgCloseResponse struct {
	MetaData *[]byte `protobuf:"bytes,1,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}
