package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type BalanceBorrowed struct {
	UsdAmount  sdk.Dec `protobuf:"bytes,1,rep,name=usd_amount,proto3" json:"usd_amount"`
	Percentage sdk.Dec `protobuf:"bytes,2,rep,name=percentage,proto3" json:"percentage"`
}
