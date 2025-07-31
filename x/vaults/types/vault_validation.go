package types

import (
	"fmt"
	"strings"
)

// Action constants for vault operations
const (
	ActionSwap     = "swap"
	ActionJoinPool = "join_pool"
	ActionExitPool = "exit_pool"
)

// IsActionAllowed checks if a specific action is allowed for the vault
func (v *Vault) IsActionAllowed(action string) bool {
	action = strings.ToLower(strings.TrimSpace(action))
	for _, allowedAction := range v.AllowedActions {
		if strings.ToLower(strings.TrimSpace(allowedAction)) == action {
			return true
		}
	}
	return false
}

// ValidateActionAllowed checks if an action is allowed and returns an error if not
func (v *Vault) ValidateActionAllowed(action string) error {
	if !v.IsActionAllowed(action) {
		return fmt.Errorf("action '%s' is not allowed for vault %d. Allowed actions: %v", action, v.Id, v.AllowedActions)
	}
	return nil
}
