package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdListFeeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-fee-info",
		Short: "Query list-fee-info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryListFeeInfoRequest{}

			res, err := queryClient.ListFeeInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
