package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/spf13/cobra"
)

func CmdCreatePendingPerpetualOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pending-perpetual-order [order]",
		Short: "Create a new pending-perpetual-order",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePendingPerpetualOrder(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdatePendingPerpetualOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pending-perpetual-order [id] [order]",
		Short: "Update a pending-perpetual-order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// TODO: Add order price definition in other task
			msg := types.NewMsgUpdatePendingPerpetualOrder(clientCtx.GetFromAddress().String(), id, &types.OrderPrice{})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCancelPerpetualOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-perpetual-orders [ids.json]",
		Short: "Cancel a pending-perpetual-orders by ids",
		Example: "elysd tx perpetual cancel-perpetual-orders ids.json --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := readPositionRequestJSON(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelPerpetualOrders(clientCtx.GetFromAddress().String(), ids)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
