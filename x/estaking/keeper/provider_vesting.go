package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	commitmentmoduletypes "github.com/elys-network/elys/v4/x/commitment/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

func (k Keeper) ClaimAndVestProviderStakingRewards(ctx sdk.Context) error {
	providerRewardAddress := authtypes.NewModuleAddress(ccvconsumertypes.ConsumerToSendToProviderName)

	// Claiming first to remove one slot to vest next
	claimVestingMsg := commitmentmoduletypes.MsgClaimVesting{Sender: providerRewardAddress.String()}
	_, err := k.commKeeper.ClaimVesting(ctx, &claimVestingMsg)
	if err != nil {
		return err
	}

	commitments := k.commKeeper.GetCommitments(ctx, providerRewardAddress)
	// We don't edenB to provider chain
	edenAmount := commitments.Claimed.AmountOf(ptypes.Eden)
	if !edenAmount.IsZero() {
		err = k.commKeeper.ProcessTokenVesting(ctx, ptypes.Eden, edenAmount, providerRewardAddress)
		if err != nil {
			if err == commitmentmoduletypes.ErrExceedMaxVestings {
				return nil
			}
			return err
		}

	}
	return nil
}
