package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DateFormat = "2006-01-02"

func (p Portfolio) GetCreatorAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Creator)
}

func NewPortfolioWithContextDate(date string, creator sdk.AccAddress, value sdkmath.LegacyDec) Portfolio {
	return Portfolio{
		Creator:   creator.String(),
		Date:      date,
		Portfolio: value,
	}
}

var (
	Bronze = MembershipTier{
		Discount:         sdkmath.LegacyZeroDec(),
		Membership:       MembershipTierType_BRONZE,
		MinimumPortfolio: sdkmath.LegacyZeroDec(),
	}
	Silver = MembershipTier{
		Discount:         sdkmath.LegacyMustNewDecFromStr("0.1"),
		Membership:       MembershipTierType_SILVER,
		MinimumPortfolio: sdkmath.LegacyMustNewDecFromStr("50000"),
	}
	Gold = MembershipTier{
		Discount:         sdkmath.LegacyMustNewDecFromStr("0.2"),
		Membership:       MembershipTierType_GOLD,
		MinimumPortfolio: sdkmath.LegacyMustNewDecFromStr("250000"),
	}
	Platinum = MembershipTier{
		Discount:         sdkmath.LegacyMustNewDecFromStr("0.3"),
		Membership:       MembershipTierType_PLATINUM,
		MinimumPortfolio: sdkmath.LegacyMustNewDecFromStr("500000"),
	}
)
