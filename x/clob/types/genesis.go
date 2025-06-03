package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                  DefaultParams(),
		SubAccounts:             nil,
		PerpetualMarkets:        nil,
		PerpetualMarketCounters: nil,
		Perpetuals:              nil,
		PerpetualOwners:         nil,
		OrderBooks:              nil,
		OrderOwners:             nil,
		TwapPrices:              nil,
		FundingRates:            nil,
		PerpetualADLs:           nil,
	}
}

func (gs GenesisState) Validate() error {
	return nil
}
