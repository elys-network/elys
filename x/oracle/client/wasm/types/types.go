package types

import (
	query "github.com/cosmos/cosmos-sdk/types/query"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

type PriceAll struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type AllPriceResponse struct {
	Price      []oracletypes.Price `protobuf:"bytes,1,rep,name=price,proto3" json:"price"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
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
