package cli

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateSpotOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-spot-order [order-type] [order-amount] [order-target-denom] [order-price]",
		Short: "Create a new spot order",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr := clientCtx.GetFromAddress().String()
			orderType := types.GetSpotOrderTypeFromString(args[0])
			orderAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			orderTargetDenom := args[2]
			orderPrice := types.OrderPrice{
				BaseDenom:  orderAmount.Denom,
				QuoteDenom: orderTargetDenom,
				Rate:       sdk.MustNewDecFromStr(args[3]),
			}

			msg := types.NewMsgCreateSpotOrder(addr, orderType, orderPrice, orderAmount, orderTargetDenom)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateSpotOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-spot-order [id] [order]",
		Short: "Update a spot order",
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
			msg := types.NewMsgUpdateSpotOrder(clientCtx.GetFromAddress().String(), id, &types.OrderPrice{})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCancelSpotOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-spot-order [order-id]",
		Short: "Broadcast message cancel-spot-order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argOrderId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelSpotOrder(
				clientCtx.GetFromAddress().String(),
				argOrderId,
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

func CmdCancelSpotOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cancel-spot-orders [ids.json]",
		Short:   "Cancel spot-orders",
		Example: "elysd tx perpetual cancel-spot-orders ids.json --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := readPositionRequestJSON(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelSpotOrders(clientCtx.GetFromAddress().String(), ids)
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
