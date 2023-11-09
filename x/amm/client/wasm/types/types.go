package types

import (
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

type QuerySwapEstimationRequest struct {
	TokenIn sdk.Coin                     `protobuf:"bytes,2,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	Routes  []ammtypes.SwapAmountInRoute `protobuf:"bytes,1,rep,name=routes,proto3" json:"routes,omitempty"`
}

type QuerySwapEstimationResponse struct {
	SpotPrice sdk.Dec  `protobuf:"bytes,1,opt,name=SpotPrice,proto3" json:"spot_price,omitempty"`
	TokenOut  sdk.Coin `protobuf:"bytes,2,opt,name=tokenOut,proto3" json:"token_out,omitempty"`
}

type QueryBalanceRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Denom   string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}

type QueryBalanceResponse struct {
	Balance sdk.Coin `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance,omitempty"`
}

type MsgSwapExactAmountIn struct {
	Sender            string                       `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Routes            []ammtypes.SwapAmountInRoute `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	TokenIn           sdk.Coin                     `protobuf:"bytes,3,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	TokenOutMinAmount cosmos_sdk_math.Int          `protobuf:"bytes,4,opt,name=tokenOutMinAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_min_amount,omitempty"`
	MetaData          *[]byte                      `protobuf:"bytes,5,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgSwapExactAmountInResponse struct {
	TokenOutAmount cosmos_sdk_math.Int `protobuf:"bytes,1,opt,name=tokenOutAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_amount,omitempty"`
	MetaData       *[]byte             `protobuf:"bytes,2,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}
