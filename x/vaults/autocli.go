package vaults

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/elys-network/elys/v6/api/elys/vaults"
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
					RpcMethod:      "Vault",
					Use:            "vault <vault-id>",
					Short:          "Shows the vault details",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vault_id"}},
				},
				{
					RpcMethod: "Vaults",
					Use:       "vaults",
					Short:     "Shows all vaults",
				},
				{
					RpcMethod:      "VaultValue",
					Use:            "vault-value <vault-id>",
					Short:          "Shows the value of the vault",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vault_id"}},
				},
				{
					RpcMethod:      "VaultPositions",
					Use:            "vault-positions <vault-id>",
					Short:          "Shows the positions of the vault",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vault_id"}},
				},
				{
					RpcMethod:      "DepositEstimation",
					Use:            "deposit-estimation <vault-id> <amount>",
					Short:          "Shows the estimated deposit amount for a vault",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vault_id"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "WithdrawEstimation",
					Use:            "withdraw-estimation <vault-id> <shares-amount>",
					Short:          "Shows the estimated withdraw amount for a vault",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vault_id"}, {ProtoField: "shares_amount"}},
				},
				// {
				// 	RpcMethod:      "Pnl",
				// 	Use:            "pnl <vault-id> <address>",
				// 	Short:          "Shows the pnl for a address",
				// 	PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "vault_id"}, {ProtoField: "address"}},
				// },
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
				{
					RpcMethod: "PerformActionJoinPool",
					Skip:      true, // Skip autocli generation for this command
				},
				{
					RpcMethod: "PerformActionExitPool",
					Skip:      true, // Skip autocli generation for this command
				},
				{
					RpcMethod: "PerformActionSwapByDenom",
					Skip:      true, // Skip autocli generation for this command
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
