package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/spf13/cobra"
)

func CmdTotalAirdropClaimed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-airdrop-claimed",
		Short: "Query total airdrop claimed",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryTotalAirDropClaimedRequest{}
			res, err := queryClient.TotalAirdropClaimed(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
