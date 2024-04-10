package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	m.keeper.SetParams(ctx, types.Params{
		LpIncentives:            nil,
		StakeIncentives:         nil,
		RewardPortionForLps:     sdk.NewDecWithPrec(60, 2),
		RewardPortionForStakers: sdk.NewDecWithPrec(30, 2),
		PoolInfos: []types.PoolInfo{
			{
				DexApr:                sdk.ZeroDec(),
				EdenApr:               sdk.ZeroDec(),
				EdenRewardAmountGiven: sdk.ZeroInt(),
				DexRewardAmountGiven:  sdk.ZeroDec(),
				Multiplier:            sdk.NewDec(1),
				NumBlocks:             sdk.NewInt(1),
				PoolId:                32767,
				RewardWallet:          "elys12dxadvd5def5gfy6jwmmt3gqs2pt5k9xupx8ja",
			},
		},
		ElysStakeSnapInterval: 10,
		DexRewardsStakers: types.DexRewardsTracker{
			NumBlocks:                     sdk.NewInt(1011431),
			Amount:                        sdk.NewDec(708624172),
			AmountCollectedByOtherTracker: sdk.ZeroDec(),
		},
		DexRewardsLps: types.DexRewardsTracker{
			NumBlocks:                     sdk.NewInt(1011431),
			Amount:                        sdk.NewDec(1446171780),
			AmountCollectedByOtherTracker: sdk.ZeroDec(),
		},
		MaxEdenRewardAprStakers: sdk.NewDecWithPrec(3, 1), // 30%
		MaxEdenRewardAprLps:     sdk.NewDecWithPrec(5, 1), // 50%
	})
	return nil
}
