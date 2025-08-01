package cli

import (
	"errors"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v7/x/clob/types"
	"github.com/spf13/cobra"
)

func CmdPlaceLimitOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "place-limit-order [market-id] [price] [quantity] [order-type]",
		Short:   "exit a new pool and withdraw the liquidity from it",
		Example: `elysd tx clob place-limit-order 1 10 5 limit_buy --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			marketId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			price, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			quantity, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}
			var orderType types.OrderType
			switch args[3] {
			case "limit_buy":
				orderType = types.OrderType_ORDER_TYPE_LIMIT_BUY
			case "limit_sell":
				orderType = types.OrderType_ORDER_TYPE_LIMIT_SELL
			default:
				return errors.New("invalid order type")
			}

			isolatedOrder, err := cmd.Flags().GetBool(FlagIsolatedOrder)
			if err != nil {
				return err
			}

			msg := types.MsgPlaceLimitOrder{
				Creator:      clientCtx.GetFromAddress().String(),
				MarketId:     marketId,
				BaseQuantity: quantity,
				OrderType:    orderType,
				Price:        price,
				IsIsolated:   isolatedOrder,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Bool(FlagIsolatedOrder, true, "place an isolated order")

	return cmd
}
