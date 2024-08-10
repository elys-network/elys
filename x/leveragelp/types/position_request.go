package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (pr PositionRequest) GetAccountAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(pr.Address)
}
