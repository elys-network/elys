package cli

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/spf13/cobra"
)

// TODO: Add message in other task
func CmdCreatePendingSpotOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pending-spot-order [order]",
		Short: "Create a new pending-spot-order",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePendingSpotOrder(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdatePendingSpotOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pending-spot-order [id] [order]",
		Short: "Update a pending-spot-order",
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
			msg := types.NewMsgUpdatePendingSpotOrder(clientCtx.GetFromAddress().String(), id, &types.OrderPrice{})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCancelSpotOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-pending-spot-orders [ids.json]",
		Short: "Cancel pending-spot-orders",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := readPositionRequestJSON(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelSpotOrders(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func readPositionRequestJSON(filename string) ([]uint64, error) {
	var positions []uint64
	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		return []uint64{}, err
	}
	err = json.Unmarshal(bz, &positions)
	if err != nil {
		return []uint64{}, err
	}

	return positions, nil
}