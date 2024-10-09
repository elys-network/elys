package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPendingPerpetualOrderForAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-perpetual-order-for-address [address]",
		Short: "Query pending-perpetual-order-for-address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]
			reqStatusStr, err := cmd.Flags().GetString(flagStatus)
			if err != nil {
				return err
			}

			params := &types.QueryPendingPerpetualOrderForAddressRequest{
				Address: reqAddress,
			}

			if reqStatusStr != "" {
				reqStatus, ok := types.Status_value[reqStatusStr]
				if !ok {
					return types.ErrInvalidStatus
				}
				params = &types.QueryPendingPerpetualOrderForAddressRequest{
					Address: reqAddress,
					Status:  types.Status(reqStatus),
				}
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.PendingPerpetualOrderForAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
