package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/spf13/cobra"
)

func CmdAirdrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-airdrop [address]",
		Short: "Query airdrops",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			req := &types.QueryAirDropRequest{
				Address: args[0],
			}
			res, err := queryClient.AirDrop(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
