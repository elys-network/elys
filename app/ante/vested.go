package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
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
	_, accountType1 := delegatorAcc.(*vestingtypes.BaseVestingAccount)
	_, accountType2 := delegatorAcc.(*vestingtypes.ContinuousVestingAccount)
	_, accountType3 := delegatorAcc.(*vestingtypes.DelayedVestingAccount)
	_, accountType4 := delegatorAcc.(*vestingtypes.PeriodicVestingAccount)
	if accountType1 || accountType2 || accountType3 || accountType4 {
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
