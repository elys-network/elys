package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/launchpad/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
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

	cmd.AddCommand(CmdBuyElys())
	cmd.AddCommand(CmdReturnElys())
	cmd.AddCommand(CmdWithdrawRaised())

	return cmd
}

func CmdBuyElys() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "buy-elys [spendingToken] [amount]",
		Short:   "Buy elys from launchpad",
		Example: `elysd tx launchpad buy-elys uatom 100000000000 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			spendingToken := args[0]
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid spending token amount")
			}

			msg := types.NewMsgBuyElys(
				clientCtx.GetFromAddress().String(),
				spendingToken,
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdReturnElys() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "return-elys [orderId] [amount]",
		Short:   "Return elys bought from launchpad",
		Example: `elysd tx launchpad return-elys 1 100000000000 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orderId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid spending token amount")
			}

			msg := types.NewMsgReturnElys(
				clientCtx.GetFromAddress().String(),
				uint64(orderId),
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdWithdrawRaised() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-raised [amounts]",
		Short:   "Withdraw raised amount from launchpad",
		Example: `elysd tx launchpad withdraw-raised 100000000000uatom --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amounts, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawRaised(
				clientCtx.GetFromAddress().String(),
				amounts,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
