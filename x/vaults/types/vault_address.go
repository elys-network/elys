package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func NewVaultAddress(vaultId uint64) sdk.AccAddress {
	poolIdModuleName := GetVaultIdModuleName(vaultId)
	return address.Module(poolIdModuleName)
}

func GetVaultIdModuleName(vaultId uint64) string {
	vaultIdStr := strconv.FormatUint(vaultId, 10)
	vaultIdModuleName := ModuleName + "/pool/account/" + vaultIdStr
	return vaultIdModuleName
}
