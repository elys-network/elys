package types

import (
	cosmos_sdk_math "cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type QueryBorrowAmountRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

type QueryDelegatorDelegationsRequest struct {
	// delegator_addr defines the delegator address to query for.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegator_address,proto3" json:"delegator_address,omitempty"`
}

// QueryDelegatorDelegationsResponse is response type for the
// Query/DelegatorDelegations RPC method.
type QueryDelegatorDelegationsResponse struct {
	// delegation_responses defines all the delegations' info of a delegator.
	DelegationResponses []stakingtypes.DelegationResponse `protobuf:"bytes,1,rep,name=delegation_responses,json=delegationResponses,proto3" json:"delegation_responses"`
}

// UnbondingDelegationEntry defines an unbonding object with relevant metadata.
type UnbondingDelegationEntry struct {
	// creation_height is the height which the unbonding took place.
	CreationHeight int64 `protobuf:"varint,1,opt,name=creation_height,json=creationHeight,proto3" json:"creation_height,omitempty"`
	// completion_time is the unix time for unbonding completion.
	CompletionTime int64 `protobuf:"bytes,2,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
	// initial_balance defines the tokens initially scheduled to receive at completion.
	InitialBalance cosmos_sdk_math.Int `protobuf:"bytes,3,opt,name=initial_balance,json=initialBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"initial_balance"`
	// balance defines the tokens to receive at completion.
	Balance cosmos_sdk_math.Int `protobuf:"bytes,4,opt,name=balance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"balance"`
	// Incrementing id that uniquely identifies this entry
	UnbondingId uint64 `protobuf:"varint,5,opt,name=unbonding_id,json=unbondingId,proto3" json:"unbonding_id,omitempty"`
}

// QueryDelegatorUnbondingDelegationsRequest is request type for the
type QueryDelegatorUnbondingDelegationsRequest struct {
	// delegator_addr defines the delegator address to query for.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegator_address,proto3" json:"delegator_address,omitempty"`
}

type UnbondingDelegation struct {
	// delegator_address is the bech32-encoded address of the delegator.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	// validator_address is the bech32-encoded address of the validator.
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	// entries are the unbonding delegation entries.
	Entries []UnbondingDelegationEntry `protobuf:"bytes,3,rep,name=entries,proto3" json:"entries"`
}

// QueryUnbondingDelegatorDelegationsResponse is response type for the
// Query/UnbondingDelegatorDelegations RPC method.
type QueryDelegatorUnbondingDelegationsResponse struct {
	UnbondingResponses []UnbondingDelegation `protobuf:"bytes,1,rep,name=unbonding_responses,json=unbondingResponses,proto3" json:"unbonding_responses"`
}
