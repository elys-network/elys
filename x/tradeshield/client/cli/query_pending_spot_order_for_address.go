package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPendingSpotOrderForAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-spot-order-for-address [address] [status]",
		Short: "Query pending-spot-order-for-address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]
			// TODO: If status is empty
			reqStatus, ok := types.Status_value[args[1]]
			if !ok {
				return types.ErrInvalidStatus
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPendingSpotOrderForAddressRequest{

				Address: reqAddress,
				Status:  types.Status(reqStatus),
			}

			res, err := queryClient.PendingSpotOrderForAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
