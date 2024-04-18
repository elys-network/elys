package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	// initialize pool infos from incentive module
	incentiveParams := m.incentiveKeeper.GetParams(ctx)
	for _, poolInfo := range incentiveParams.PoolInfos {
		m.masterchefKeeper.SetPool(ctx, mastercheftypes.PoolInfo{
			PoolId:                poolInfo.PoolId,
			RewardWallet:          poolInfo.RewardWallet,
			Multiplier:            poolInfo.Multiplier,
			NumBlocks:             poolInfo.NumBlocks,
			DexRewardAmountGiven:  poolInfo.DexRewardAmountGiven,
			EdenRewardAmountGiven: poolInfo.EdenRewardAmountGiven,
			EdenApr:               poolInfo.EdenApr,
			DexApr:                poolInfo.DexApr,
			ExternalIncentiveApr:  sdk.ZeroDec(),
			ExternalRewardDenoms:  []string{},
		})
	}

	// initiate masterchef params
	m.masterchefKeeper.SetParams(ctx, mastercheftypes.NewParams(
		nil, // TODO:
		sdk.NewDecWithPrec(60, 2),
		sdk.NewDecWithPrec(25, 2),
		mastercheftypes.DexRewardsTracker{
			NumBlocks: sdk.NewInt(1),
			Amount:    sdk.ZeroDec(),
		},
		sdk.NewDecWithPrec(5, 1),
		"elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
	))

	// initiate estaking module data
	m.estakingKeeper.InitGenesis(ctx, estakingtypes.GenesisState{
		Params: estakingtypes.Params{
			StakeIncentives:         nil, // TODO:
			EdenCommitVal:           "",
			EdenbCommitVal:          "",
			MaxEdenRewardAprStakers: sdk.NewDecWithPrec(3, 1), // 30%
			EdenBoostApr:            sdk.OneDec(),
			DexRewardsStakers: estakingtypes.DexRewardsTracker{
				NumBlocks: sdk.OneInt(),
				Amount:    sdk.ZeroDec(),
			},
		},
	})

	// initiate missing distribution module data
	m.distrKeeper.InitGenesis(ctx, *distrtypes.DefaultGenesisState())

	// execute missing validator creation hooks
	validators := m.estakingKeeper.Keeper.GetAllValidators(ctx)
	for _, val := range validators {
		err := m.estakingKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		if err != nil {
			panic(err)
		}
	}

	// execute missing delegation creation hooks
	allDelegations := m.estakingKeeper.Keeper.GetAllDelegations(ctx)
	for _, delegation := range allDelegations {
		delAddr := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		err = m.estakingKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr)
		if err != nil {
			panic(err)
		}
	}

	// Update all commitments (move all unclaimed into claimed)
	// and execute missing eden/edenb commitment hooks
	edenValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.Eden))
	edenBValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.EdenB))
	legacyCommitments := m.commitmentKeeper.GetAllLegacyCommitments(ctx)
	for _, legacy := range legacyCommitments {
		creator := legacy.Creator
		addr, err := sdk.AccAddressFromBech32(creator)
		if err != nil {
			// This is validator address
			m.commitmentKeeper.RemoveCommitments(ctx, creator)
			continue
		}

		commitments := commitmenttypes.Commitments{
			Creator:         legacy.Creator,
			CommittedTokens: legacy.CommittedTokens,
			Claimed:         legacy.Claimed.Add(legacy.RewardsUnclaimed...),
			VestingTokens:   legacy.VestingTokens,
		}
		m.commitmentKeeper.SetCommitments(ctx, commitments)
		for _, committed := range commitments.CommittedTokens {
			if committed.Denom == ptypes.Eden {
				err = m.estakingKeeper.Hooks().AfterDelegationModified(ctx, addr, edenValAddr)
				if err != nil {
					panic(err)
				}
			}
			if committed.Denom == ptypes.EdenB {
				err = m.estakingKeeper.Hooks().AfterDelegationModified(ctx, addr, edenBValAddr)
				if err != nil {
					panic(err)
				}
			}
		}

	}

	return nil
}
