package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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
	Basic = MembershipTier{
		Discount:         sdkmath.LegacyZeroDec(),
		Membership:       MembershipTierType_BASIC,
		MinimumPortfolio: sdkmath.LegacyZeroDec(),
	}
	Bronze = MembershipTier{
		Discount:         sdkmath.LegacyMustNewDecFromStr("0.05"),
		Membership:       MembershipTierType_BRONZE,
		MinimumPortfolio: sdkmath.LegacyMustNewDecFromStr("25000"),
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

func (m MembershipTier) GetBigDecDiscount() osmomath.BigDec {
	return osmomath.BigDecFromDec(m.Discount)
}

func (m MembershipTier) GetBigDecMinimumPortfolio() osmomath.BigDec {
	return osmomath.BigDecFromDec(m.MinimumPortfolio)
}

func (p *Portfolio) GetBigDecPortFolio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.Portfolio)
}
