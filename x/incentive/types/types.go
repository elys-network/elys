package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TotalCommitmentInfo
// Stores the
type TotalCommitmentInfo struct {
	// Total Elys staked
	TotalElysBonded sdk.Int
	// Total Eden + Eden boost committed
	TotalEdenEdenBoostCommitted sdk.Int
	// Gas fees collected and DEX revenus
	TotalFeesCollected sdk.Coins
	// Total Lp Token committed
	TotalLpTokensCommitted map[string]sdk.Int
	// Revenue tracking per pool, key => (poolId)
	PoolRevenueTrack map[string]sdk.Dec
}

// Returns the pool revenue tracking key.
// Unique per pool per epoch, clean once complete the calculation.
func GetPoolRevenueTrackKey(poolId uint64) string {
	return fmt.Sprintf("pool_revenue_%d", poolId)
}

const (
	MAX_RETRY_VALIDATORS = uint32(200)
)

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
