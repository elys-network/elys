package types

const (
	TypeEvtAddExternalRewardDenom = "add_external_reward_denom"
	TypeEvtAddExternalIncentive   = "add_external_incentive"
	TypeEvtClaimRewards           = "claim_rewards"
	TypeEvtSetPoolMultiplier      = "set_pool_multiplier"

	AttributeRewardDenom    = "reward_denom"
	AttributeMinAmount      = "min_amount"
	AttributeSupported      = "supported"
	AttributePoolId         = "pool_id"
	AttributeFromBlock      = "from_block"
	AttributeToBlock        = "to_block"
	AttributeAmountPerBlock = "amount_per_block"
	AttributeSender         = "sender"
	AttributeMultiplier     = "multiplier"
)
