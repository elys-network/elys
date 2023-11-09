package types

import (
	cosmos_sdk_math "cosmossdk.io/math"
)

type MsgStake struct {
	Address          string              `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount           cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Asset            string              `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	ValidatorAddress string              `protobuf:"bytes,4,opt,name=validator_address,proto3" json:"validator_address,omitempty"`
}

type MsgUnstake struct {
	Address          string              `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount           cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Asset            string              `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	ValidatorAddress string              `protobuf:"bytes,4,opt,name=validator_address,proto3" json:"validator_address,omitempty"`
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}
