package auth

import (
	"encoding/json"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
)

func (oq *Querier) queryAddresses(ctx sdk.Context, req *authtypes.QueryAccountsRequest) ([]byte, error) {
	queryServer := authkeeper.NewQueryServer(*oq.keeper)
	res, err := queryServer.Accounts(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get accounts")
	}

	var addresses []string
	for _, account := range res.Accounts {
		iaccount := account.GetCachedValue().(sdk.AccountI)
		if iaccount == nil {
			continue
		}
		addresses = append(addresses, iaccount.GetAddress().String())
	}

	authAddressesResponse := wasmbindingstypes.AuthAddressesResponse{
		Addresses:  addresses,
		Pagination: res.Pagination,
	}

	responseBytes, err := json.Marshal(authAddressesResponse)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize addresses response")
	}
	return responseBytes, nil
}
