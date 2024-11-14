package auth

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
)

type QueryAddressesParams struct {
	Context sdk.Context
	Request *authtypes.QueryAccountsRequest
}

func (oq *Querier) queryAddresses(params QueryAddressesParams) ([]byte, error) {
	if params.Request == nil {
		return nil, fmt.Errorf("query request cannot be nil")
	}

	queryServer := authkeeper.NewQueryServer(*oq.keeper)
	res, err := queryServer.Accounts(sdk.WrapSDKContext(params.Context), params.Request)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to get accounts for request: %v", params.Request)
	}

	// Pre-allocate slice with capacity matching the number of accounts
	addresses := make([]string, 0, len(res.Accounts))

	for _, account := range res.Accounts {
		iaccount, ok := account.GetCachedValue().(authtypes.AccountI)
		if !ok || iaccount == nil {
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
		return nil, errorsmod.Wrapf(err, "failed to serialize addresses response: %+v", authAddressesResponse)
	}

	return responseBytes, nil
}
