package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdShowFeeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fee-info [date]",
		Short: "Query show-fee-info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDate := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryShowFeeInfoRequest{

				Date: reqDate,
			}

			res, err := queryClient.ShowFeeInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
