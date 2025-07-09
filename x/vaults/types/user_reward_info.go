package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (u UserRewardInfo) GetUserAccount() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(u.User)
}
