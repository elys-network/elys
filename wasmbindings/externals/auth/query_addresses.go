package auth

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (oq *Querier) queryAddresses(ctx sdk.Context, req *authtypes.QueryAccountsRequest) ([]byte, error) {
	res, err := oq.keeper.Accounts(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get accounts")
	}

	var addresses []string
	for _, account := range res.Accounts {
		iaccount := account.GetCachedValue().(authtypes.AccountI)
		if iaccount == nil {
			continue
		}
		addresses = append(addresses, iaccount.GetAddress().String())
	}

	responseBytes, err := json.Marshal(addresses)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize addresses response")
	}
	return responseBytes, nil
}
