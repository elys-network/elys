package types_test

import (
	"testing"

	"github.com/elys-network/elys/v7/x/vaults/types"
	"github.com/stretchr/testify/require"
)

func TestVaultActionValidation(t *testing.T) {
	// Test vault with no allowed actions (should allow all actions)
	vault1 := &types.Vault{
		Id:             1,
		AllowedActions: []string{},
	}

	// Test vault with specific allowed actions
	vault2 := &types.Vault{
		Id:             2,
		AllowedActions: []string{"swap", "join_pool"},
	}

	// Test vault with case-insensitive actions
	vault3 := &types.Vault{
		Id:             3,
		AllowedActions: []string{"SWAP", "JOIN_POOL"},
	}

	tests := []struct {
		name        string
		vault       *types.Vault
		action      string
		shouldAllow bool
	}{
		{
			name:        "vault with no restrictions should allow swap",
			vault:       vault1,
			action:      "swap",
			shouldAllow: true,
		},
		{
			name:        "vault with no restrictions should allow join_pool",
			vault:       vault1,
			action:      "join_pool",
			shouldAllow: true,
		},
		{
			name:        "vault with specific actions should allow swap",
			vault:       vault2,
			action:      "swap",
			shouldAllow: true,
		},
		{
			name:        "vault with specific actions should allow join_pool",
			vault:       vault2,
			action:      "join_pool",
			shouldAllow: true,
		},
		{
			name:        "vault with specific actions should not allow exit_pool",
			vault:       vault2,
			action:      "exit_pool",
			shouldAllow: false,
		},
		{
			name:        "case insensitive should work for swap",
			vault:       vault3,
			action:      "swap",
			shouldAllow: true,
		},
		{
			name:        "case insensitive should work for SWAP",
			vault:       vault3,
			action:      "SWAP",
			shouldAllow: true,
		},
		{
			name:        "whitespace should be trimmed",
			vault:       vault2,
			action:      "  swap  ",
			shouldAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := tt.vault.IsActionAllowed(tt.action)
			require.Equal(t, tt.shouldAllow, allowed)

			if tt.shouldAllow {
				err := tt.vault.ValidateActionAllowed(tt.action)
				require.NoError(t, err)
			} else {
				err := tt.vault.ValidateActionAllowed(tt.action)
				require.Error(t, err)
				require.Contains(t, err.Error(), "is not allowed")
			}
		})
	}
}
