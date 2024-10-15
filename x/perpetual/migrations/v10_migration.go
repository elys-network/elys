package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	// Update params
	params := m.keeper.GetLegacyParams(ctx)

	newParams := types.Params{
		LeverageMax:                                    params.LeverageMax,
		BorrowInterestRateMax:                          params.BorrowInterestRateMax,
		BorrowInterestRateMin:                          params.BorrowInterestRateMin,
		MinBorrowInterestAmount:                        sdk.ZeroInt(),
		BorrowInterestRateIncrease:                     params.BorrowInterestRateIncrease,
		BorrowInterestRateDecrease:                     params.BorrowInterestRateDecrease,
		HealthGainFactor:                               params.HealthGainFactor,
		EpochLength:                                    params.EpochLength,
		MaxOpenPositions:                               params.MaxOpenPositions,
		PoolOpenThreshold:                              params.PoolOpenThreshold,
		ForceCloseFundPercentage:                       params.ForceCloseFundPercentage,
		ForceCloseFundAddress:                          params.ForceCloseFundAddress,
		IncrementalBorrowInterestPaymentFundPercentage: params.IncrementalBorrowInterestPaymentFundPercentage,
		IncrementalBorrowInterestPaymentFundAddress:    params.IncrementalBorrowInterestPaymentFundAddress,
		SafetyFactor:                                   params.SafetyFactor,
		IncrementalBorrowInterestPaymentEnabled:        params.IncrementalBorrowInterestPaymentEnabled,
		WhitelistingEnabled:                            params.WhitelistingEnabled,
		InvariantCheckEpoch:                            params.InvariantCheckEpoch,
		TakeProfitBorrowInterestRateMin:                params.TakeProfitBorrowInterestRateMin,
		SwapFee:                                        params.SwapFee,
		MaxLimitOrder:                                  params.MaxLimitOrder,
		FixedFundingRate:                               sdk.NewDecWithPrec(30, 2), // set to 30%
	}

	m.keeper.SetParams(ctx, &newParams)

	mtps := m.keeper.GetAllLegacyMTP(ctx)

	ctx.Logger().Info("Migrating positions from legacy to new format")

	for _, mtp := range mtps {
		newMtp := types.MTP{
			Address:                        mtp.Address,
			CollateralAsset:                mtp.CollateralAsset,
			TradingAsset:                   mtp.TradingAsset,
			LiabilitiesAsset:               mtp.LiabilitiesAsset,
			CustodyAsset:                   mtp.CustodyAsset,
			Collateral:                     mtp.Collateral,
			Liabilities:                    mtp.Liabilities,
			BorrowInterestPaidCollateral:   mtp.BorrowInterestPaidCollateral,
			BorrowInterestPaidCustody:      mtp.BorrowInterestPaidCustody,
			BorrowInterestUnpaidCollateral: mtp.BorrowInterestUnpaidCollateral,
			Custody:                        mtp.Custody,
			TakeProfitLiabilities:          mtp.TakeProfitLiabilities,
			TakeProfitCustody:              mtp.TakeProfitCustody,
			MtpHealth:                      mtp.MtpHealth,
			Position:                       mtp.Position,
			Id:                             mtp.Id,
			AmmPoolId:                      mtp.AmmPoolId,
			TakeProfitPrice:                mtp.TakeProfitPrice,
			TakeProfitBorrowRate:           mtp.TakeProfitBorrowRate,
			FundingFeePaidCollateral:       mtp.FundingFeePaidCollateral,
			FundingFeePaidCustody:          mtp.FundingFeePaidCustody,
			FundingFeeReceivedCollateral:   mtp.FundingFeeReceivedCollateral,
			FundingFeeReceivedCustody:      mtp.FundingFeeReceivedCustody,
			OpenPrice:                      mtp.OpenPrice,
			StopLossPrice:                  mtp.StopLossPrice,
			LastInterestCalcTime:           mtp.LastInterestCalcTime,
			LastInterestCalcBlock:          mtp.LastInterestCalcBlock,
			LastFundingCalcTime:            mtp.LastFundingCalcTime,
			LastFundingCalcBlock:           mtp.LastFundingCalcBlock,
		}
		m.keeper.DeleteLegacyMTP(ctx, mtp.Address, mtp.Id)
		m.keeper.SetMTP(ctx, &newMtp)

		baseCurrency, _ := m.keeper.GetBaseCurreny(ctx)
		pool, poolFound := m.keeper.GetPool(ctx, newMtp.AmmPoolId)
		if !poolFound {
			continue
		}
		ammPool, poolErr := m.keeper.GetAmmPool(ctx, newMtp.AmmPoolId, newMtp.TradingAsset)
		if poolErr != nil {
			continue
		}

		m.keeper.CheckAndLiquidateUnhealthyPosition(ctx, &newMtp, pool, ammPool, baseCurrency.Denom, baseCurrency.Decimals)

		pools := m.keeper.GetAllLegacyPools(ctx)

		ctx.Logger().Info("Migrating pool")

		for _, pool := range pools {
			newPool := types.Pool{
				AmmPoolId:                            pool.AmmPoolId,
				Health:                               pool.Health,
				Enabled:                              pool.Enabled,
				Closed:                               pool.Closed,
				BorrowInterestRate:                   pool.BorrowInterestRate,
				PoolAssetsLong:                       pool.PoolAssetsLong,
				PoolAssetsShort:                      pool.PoolAssetsShort,
				LastHeightBorrowInterestRateComputed: pool.LastHeightBorrowInterestRateComputed,
				FundingRate:                          pool.FundingRate,
				FeesCollected:                        sdk.Coins{},
			}
			m.keeper.RemoveLegacyPool(ctx, newPool.AmmPoolId)
			m.keeper.SetPool(ctx, newPool)
		}
	}

	new_mtps := m.keeper.GetAllMTPs(ctx)
	for _, mtp := range new_mtps {
		m.keeper.CheckSamePositionAndConsolidate(ctx, &mtp)
	}

	return nil
}
