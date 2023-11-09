package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

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
