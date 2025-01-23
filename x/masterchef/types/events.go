package types

const (
	TypeEvtAddExternalRewardDenom = "add_external_reward_denom"
	TypeEvtAddExternalIncentive   = "add_external_incentive"
	TypeEvtClaimRewards           = "claim_rewards"
	TypeEvtSetPoolMultiplier      = "set_pool_multiplier"
	TypeEvtSkipSwap               = "skip_conversion_of_fee"
	TypeEvtUsdcFee                = "fee_collected_in_usdc"

	AttributeRewardDenom    = "reward_denom"
	AttributeMinAmount      = "min_amount"
	AttributeSupported      = "supported"
	AttributePoolId         = "pool_id"
	AttributePoolIds        = "pool_ids"
	AttributeFromBlock      = "from_block"
	AttributeToBlock        = "to_block"
	AttributeAmountPerBlock = "amount_per_block"
	AttributeSender         = "sender"
	AttributeRecipient      = "recipient"
	AttributeMultiplier     = "multiplier"
)
