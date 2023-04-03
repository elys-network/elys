package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tokenomics/types"
	"github.com/spf13/cobra"
)

func CmdListAirdrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-airdrop",
		Short: "list all airdrop",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllAirdropRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AirdropAll(context.Background(), params)
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

func CmdShowAirdrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-airdrop [intent]",
		Short: "shows a airdrop",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIntent := args[0]

			params := &types.QueryGetAirdropRequest{
				Intent: argIntent,
			}

			res, err := queryClient.Airdrop(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
