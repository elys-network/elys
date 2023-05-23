package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/parameter/types"
	"github.com/spf13/cobra"
)

func CmdListAnteHandlerParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-ante-handler-param",
		Short: "list all ante-handler-param",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllAnteHandlerParamRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AnteHandlerParamAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
