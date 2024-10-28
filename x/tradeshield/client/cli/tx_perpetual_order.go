package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perptypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/spf13/cobra"
)

func CmdCreatePerpetualOpenOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-perpetual-open-order [order-type] [position] [leverage] [trading-asset] [collateral] [trigger-price] [take-profit-price] [stop-loss-price]",
		Short: "Create a new perpetual open order",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr := clientCtx.GetFromAddress().String()
			orderType := types.GetPerpetualOrderTypeFromString(args[0])
			position := types.PerpetualPosition(perptypes.GetPositionFromString(args[1]))
			leverage := sdk.MustNewDecFromStr(args[2])
			tradingAsset := args[3]
			collateral, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}
			triggerPrice := types.OrderPrice{
				BaseDenom:  collateral.Denom,
				QuoteDenom: tradingAsset,
				Rate:       sdk.MustNewDecFromStr(args[5]),
			}
			takeProfitPrice := sdk.MustNewDecFromStr(args[6])
			stopLossPrice := sdk.MustNewDecFromStr(args[7])

			msg := types.NewMsgCreatePerpetualOpenOrder(addr, orderType, triggerPrice, collateral, tradingAsset, position, leverage, takeProfitPrice, stopLossPrice)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreatePerpetualCloseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-perpetual-close-order [order-type] [trigger-price] [position-id]",
		Short: "Create a new perpetual close order",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr := clientCtx.GetFromAddress().String()
			orderType := types.GetPerpetualOrderTypeFromString(args[0])

			// base and quote denom will be filled by message handler
			triggerPrice := types.OrderPrice{
				BaseDenom:  "",
				QuoteDenom: "",
				Rate:       sdk.MustNewDecFromStr(args[1]),
			}

			positionId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePerpetualCloseOrder(addr, orderType, triggerPrice, positionId)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdatePerpetualOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-perpetual-order [id] [order]",
		Short: "Update a perpetual order",
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
			msg := types.NewMsgUpdatePerpetualOrder(clientCtx.GetFromAddress().String(), id, &types.OrderPrice{})
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
		Use:     "cancel-perpetual-orders [ids.json]",
		Short:   "Cancel a perpetual orders by ids",
		Example: "elysd tx perpetual cancel-perpetual-orders ids.json --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000",
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
