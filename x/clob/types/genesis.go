package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:           DefaultParams(),
		SubAccounts:      nil,
		PerpetualMarkets: nil,
		Perpetuals:       nil,
		PerpetualOwners:  nil,
		OrderBooks:       nil,
		TwapPrices:       nil,
	}
}

func (gs GenesisState) Validate() error {
	return nil
}
