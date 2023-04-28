package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdValidatorOutstandingRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-outstanding-rewards [validator-address]",
		Short: "Query validator-outstanding-rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryValidatorOutstandingRewardsRequest{

				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.ValidatorOutstandingRewards(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
