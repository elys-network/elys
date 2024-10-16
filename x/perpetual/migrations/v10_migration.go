package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	m.keeper.DeleteAllFundingRate(ctx)
	m.keeper.DeleteAllInterestRate(ctx)

	newParams := types.Params{
		LeverageMax:                                    sdk.MustNewDecFromStr("25"),
		BorrowInterestRateMax:                          sdk.NewDecWithPrec(27, 7), // 0.000002700000000000
		BorrowInterestRateMin:                          sdk.NewDecWithPrec(3, 8),  // 0.000000030000000000
		MinBorrowInterestAmount:                        sdk.ZeroInt(),
		BorrowInterestRateIncrease:                     sdk.NewDecWithPrec(33, 10), // 0.000000003300000000
		BorrowInterestRateDecrease:                     sdk.NewDecWithPrec(33, 10), // 0.000000003300000000
		HealthGainFactor:                               sdk.NewDecWithPrec(22, 8),  // 0.000000220000000000
		EpochLength:                                    1,
		MaxOpenPositions:                               3000,
		PoolOpenThreshold:                              sdk.MustNewDecFromStr("0.65"),
		ForceCloseFundPercentage:                       sdk.MustNewDecFromStr("1"),
		ForceCloseFundAddress:                          "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l",
		IncrementalBorrowInterestPaymentFundPercentage: sdk.MustNewDecFromStr("0.1"),
		IncrementalBorrowInterestPaymentFundAddress:    "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l",
		SafetyFactor:                                   sdk.MustNewDecFromStr("1.05"),
		IncrementalBorrowInterestPaymentEnabled:        true,
		WhitelistingEnabled:                            false,
		InvariantCheckEpoch:                            "day",
		TakeProfitBorrowInterestRateMin:                sdk.MustNewDecFromStr("1"),
		SwapFee:                                        sdk.MustNewDecFromStr("0.001"),
		MaxLimitOrder:                                  500,
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
