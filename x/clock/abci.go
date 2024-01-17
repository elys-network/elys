package clock

import (
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/clock/keeper"
	"github.com/elys-network/elys/x/clock/types"
)

// EndBlocker executes on contracts at the end of the block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	message := []byte(types.EndBlockSudoMessage)

	p := k.GetParams(ctx)

	errorExecs := make([]string, 0)

	for idx, addr := range p.ContractAddresses {
		contract, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			errorExecs[idx] = addr
			continue
		}

		childCtx := ctx.WithGasMeter(sdk.NewGasMeter(p.ContractGasLimit))
		_, err = k.GetContractKeeper().Sudo(childCtx, contract, message)
		if err != nil {
			errorExecs = append(errorExecs, err.Error())
			continue
		}
	}

	if len(errorExecs) > 0 {
		log.Printf("[x/clock] Execute Errors: %v", errorExecs)
	}
}
