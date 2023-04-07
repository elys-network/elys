package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdShowCommitments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-commitments [creator]",
		Short: "Query show-commitments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqCreator := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryShowCommitmentsRequest{

				Creator: reqCreator,
			}

			res, err := queryClient.ShowCommitments(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
