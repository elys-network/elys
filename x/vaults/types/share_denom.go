package types

import (
	"fmt"
	"strconv"
	"strings"
)

func GetShareDenomForVault(vaultId uint64) string {
	return "vault/share/" + strconv.FormatUint(vaultId, 10)
}

// GetVaultIDFromPath retrieves the vaultid from the given path in the format "vault/share/vaultid".
func GetVaultIDFromPath(path string) (uint64, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 || parts[0] != "vault" || parts[1] != "share" {
		return 0, fmt.Errorf("invalid path format")
	}
	vaultID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid vault ID: %v", err)
	}
	return vaultID, nil
}
