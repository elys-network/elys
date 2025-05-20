package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v4/x/oracle/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateAssetInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-asset-info [denom] [display] [band-ticker] [elys-ticker] [decimal]",
		Short: "create a new asset info",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argDisplay := args[1]
			argBandTicker := args[2]
			argElysTicker := args[3]
			argDecimal, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAssetInfo(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argDisplay,
				argBandTicker,
				argElysTicker,
				argDecimal,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
