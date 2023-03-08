package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tokenomics/types"
	"github.com/spf13/cobra"
)

func CmdShowGenesisInflation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-genesis-inflation",
		Short: "shows genesis-inflation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGenesisInflationRequest{}

			res, err := queryClient.GenesisInflation(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
