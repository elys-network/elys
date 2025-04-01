package cli

import (
	"cosmossdk.io/math"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/clob/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdPlaceMarketOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "place-market-order [market-id] [quantity] [order-type]",
		Short:   "exit a new pool and withdraw the liquidity from it",
		Example: `elysd tx amm exit-pool 0 1000uatom,1000uusdc 200000000000000000 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			marketId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			quantity, ok := math.NewIntFromString(args[2])
			if !ok {
				return errors.New("invalid quantity")
			}
			var orderType types.OrderType
			switch args[3] {
			case "market_buy":
				orderType = types.OrderType_ORDER_TYPE_MARKET_BUY
			case "market_sell":
				orderType = types.OrderType_ORDER_TYPE_MARKET_SELL
			default:
				return errors.New("invalid order type")
			}

			msg := types.MsgPlaceMarketOrder{
				Creator:      clientCtx.GetFromAddress().String(),
				MarketId:     marketId,
				BaseQuantity: quantity,
				OrderType:    orderType,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
