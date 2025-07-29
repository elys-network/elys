package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdClaimAllUserRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim-all-user-rewards [flags]",
		Short:   "Claims rewards from all leveragelp positions, capped to maxPageLimit positions",
		Example: `elysd tx leveragelp claim-all-user-rewards --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			msg := &types.MsgClaimAllUserRewards{
				Sender: signer.String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
