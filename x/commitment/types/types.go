package types

import (
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// QueryValidatorsRequest is request type for Query/Validators RPC method.
type QueryValidatorsRequest struct {
	// status enables to query for validators matching a given status.
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,proto3" json:"delegator_address,omitempty"`
}

// QueryDelegatorValidatorsResponse is response type for the
// Query/DelegatorValidators RPC method.
type QueryDelegatorValidatorsResponse struct {
	// validators defines the validators' info of a delegator.
	Validators []ValidatorDetail `protobuf:"bytes,1,rep,name=validators,proto3" json:"validators"`
}

type BalanceAvailable struct {
	Amount    uint64  `protobuf:"bytes,1,rep,name=amount,proto3" json:"amount"`
	UsdAmount sdk.Dec `protobuf:"bytes,2,rep,name=usd_amount,proto3" json:"usd_amount"`
}

type ValidatorDetail struct {
	// The validator address.
	Address string `protobuf:"bytes,2,rep,name=address,proto3" json:"address"`
	// The validator name.
	Name string `protobuf:"bytes,3,rep,name=name,proto3" json:"name"`
	// Voting power percentage for this validator.
	VotingPower sdk.Dec `protobuf:"bytes,4,rep,name=voting_power,proto3" json:"voting_power"`
	// Comission percentage for the validator.
	Commission sdk.Dec `protobuf:"bytes,5,rep,name=commission,proto3" json:"commission"`
	// The url of the validator profile picture
	ProfilePictureSrc string `protobuf:"bytes,6,rep,name=profile_picture_src,proto3" json:"profile_picture_src"`
	// The staked amount the user has w/ this validator
	// Only available if there's some and if address.
	// is sent in request object.
	Staked BalanceAvailable `protobuf:"bytes,7,rep,name=staked,proto3" json:"staked"`
}

type StakingValidator struct {
	// The validator address.
	Address string `protobuf:"bytes,1,rep,name=address,proto3" json:"address"`
	// The validator name.
	Name string `protobuf:"bytes,2,rep,name=name,proto3" json:"name"`
	// Voting power percentage for this validator.
	VotingPower sdk.Dec `protobuf:"bytes,3,rep,name=voting_power,proto3" json:"voting_power"`
	// Comission percentage for the validator.
	Commission sdk.Dec `protobuf:"bytes,4,rep,name=commission,proto3" json:"commission"`
	// The url of the validator profile picture
	ProfilePictureSrc string `protobuf:"bytes,5,rep,name=profile_picture_src,proto3" json:"profile_picture_src"`
}

type StakedPosition struct {
	// The position ID.
	Id string `protobuf:"bytes,1,rep,name=id,proto3" json:"id"`
	// The validator that's being unstaked from.
	Validator StakingValidator `protobuf:"bytes,2,rep,name=validator,proto3" json:"validator"`
	// The amount that's being staked.
	Staked BalanceAvailable `protobuf:"bytes,3,rep,name=staked,proto3" json:"staked"`
}

type QueryStakedPositionResponse struct {
	StakedPosition []StakedPosition `protobuf:"bytes,1,rep,name=staked_position,proto3" json:"staked_position"`
}

type UnstakedPosition struct {
	// The position ID.
	Id string `protobuf:"bytes,1,rep,name=id,proto3" json:"id"`
	// The validator that's being unstaked from.
	Validator StakingValidator `protobuf:"bytes,2,rep,name=validator,proto3" json:"validator"`
	// Remaining time to unstake in days.
	RemainingTime uint64 `protobuf:"bytes,3,rep,name=remaining_time,proto3" json:"remaining_time"`
	// The amount that's being staked.
	Unstaked BalanceAvailable `protobuf:"bytes,4,rep,name=unstaked,proto3" json:"staked"`
}

type QueryUnstakedPositionResponse struct {
	UnstakedPosition []UnstakedPosition `protobuf:"bytes,1,rep,name=unstaked_position,proto3" json:"unstaked_position"`
}
