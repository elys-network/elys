package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DateFormat = "2006-01-02"

func (p Portfolio) GetCreatorAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Creator)
}

func NewPortfolioWithContextDate(date string, creator sdk.AccAddress, value sdk.Dec) Portfolio {
	return Portfolio{
		Creator:   creator.String(),
		Date:      date,
		Portfolio: value,
	}
}
