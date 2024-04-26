package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (oq *Querier) queryStakedBalanceOfDenom(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	edenDenomPrice := sdk.ZeroDec()
	entry, found := oq.assetKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if found {
		edenDenomPrice = oq.ammKeeper.GetEdenDenomPrice(ctx, entry.Denom)
	}
	baseCurrency := entry.Denom

	denom := query.Denom
	addr := query.Address
	address, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	bondedAmt := oq.stakingKeeper.GetDelegatorBonded(ctx, address)
	balance := sdk.NewCoin(denom, bondedAmt)
	usdAmount := edenDenomPrice.MulInt(balance.Amount)
	lockups := make([]commitmenttypes.Lockup, 0)

	if denom != ptypes.Elys {
		commitment := oq.keeper.GetCommitments(ctx, addr)
		if denom == ptypes.BaseCurrency {
			denom = stabletypes.GetShareDenom()
		}

		committedToken := commitment.GetCommittedAmountForDenom(denom)
		lockups = commitment.GetCommittedLockUpsForDenom(denom)
		balance = sdk.NewCoin(denom, committedToken)
		if denom == ptypes.Eden || denom == ptypes.EdenB {
			usdAmount = edenDenomPrice.MulInt(balance.Amount)
		}

		if denom == stabletypes.GetShareDenom() {
			stableShareDenomPrice := oq.stableKeeper.ShareDenomPrice(ctx, oq.oracleKeeper, baseCurrency)
			usdAmount = stableShareDenomPrice.MulInt(balance.Amount)
		}
	}

	resp := commitmenttypes.StakedAvailable{
		Amount:    balance.Amount,
		UsdAmount: usdAmount,
		Lockups:   lockups,
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get staked balance response")
	}

	return responseBytes, nil
}
