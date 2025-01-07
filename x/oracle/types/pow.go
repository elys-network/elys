package types

import (
	elystypes "github.com/elys-network/elys/types"
)

func Pow10(decimal uint64) elystypes.Dec34 {
	value := elystypes.OneDec34()
	for i := 0; i < int(decimal); i++ {
		value = value.Mul(elystypes.NewDec34FromInt64(10))
	}
	return value
}
