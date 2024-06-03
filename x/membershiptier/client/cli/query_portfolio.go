package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/elys-network/elys/x/membershiptier/types"
)

func CmdListPortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-portfolio",
		Short: "list all portfolio",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPortfolioRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PortfolioAll(cmd.Context(), params)
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

func CmdShowPortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-portfolio [user] [asset-type]",
		Short: "shows a portfolio",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argUser := args[0]
			argType := args[1]

			params := &types.QueryGetPortfolioRequest{
				User:      argUser,
				AssetType: argType,
			}

			res, err := queryClient.Portfolio(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
