package auth

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (oq *Querier) queryAccounts(ctx sdk.Context, req *types.QueryAccountsRequest) ([]byte, error) {
	res, err := oq.keeper.Accounts(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get accounts")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize accounts response")
	}
	return responseBytes, nil
}
