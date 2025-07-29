package cli

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perpcli "github.com/elys-network/elys/v7/x/perpetual/client/cli"
	perptypes "github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreatePerpetualOpenOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-perpetual-open-order [position] [leverage] [pool-id] [collateral] [trigger-price]",
		Short: "Create a new perpetual open order",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			position := types.PerpetualPosition(perptypes.GetPositionFromString(args[0]))

			leverage, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			triggerPrice := math.LegacyMustNewDecFromStr(args[4])

			takeProfitPriceStr, err := cmd.Flags().GetString(perpcli.FlagTakeProfitPrice)
			if err != nil {
				return err
			}

			var takeProfitPrice math.LegacyDec
			if takeProfitPriceStr != perptypes.InfinitePriceString {
				takeProfitPrice, err = math.LegacyNewDecFromStr(takeProfitPriceStr)
				if err != nil {
					return errors.New("invalid take profit price")
				}
			} else {
				takeProfitPrice = perptypes.TakeProfitPriceDefault
			}

			stopLossPriceStr, err := cmd.Flags().GetString(perpcli.FlagStopLossPrice)
			if err != nil {
				return err
			}

			var stopLossPrice math.LegacyDec
			if stopLossPriceStr != perptypes.ZeroPriceString {
				stopLossPrice, err = math.LegacyNewDecFromStr(stopLossPriceStr)
				if err != nil {
					return errors.New("invalid stop loss price")
				}
			} else {
				stopLossPrice = perptypes.StopLossPriceDefault
			}

			msg := types.NewMsgCreatePerpetualOpenOrder(
				signer.String(),
				triggerPrice,
				collateral,
				position,
				leverage,
				takeProfitPrice,
				stopLossPrice,
				poolId,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(perpcli.FlagTakeProfitPrice, perptypes.InfinitePriceString, "Optional take profit price")
	cmd.Flags().String(perpcli.FlagStopLossPrice, perptypes.ZeroPriceString, "Optional stop loss price")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreatePerpetualCloseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-perpetual-close-order [trigger-price] [position-id]",
		Short: "Create a new perpetual close order",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr := clientCtx.GetFromAddress().String()

			// trading asset will be filled by message handler
			triggerPrice := math.LegacyMustNewDecFromStr(args[0])

			positionId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePerpetualCloseOrder(addr, triggerPrice, positionId)
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
		Use:   "update-perpetual-order [order-id] [trigger-price]",
		Short: "Update a perpetual order",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orderId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			triggerPrice := math.LegacyMustNewDecFromStr(args[1])

			msg := types.NewMsgUpdatePerpetualOrder(clientCtx.GetFromAddress().String(), orderId, triggerPrice)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCancelPerpetualOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-perpetual-order [order-id]",
		Short: "Broadcast message cancel-perpetual-order",
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

			msg := types.NewMsgCancelPerpetualOrder(
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

func CmdCancelPerpetualOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cancel-perpetual-orders [ids.json]",
		Short:   "Cancel a perpetual orders by ids",
		Example: "elysd tx tradeshield cancel-perpetual-orders ids.json --from=bob --yes --gas=1000000",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orders, err := readPerpetualOrderRequestJSON(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelPerpetualOrders(clientCtx.GetFromAddress().String(), orders)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func readPerpetualOrderRequestJSON(filename string) ([]types.PerpetualOrderPoolKey, error) {
	var orders []types.PerpetualOrderPoolKey
	bz, err := os.ReadFile(filename)
	if err != nil {
		return []types.PerpetualOrderPoolKey{}, err
	}
	err = json.Unmarshal(bz, &orders)
	if err != nil {
		return []types.PerpetualOrderPoolKey{}, err
	}

	return orders, nil
}

func CmdCancelAllPerpetualOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cancel-all-perpetual-orders",
		Short:   "Cancel all pending perpetual orders for the user",
		Example: "elysd tx tradeshield cancel-all-perpetual-orders --from=bob --yes --gas=1000000",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelAllPerpetualOrders(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
