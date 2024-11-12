package wasm

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (oq *Querier) queryStakedBalanceOfDenom(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	edenDenomPrice := sdkmath.LegacyZeroDec()
	baseCurrency, found := oq.assetKeeper.GetUsdcDenom(ctx)
	if found {
		edenDenomPrice = oq.ammKeeper.GetEdenDenomPrice(ctx, baseCurrency)
	}

	address, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	bondedAmt, err := oq.stakingKeeper.GetDelegatorBonded(ctx, address)
	if err != nil {
		return nil, err
	}
	balance := sdk.NewCoin(query.Denom, bondedAmt)
	usdAmount := edenDenomPrice.MulInt(balance.Amount)
	lockups := make([]commitmenttypes.Lockup, 0)

	if query.Denom != ptypes.Elys {
		commitment := oq.keeper.GetCommitments(ctx, address)
		if query.Denom == ptypes.BaseCurrency {
			query.Denom = stabletypes.GetShareDenom()
		}

		committedToken := commitment.GetCommittedAmountForDenom(query.Denom)
		lockups = commitment.GetCommittedLockUpsForDenom(query.Denom)
		balance = sdk.NewCoin(query.Denom, committedToken)
		if query.Denom == ptypes.Eden || query.Denom == ptypes.EdenB {
			usdAmount = edenDenomPrice.MulInt(balance.Amount)
		}

		if query.Denom == stabletypes.GetShareDenom() {
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
