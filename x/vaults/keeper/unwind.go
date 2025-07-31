package keeper

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPoolIdFromShareDenom(shareDenom string) (uint64, error) {
	poolId, err := strconv.Atoi(strings.TrimPrefix(shareDenom, "amm/pool/"))
	if err != nil {
		return 0, err
	}
	return uint64(poolId), nil
}

func (k Keeper) GetAllPoolIds(ctx sdk.Context, vaultAddress sdk.AccAddress) []uint64 {
	commitments := k.commitment.GetCommitments(ctx, vaultAddress)
	poolIds := make([]uint64, 0)
	for _, commitment := range commitments.CommittedTokens {
		if strings.HasPrefix(commitment.Denom, "amm/pool/") {
			poolId, err := GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				return nil
			}
			poolIds = append(poolIds, poolId)
		}
	}
	return poolIds
}
