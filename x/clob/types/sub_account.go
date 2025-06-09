package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"math"
)

const CrossMarginSubAccountId = uint64(math.MaxInt16) // 32767

func (s SubAccount) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(s.Owner)
}

func (s SubAccount) GetTradingAccountAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/%s/%d", s.Owner, s.Id))
}

func (s SubAccount) IsCrossMargin() bool {
	return s.Id == CrossMarginSubAccountId
}

func (s SubAccount) IsIsolated() bool {
	return !s.IsCrossMargin()
}
