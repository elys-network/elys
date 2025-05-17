package app

import (
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateUpgradeVersion(t *testing.T) {
	// mainnet
	version.Version = "v4.0.0"
	require.Equal(t, "v4", generateUpgradeVersion())

	version.Version = "v4.1.0"
	require.Equal(t, "v4.1", generateUpgradeVersion())

	version.Version = "v4.2.1"
	require.Equal(t, "v4.2", generateUpgradeVersion())

	// SUT
	version.Version = "v"
	require.Equal(t, "v999999", generateUpgradeVersion())

	// testnet
	version.Version = "v4.0.0-rc0"
	require.Equal(t, "v4-rc0", generateUpgradeVersion())

	version.Version = "v4.0.0-rc1"
	require.Equal(t, "v4-rc1", generateUpgradeVersion())

	// devnet
	version.Version = "2e561b347baaad345e9a73f4cbfcdcbf3c958d20"
	require.Equal(t, "2e561b347baaad345e9a73f4cbfcdcbf3c958d20", generateUpgradeVersion())
}
