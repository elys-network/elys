package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/elys-network/elys/v7/x/estaking/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		CmdWithdrawAllRewards(),
		CmdWithdrawElysStakingRewards(),
		CmdUnjailGovernor(),
	)

	return cmd
}

func CmdWithdrawAllRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-all-rewards",
		Short: "Withdraw all rewards for delegations and Eden/EdenB commit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw all rewards for delegations and Eden/EdenB commit,
Example:
$ %s tx estaking withdraw-all-rewards --from mykey 
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			msg := &types.MsgWithdrawAllRewards{
				DelegatorAddress: delAddr.String(),
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUnjailGovernor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail-governor",
		Short: "Unjail a jailed governor",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unjail a jailed governor,
Example:
$ %s tx estaking unjail-governor --from governor 
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			address := clientCtx.GetFromAddress()

			msg := &types.MsgUnjailGovernor{
				Address: address.String(),
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdWithdrawElysStakingRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-elys-staking-rewards",
		Short: "Withdraw rewards for delegations",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw rewards for delegations,
Example:
$ %s tx estaking withdraw-elys-staking-rewards --from mykey 
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			msg := &types.MsgWithdrawAllRewards{
				DelegatorAddress: delAddr.String(),
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
