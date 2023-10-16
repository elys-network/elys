package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}

	commitments, found := k.commKeeper.GetCommitments(ctx, mtp.GetMTPAddress().String())
	if !found {
		return sdk.ZeroDec(), nil
	}

	mtpVal := sdk.ZeroDec()
	for _, commitment := range commitments.CommittedTokens {
		ammPoolTvl, err := ammPool.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		mtpVal = mtpVal.Add(
			ammPoolTvl.
				Mul(math.LegacyNewDecFromInt(commitment.Amount)).
				Quo(math.LegacyNewDecFromInt(ammPool.TotalShares.Amount)))
	}

	lr := mtpVal.Quo(sdk.NewDecFromBigInt(xl.BigInt()))
	return lr, nil
}
