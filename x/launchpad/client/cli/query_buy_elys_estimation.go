package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/spf13/cobra"
)

func CmdQueryBuyElysEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "buy-elys-estimation [spendingToken] [amount]",
		Short:   "Query BuyElys estimation",
		Example: "elysd q launchpad buy-elys-estimation [spendingToken] [amount]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid spending token amount")
			}

			params := &types.QueryBuyElysEstRequest{
				SpendingToken: args[0],
				TokenAmount:   amount,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.BuyElysEst(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
