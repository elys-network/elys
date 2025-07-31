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
	vaultIdModuleName := ModuleName + "/account/" + vaultIdStr
	return vaultIdModuleName
}

func NewVaultRewardCollectorAddress(vaultId uint64) sdk.AccAddress {
	vaultIdModuleName := GetVaultIdModuleName(vaultId) + "/reward-collector"
	return address.Module(vaultIdModuleName)
}

func NewVaultRewardCollectorAddressString(vaultId uint64) string {
	vaultIdModuleName := GetVaultIdModuleName(vaultId) + "/reward-collector"
	return vaultIdModuleName
}
