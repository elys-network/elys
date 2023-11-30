package types_test

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateOpenEvent(t *testing.T) {
	// Mock data for testing
	var testMTP = types.MTP{
		Address:                        "elys1x0jyazg9qzys8x9m2x8q3q3x0jyazg9qzys8x9",
		Collateral:                     sdk.ZeroInt(),
		Liabilities:                    sdk.ZeroInt(),
		BorrowInterestPaidCollateral:   sdk.ZeroInt(),
		BorrowInterestPaidCustody:      sdk.ZeroInt(),
		BorrowInterestUnpaidCollateral: sdk.ZeroInt(),
		Custody:                        sdk.ZeroInt(),
		TakeProfitLiabilities:          sdk.ZeroInt(),
		TakeProfitCustody:              sdk.ZeroInt(),
		Leverage:                       sdk.NewDec(10),
		MtpHealth:                      sdk.ZeroDec(),
		Position:                       types.Position_LONG,
		AmmPoolId:                      1,
		Id:                             1,
		ConsolidateLeverage:            sdk.NewDec(10),
		SumCollateral:                  sdk.ZeroInt(),
		TakeProfitPrice:                sdk.NewDec(10),
		TakeProfitBorrowRate:           sdk.OneDec(),
		FundingFeePaidCollateral:       sdk.ZeroInt(),
		FundingFeePaidCustody:          sdk.ZeroInt(),
		FundingFeeReceivedCollateral:   sdk.ZeroInt(),
		FundingFeeReceivedCustody:      sdk.ZeroInt(),
	}

	event := types.GenerateOpenEvent(&testMTP)

	// Assert that the event type is correct
	assert.Equal(t, types.EventOpen, event.Type)

	// Assert that all the attributes are correctly set
	assert.Equal(t, strconv.FormatInt(int64(testMTP.Id), 10), getAttributeValue(event, "id"))
	assert.Equal(t, testMTP.Position.String(), getAttributeValue(event, "position"))
	assert.Equal(t, testMTP.Address, getAttributeValue(event, "address"))
	assert.Equal(t, testMTP.Collateral.String(), getAttributeValue(event, "collateral"))
	assert.Equal(t, testMTP.Custody.String(), getAttributeValue(event, "custody"))
	assert.Equal(t, fmt.Sprintf("%s", testMTP.Leverage), getAttributeValue(event, "leverage"))
	assert.Equal(t, testMTP.Liabilities.String(), getAttributeValue(event, "liabilities"))
	assert.Equal(t, testMTP.BorrowInterestPaidCollateral.String(), getAttributeValue(event, "borrow_interest_paid_collateral"))
	assert.Equal(t, testMTP.BorrowInterestPaidCustody.String(), getAttributeValue(event, "borrow_interest_paid_custody"))
	assert.Equal(t, testMTP.BorrowInterestUnpaidCollateral.String(), getAttributeValue(event, "borrow_interest_unpaid_collateral"))
	assert.Equal(t, testMTP.MtpHealth.String(), getAttributeValue(event, "health"))
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
