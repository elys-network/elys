package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tokenomics/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdListTimeBasedInflation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-time-based-inflation",
		Short: "list all time-based-inflation",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllTimeBasedInflationRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.TimeBasedInflationAll(context.Background(), params)
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

func CmdShowTimeBasedInflation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-time-based-inflation [start-block-height] [end-block-height]",
		Short: "shows a time-based-inflation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argStartBlockHeight, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argEndBlockHeight, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryGetTimeBasedInflationRequest{
				StartBlockHeight: argStartBlockHeight,
				EndBlockHeight:   argEndBlockHeight,
			}

			res, err := queryClient.TimeBasedInflation(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
