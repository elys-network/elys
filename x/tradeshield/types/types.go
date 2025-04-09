package types

import "github.com/osmosis-labs/osmosis/osmomath"

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
	case "stoploss":
		return PerpetualOrderType_STOPLOSSPERP
	default:
		panic("invalid perpetual order type")
	}
}

func (o PerpetualOrder) GetBigDecTriggerPrice() osmomath.BigDec {
	return osmomath.BigDecFromDec(o.TriggerPrice)
}

func (o SpotOrder) GetBigDecOrderPrice() osmomath.BigDec {
	return osmomath.BigDecFromDec(o.OrderPrice)
}
