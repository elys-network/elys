package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdLiquidationPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidation-price [address] [position_id]",
		Short: "Query liquidation price",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]
			_, err = sdk.AccAddressFromBech32(reqAddress)
			if err != nil {
				return err
			}

			reqId := args[1]
			positionId, err := strconv.ParseUint(reqId, 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLiquidationPriceRequest{
				Address:    reqAddress,
				PositionId: positionId,
			}

			res, err := queryClient.LiquidationPrice(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
