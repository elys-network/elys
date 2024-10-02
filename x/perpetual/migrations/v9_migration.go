package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	m.keeper.DeleteAllNegativeCustomMTP(ctx)

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
			ConsolidateLeverage:            mtp.ConsolidateLeverage,
			SumCollateral:                  mtp.SumCollateral,
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
	}

	return nil
}
