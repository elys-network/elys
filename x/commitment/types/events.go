package types

// epochs events
const (
	EventTypeCommitmentChanged = "commitment_changed"
	EventTypeClaimVesting      = "claim_vesting"
	EventTypeMintCoins         = "mint_coins"
	EventTypeBurnCoins         = "burn_coins"
	EventTypeSendCoins         = "send_coins"

	AttributeCreator = "creator"
	AttributeAmount  = "token_amount"
	AttributeDenom   = "token_denom"
)
