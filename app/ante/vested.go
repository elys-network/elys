package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// CustomAnteHandlerDecorator checks for the use of vested tokens in various operations
type VestedAnteHandlerDecorator struct {
	ak authante.AccountKeeper
	bk bankkeeper.Keeper
}

func NewVestedAnteHandlerDecorator(ak authante.AccountKeeper, bk bankkeeper.Keeper) VestedAnteHandlerDecorator {
	return VestedAnteHandlerDecorator{ak: ak, bk: bk}
}

func (cad VestedAnteHandlerDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	// Iterate over all messages in the transaction
	for _, m := range tx.GetMsgs() {
		switch msg := m.(type) {
		case *stakingtypes.MsgDelegate:
			if err := cad.checkDelegation(ctx, msg); err != nil {
				return ctx, err
			}
		}
	}
	return next(ctx, tx, simulate)
}

func (cad VestedAnteHandlerDecorator) checkDelegation(ctx sdk.Context, msg *stakingtypes.MsgDelegate) error {
	// Retrieve the delegator account
	delegatorAcc := cad.ak.GetAccount(ctx, sdk.AccAddress(msg.DelegatorAddress))

	// Check if the account is a vesting account

	if _, ok := delegatorAcc.(types.VestingAccount); ok {
		// Calculate the spendable coins
		spendableCoins := cad.bk.SpendableCoins(ctx, sdk.AccAddress(msg.DelegatorAddress))

		// Ensure the amount to be delegated is less than or equal to the spendable coins
		if msg.Amount.Amount.GT(spendableCoins.AmountOf(msg.Amount.Denom)) {
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot delegate vested tokens")
		}
	}

	// For non-vesting accounts, no additional checks are required
	return nil
}
