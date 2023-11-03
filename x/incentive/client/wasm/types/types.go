package types

import (
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBeginRedelegate struct {
	DelegatorAddress    string   `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorSrcAddress string   `protobuf:"bytes,2,opt,name=validator_src_address,json=validatorSrcAddress,proto3" json:"validator_src_address,omitempty"`
	ValidatorDstAddress string   `protobuf:"bytes,3,opt,name=validator_dst_address,json=validatorDstAddress,proto3" json:"validator_dst_address,omitempty"`
	Amount              sdk.Coin `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
}

type MsgCancelUnbondingDelegation struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	// amount is always less than or equal to unbonding delegation entry balance
	Amount sdk.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
	// creation_height is the height which the unbonding took place.
	CreationHeight int64 `protobuf:"varint,4,opt,name=creation_height,json=creationHeight,proto3" json:"creation_height,omitempty"`
}

type MsgVest struct {
	Creator string              `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Amount  cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount"`
	Denom   string              `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

type MsgCancelVest struct {
	Creator string              `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Amount  cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount"`
	Denom   string              `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

type MsgWithdrawTokens struct {
	Creator string              `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Amount  cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount"`
	Denom   string              `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

type MsgWithdrawRewards struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	Denom            string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}

type MsgWithdrawValidatorCommission struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Denom            string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}
