package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDelegationRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-rewards [delegator-address] [validator-address]",
		Short: "Query delegation-rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDelegatorAddress := args[0]
			reqValidatorAddress := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegationRewardsRequest{

				DelegatorAddress: reqDelegatorAddress,
				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.DelegationRewards(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
