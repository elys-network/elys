package types

func GetSpotOrderTypeFromString(s string) SpotOrderType {
	switch s {
	case "stoploss":
		return SpotOrderType_STOPLOSS
	case "limitsell":
		return SpotOrderType_LIMITSELL
	case "limitbuy":
		return SpotOrderType_LIMITBUY
	case "marketbuy":
		return SpotOrderType_MARKETBUY
	default:
		panic("invalid spot order type")
	}
}

func GetPerpetualOrderTypeFromString(s string) PerpetualOrderType {
	switch s {
	case "limitopen":
		return PerpetualOrderType_LIMITOPEN
	case "limitclose":
		return PerpetualOrderType_LIMITCLOSE
	case "marketopen":
		return PerpetualOrderType_MARKETOPEN
	case "marketclose":
		return PerpetualOrderType_MARKETCLOSE
	case "stoploss":
		return PerpetualOrderType_STOPLOSSPERP
	default:
		panic("invalid perpetual order type")
	}
}
