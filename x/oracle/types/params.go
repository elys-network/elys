package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance
func NewParams(
	bandEpoch string,
	oracleScriptID OracleScriptID,
	multiplier uint64,
	bandChannelSrc string,
	askCount uint64,
	minCount uint64,
	feeLimit sdk.Coins,
	prepareGas uint64,
	executeGas uint64,
	priceExpiryTime uint64,
	lifeTimeInBlocks uint64,
) Params {
	return Params{
		BandEpoch:         bandEpoch,
		ClientID:          BandPriceClientIDKey,
		OracleScriptID:    uint64(oracleScriptID),
		Multiplier:        multiplier,
		BandChannelSource: bandChannelSrc,
		AskCount:          askCount,
		MinCount:          minCount,
		FeeLimit:          feeLimit,
		PrepareGas:        prepareGas,
		ExecuteGas:        executeGas,
		PriceExpiryTime:   priceExpiryTime,
		LifeTimeInBlocks:  lifeTimeInBlocks,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"band_epoch",
		37,
		18,          // decimal 18
		"channel-1", // used on dockernet
		4,
		3,
		sdk.NewCoins(sdk.NewInt64Coin("uband", 30)),
		600000,
		600000,
		86400, // 1 day old data
		1,     // 1 block old data
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.BandEpoch == "" {
		return fmt.Errorf("band epoch must not be empty: %s", p.BandEpoch)
	}
	if p.ClientID == "" {
		return fmt.Errorf("clientID must not be empty: %s", p.ClientID)
	}
	if p.Multiplier == 0 {
		return fmt.Errorf("multiplier should be positive: %d", p.Multiplier)
	}
	if p.BandChannelSource == "" {
		return fmt.Errorf("channel should not be empty: %s", p.BandChannelSource)
	}
	if p.AskCount == 0 {
		return fmt.Errorf("ask count should not be zero: %d", p.AskCount)
	}
	if p.MinCount == 0 {
		return fmt.Errorf("min count should not be zero: %d", p.MinCount)
	}
	if err := p.FeeLimit.Validate(); err != nil {
		return err
	}

	return nil
}
