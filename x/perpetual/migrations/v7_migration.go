package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	mtps := m.keeper.GetAllLegacyMTP(ctx)

	ctx.Logger().Info("Migrating positions from legacy to new format")

	for _, mtp := range mtps {
		new_mtp := types.MTP{
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
			Leverage:                       mtp.Leverage,
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
			StopLossPrice:                  sdk.NewDec(0),
		}
		m.keeper.DeleteLegacyMTP(ctx, mtp.Address, mtp.Id)
		m.keeper.SetMTP(ctx, &new_mtp)
	}
	return nil
}