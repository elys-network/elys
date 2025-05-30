package vaults

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/elys-network/elys/v5/api/elys/vaults"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "VaultPositions",
					Use:       "vault-positions <vault-id>",
					Short:     "Shows the positions of the vault",
				},
				{
					RpcMethod: "Vault",
					Use:       "vault <vault-id>",
					Short:     "Shows the vault",
				},
				{
					RpcMethod: "Vaults",
					Use:       "vaults",
					Short:     "Shows the vaults",
				},
				{
					RpcMethod: "VaultValue",
					Use:       "vault-value <vault-id>",
					Short:     "Shows the value of the vault",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateVaultCoins",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateVaultFees",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateVaultLockupPeriod",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateVaultMaxAmountUsd",
					Skip:      true, // skipped because authority gated
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
