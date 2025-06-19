package types

// epochs events
const (
	EventTypeCommitmentChanged  = "commitment_changed"
	EventTypeClaimVesting       = "claim_vesting"
	EventTypeMintCoins          = "mint_coins"
	EventTypeBurnCoins          = "burn_coins"
	EventTypeSendCoins          = "send_coins"
	EventTypeClaimRewardProgram = "claim_reward_program"

	AttributeCreator             = "creator"
	AttributeAmount              = "token_amount"
	AttributeCancelledAmount     = "cancelled_amount"
	AttributeDenom               = "token_denom"
	AttributeKeyClaimAddress     = "claim_address"
	AttributeKeyEdenAmount       = "eden_amount"
	AttributeKeyTotalEdenClaimed = "total_eden_claimed"
)
