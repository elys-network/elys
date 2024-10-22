package types

import (
	"cosmossdk.io/math"
	"fmt"
)

func (ap AccountedPool) GetNonAmmTokenBalance(denom string) (math.Int, error) {
	for _, nonAmmToken := range ap.NonAmmPoolTokens {
		if nonAmmToken.Denom == denom {
			return nonAmmToken.Amount, nil
		}
	}
	return math.ZeroInt(), fmt.Errorf("denom %s not exist in accounted pool", denom)
}
