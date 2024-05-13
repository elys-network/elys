package migrations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	fmt.Println("Running incentive v11 migration ...")
	fmt.Println("1) Running PoolInfos migration ...")
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

	fmt.Println("2) Setting masterchef params ...")
	// initiate masterchef params
	m.masterchefKeeper.SetParams(ctx, mastercheftypes.NewParams(
		nil,
		sdk.NewDecWithPrec(60, 2),
		sdk.NewDecWithPrec(25, 2),
		mastercheftypes.DexRewardsTracker{
			NumBlocks: sdk.NewInt(1),
			Amount:    sdk.ZeroDec(),
		},
		sdk.NewDecWithPrec(5, 1),
		"elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
	))

	fmt.Println("3) Running init genesis for estaking ...")
	// initiate estaking module data
	m.estakingKeeper.InitGenesis(ctx, estakingtypes.GenesisState{
		Params: estakingtypes.Params{
			StakeIncentives:         nil,
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

	fmt.Println("4) Moving staking snapshots to estaking ...")
	// initiate delegation snapshot
	stakedSnapshots := m.incentiveKeeper.GetAllElysStaked(ctx)
	for _, snap := range stakedSnapshots {
		m.estakingKeeper.SetElysStaked(ctx, estakingtypes.ElysStaked{
			Address: snap.Address,
			Amount:  snap.Amount,
		})
	}

	fmt.Println("5) Running InitGenesis for distribution ...")
	// initiate missing distribution module data
	m.distrKeeper.InitGenesis(ctx, *distrtypes.DefaultGenesisState())

	fmt.Println("6) Running validator creation hooks ...")
	// execute missing validator creation hooks
	validators := m.estakingKeeper.Keeper.GetAllValidators(ctx)
	for _, val := range validators {
		err := m.estakingKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("7) Running delegation hooks ...")
	// execute missing delegation creation hooks
	allDelegations := m.estakingKeeper.Keeper.GetAllDelegations(ctx)
	for _, delegation := range allDelegations {
		delAddr := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		err = m.estakingKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr)
		if err != nil {
			panic(err)
		}
		err = m.estakingKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("8) Running commitment migrations ...")
	// Update all commitments (move all unclaimed into claimed)
	// and execute missing eden/edenb commitment hooks
	edenValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.Eden))
	edenBValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.EdenB))
	legacyCommitments := m.commitmentKeeper.GetAllLegacyCommitments(ctx)
	commParams := m.commitmentKeeper.GetLegacyParams(ctx)
	numberOfCommitments := uint64(0)
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
			if committed.Denom == ptypes.Eden && commParams.TotalCommitted.AmountOf(ptypes.Eden).IsPositive() {
				err = m.estakingKeeper.Hooks().BeforeDelegationCreated(ctx, addr, edenValAddr)
				if err != nil {
					return err
				}
				err = m.estakingKeeper.Hooks().AfterDelegationModified(ctx, addr, edenValAddr)
				if err != nil {
					return err
				}
			}
			if committed.Denom == ptypes.EdenB && commParams.TotalCommitted.AmountOf(ptypes.EdenB).IsPositive() {
				err = m.estakingKeeper.Hooks().BeforeDelegationCreated(ctx, addr, edenBValAddr)
				if err != nil {
					return err
				}
				err = m.estakingKeeper.Hooks().AfterDelegationModified(ctx, addr, edenBValAddr)
				if err != nil {
					return err
				}
			}

			// Execute hook for normal amm pool deposit
			poolId, err := ammtypes.GetPoolIdFromShareDenom(committed.Denom)
			if err == nil {
				m.masterchefKeeper.AfterDeposit(ctx, poolId, addr.String(), committed.Amount)
			}

			// Execute hook for stablestake deposit
			if committed.Denom == stablestaketypes.GetShareDenom() {
				m.masterchefKeeper.AfterDeposit(ctx, stablestaketypes.PoolId, addr.String(), committed.Amount)
			}
		}
		numberOfCommitments++
	}
	m.commitmentKeeper.SetParams(ctx, commitmenttypes.Params{
		VestingInfos:        commParams.VestingInfos,
		TotalCommitted:      commParams.TotalCommitted,
		NumberOfCommitments: numberOfCommitments,
	})

	fmt.Println("Finished incentive v11 migration  ...")
	return nil
}
