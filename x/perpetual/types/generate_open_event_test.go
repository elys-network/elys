package types_test

import (
	sdkmath "cosmossdk.io/math"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateOpenEvent(t *testing.T) {
	// Mock data for testing
	testMTP := types.MTP{
		Address:                       "elys1x0jyazg9qzys8x9m2x8q3q3x0jyazg9qzys8x9",
		CollateralAsset:               "uusdc",
		TradingAsset:                  "uatom",
		LiabilitiesAsset:              "uusdc",
		CustodyAsset:                  "uatom",
		Collateral:                    sdkmath.OneInt(),
		Liabilities:                   sdkmath.OneInt(),
		BorrowInterestPaidCustody:     sdkmath.OneInt(),
		BorrowInterestUnpaidLiability: sdkmath.OneInt(),
		Custody:                       sdkmath.OneInt(),
		TakeProfitLiabilities:         sdkmath.OneInt(),
		TakeProfitCustody:             sdkmath.OneInt(),
		MtpHealth:                     sdkmath.LegacyZeroDec(),
		Position:                      types.Position_LONG,
		Id:                            1,
		AmmPoolId:                     1,
		TakeProfitPrice:               sdkmath.LegacyNewDec(10),
		TakeProfitBorrowFactor:        sdkmath.LegacyOneDec(),
		FundingFeePaidCustody:         sdkmath.OneInt(),
		FundingFeeReceivedCustody:     sdkmath.OneInt(),
		OpenPrice:                     sdkmath.LegacyNewDec(10),
		StopLossPrice:                 sdkmath.LegacyNewDec(0),
	}

	event := types.GenerateOpenEvent(&testMTP)

	// Assert that the event type is correct
	assert.Equal(t, types.EventOpen, event.Type)

	// Assert that all the attributes are correctly set
	assert.Equal(t, testMTP.Address, getAttributeValue(event, "address"))
	assert.Equal(t, testMTP.CollateralAsset, getAttributeValue(event, "collateral_asset"))
	assert.Equal(t, testMTP.TradingAsset, getAttributeValue(event, "trading_asset"))
	assert.Equal(t, testMTP.LiabilitiesAsset, getAttributeValue(event, "liabilities_asset"))
	assert.Equal(t, testMTP.CustodyAsset, getAttributeValue(event, "custody_asset"))
	assert.Equal(t, testMTP.Collateral.String(), getAttributeValue(event, "collateral"))
	assert.Equal(t, testMTP.Liabilities.String(), getAttributeValue(event, "liabilities"))
	assert.Equal(t, testMTP.BorrowInterestPaidCustody.String(), getAttributeValue(event, "borrow_interest_paid_custody"))
	assert.Equal(t, testMTP.BorrowInterestUnpaidLiability.String(), getAttributeValue(event, "borrow_interest_unpaid_liability"))
	assert.Equal(t, testMTP.Custody.String(), getAttributeValue(event, "custody"))
	assert.Equal(t, testMTP.TakeProfitLiabilities.String(), getAttributeValue(event, "take_profit_liabilities"))
	assert.Equal(t, testMTP.TakeProfitCustody.String(), getAttributeValue(event, "take_profit_custody"))
	assert.Equal(t, testMTP.MtpHealth.String(), getAttributeValue(event, "mtp_health"))
	assert.Equal(t, testMTP.Position.String(), getAttributeValue(event, "position"))
	assert.Equal(t, strconv.FormatInt(int64(testMTP.Id), 10), getAttributeValue(event, "id"))
	assert.Equal(t, strconv.FormatInt(int64(testMTP.AmmPoolId), 10), getAttributeValue(event, "amm_pool_id"))
	assert.Equal(t, testMTP.TakeProfitPrice.String(), getAttributeValue(event, "take_profit_price"))
	assert.Equal(t, testMTP.TakeProfitBorrowFactor.String(), getAttributeValue(event, "take_profit_borrow_factor"))
	assert.Equal(t, testMTP.FundingFeePaidCustody.String(), getAttributeValue(event, "funding_fee_paid_custody"))
	assert.Equal(t, testMTP.FundingFeeReceivedCustody.String(), getAttributeValue(event, "funding_fee_received_custody"))
	assert.Equal(t, testMTP.OpenPrice.String(), getAttributeValue(event, "open_price"))
}

// Helper function to get attribute value from an event
func getAttributeValue(event sdk.Event, key string) string {
	for _, attr := range event.Attributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}
