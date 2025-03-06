package cli

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/clob/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdPlaceLimitOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "place-limit-order [sub-account-id] [market-id] [price] [quantity] [order-type]",
		Short:   "exit a new pool and withdraw the liquidity from it",
		Example: `elysd tx amm exit-pool 0 1000uatom,1000uusdc 200000000000000000 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			subAccountId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			marketId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			price, err := sdkmath.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}

			quantity, err := sdkmath.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}

			var orderType types.OrderType
			switch args[4] {
			case "limit_buy":
				orderType = types.OrderType_ORDER_TYPE_LIMIT_BUY
			case "limit_sell":
				orderType = types.OrderType_ORDER_TYPE_LIMIT_SELL
			default:
				return errors.New("invalid order type")
			}

			msg := types.MsgPlaceLimitOrder{
				Creator:      clientCtx.GetFromAddress().String(),
				SubAccountId: subAccountId,
				MarketId:     marketId,
				BaseQuantity: quantity,
				OrderType:    orderType,
				Price:        price,
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
