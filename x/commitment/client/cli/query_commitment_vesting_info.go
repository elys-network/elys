package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCommitmentVestingInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commitment-vesting-info [address]",
		Short: "Query commitment-vesting-info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCommitmentVestingInfoRequest{

				Address: reqAddress,
			}

			res, err := queryClient.CommitmentVestingInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
