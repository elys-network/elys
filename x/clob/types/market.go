package types

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	MinimumIMRByMMR = math.LegacyMustNewDecFromStr("1.5")
	MaximumIMRByMMR = math.LegacyMustNewDecFromStr("2.5")
)

func NewPerpetualMarket(id uint64, baseDenom, quoteDenom, admin string, imr, mmr, makerFee, takerFee, liquidationFeeShare, minPriceTickSize, minQuantityTickSize, minNotional, maxFRChange, maxFR math.LegacyDec, twapPriceWindow uint64, status PerpetualMarketStatus) PerpetualMarket {
	return PerpetualMarket{
		Id:                      id,
		BaseDenom:               baseDenom,
		QuoteDenom:              quoteDenom,
		InitialMarginRatio:      imr,
		MaintenanceMarginRatio:  mmr,
		MakerFeeRate:            makerFee,
		TakerFeeRate:            takerFee,
		LiquidationFeeShareRate: liquidationFeeShare,
		Status:                  status,
		MinPriceTickSize:        minPriceTickSize,
		MinQuantityTickSize:     minQuantityTickSize,
		MinNotional:             minNotional,
		Admin:                   admin,
		TotalOpen:               math.LegacyZeroDec(),
		MaxAbsFundingRateChange: maxFRChange,
		MaxAbsFundingRate:       maxFR,
		TwapPricesWindow:        twapPriceWindow,
	}
}

func (market PerpetualMarket) ValidateOpenPositionRequest(marketId uint64, price, quantity math.LegacyDec, isMarketOrder bool) error {
	if market.Id != marketId {
		return ErrPerpetualMarketNotFound
	}
	if price.Mul(quantity).LT(market.MinNotional) {
		return errors.New("trade value less than minimum notional value")
	}
	if quantity.LT(market.MinQuantityTickSize) {
		return errors.New("quantity less than minimum quantity tick size")
	}
	if !quantity.Quo(market.MinQuantityTickSize).IsInteger() {
		return errors.New("quantity is not of proper tick size")
	}
	if !price.Quo(market.MinPriceTickSize).IsInteger() {
		return errors.New("price is not of proper tick size")
	}
	return nil
}

func (market PerpetualMarket) GetAccount() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/perpetual/%d", market.Id))
}

func (market PerpetualMarket) GetInsuranceAccount() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/perpetual/insurance/%d", market.Id))
}

func (market *PerpetualMarket) UpdateTotalOpenInterest(buyerBefore, sellerBefore, tradeSize math.LegacyDec) {
	if tradeSize.LTE(math.LegacyZeroDec()) {
		panic("trade size cannot be 0 or negative")
	}

	// Calculate final positions for buyer and seller
	buyerAfter := buyerBefore.Add(tradeSize)
	sellerAfter := sellerBefore.Sub(tradeSize)

	// Determine if each party increased or decreased their exposure magnitude
	buyerExposureIncreased := buyerAfter.Abs().GT(buyerBefore.Abs())
	sellerExposureIncreased := sellerAfter.Abs().GT(sellerBefore.Abs())
	buyerExposureDecreased := buyerAfter.Abs().LT(buyerBefore.Abs())
	sellerExposureDecreased := sellerAfter.Abs().LT(sellerBefore.Abs())

	deltaOI := math.LegacyZeroDec()

	// OI increases only if BOTH increase exposure magnitude
	if buyerExposureIncreased && sellerExposureIncreased {
		deltaOI = tradeSize
	}

	// OI decreases only if BOTH decrease exposure magnitude
	if buyerExposureDecreased && sellerExposureDecreased {
		deltaOI = tradeSize.Neg()
	}

	// In all other cases (one increases, one decreases), OI is unchanged (deltaOI remains zero)

	market.TotalOpen = market.TotalOpen.Add(deltaOI)
}
