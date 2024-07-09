package ante

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

type MinCommissionDecorator struct {
	sk  *stakingkeeper.Keeper
	bk  bankkeeper.Keeper
	cdc codec.BinaryCodec
	pk  parameterkeeper.Keeper
}

func NewMinCommissionDecorator(cdc codec.BinaryCodec, sk *stakingkeeper.Keeper, bk bankkeeper.Keeper, pk parameterkeeper.Keeper) MinCommissionDecorator {
	return MinCommissionDecorator{cdc: cdc, sk: sk, bk: bk, pk: pk}
}

// getValidator returns the validator belonging to a given bech32 validator address
func (min MinCommissionDecorator) getValidator(ctx sdk.Context, bech32ValAddr string) (stakingtypes.Validator, error) {
	valAddr, err := sdk.ValAddressFromBech32(bech32ValAddr)
	if err != nil {
		return stakingtypes.Validator{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, bech32ValAddr)
	}

	val, found := min.sk.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.Validator{}, errorsmod.Register("ante", 12, "validator does not exist")
	}

	return val, nil
}

func (min MinCommissionDecorator) getTotalDelegatedTokens(ctx sdk.Context) math.Int {
	bondDenom := min.sk.BondDenom(ctx)
	bondedPool := min.sk.GetBondedPool(ctx)
	notBondedPool := min.sk.GetNotBondedPool(ctx)

	notBondedAmount := min.bk.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount
	bondedAmount := min.bk.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount

	return notBondedAmount.Add(bondedAmount)
}

// Returns the projected voting power as a percentage (not a fraction)
func (min MinCommissionDecorator) CalculateValidatorProjectedVotingPower(ctx sdk.Context, delegateAmount sdk.Dec) sdk.Dec {
	totalDelegatedTokens := sdk.NewDecFromInt(min.getTotalDelegatedTokens(ctx))
	// If I am the first validator, then accept 100% voting power
	if totalDelegatedTokens.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}

	projectedTotalDelegatedTokens := totalDelegatedTokens.Add(delegateAmount)
	projectedValidatorTokens := delegateAmount

	// Ensure projectedTotalDelegatedTokens is not zero to avoid division by zero
	if projectedTotalDelegatedTokens.IsZero() {
		return sdk.ZeroDec()
	}

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}

// Returns the projected voting power as a percentage (not a fraction)
func (min MinCommissionDecorator) CalculateDelegateProjectedVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdk.Dec) sdk.Dec {
	validatorTokens := sdk.NewDecFromInt(validator.GetTokens())
	totalDelegatedTokens := sdk.NewDecFromInt(min.getTotalDelegatedTokens(ctx))

	projectedTotalDelegatedTokens := totalDelegatedTokens.Add(delegateAmount)
	projectedValidatorTokens := validatorTokens.Add(delegateAmount)

	// Ensure projectedTotalDelegatedTokens is not zero to avoid division by zero
	if projectedTotalDelegatedTokens.IsZero() {
		return sdk.ZeroDec()
	}

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}

// Returns the projected voting power as a percentage (not a fraction)
func (min MinCommissionDecorator) CalculateRedelegateProjectedVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdk.Dec) sdk.Dec {
	validatorTokens := sdk.NewDecFromInt(validator.GetTokens())
	projectedTotalDelegatedTokens := sdk.NewDecFromInt(min.getTotalDelegatedTokens(ctx)) // no additional delegated tokens

	projectedValidatorTokens := validatorTokens.Add(delegateAmount)

	// Ensure projectedTotalDelegatedTokens is not zero to avoid division by zero
	if projectedTotalDelegatedTokens.IsZero() {
		return sdk.ZeroDec()
	}

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}

func (min MinCommissionDecorator) AnteHandle(
	ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	// Fetch parameter from parameter module
	params := min.pk.GetParams(ctx)
	minCommissionRate := params.MinCommissionRate
	maxVotingPower := params.MaxVotingPower

	validMsg := func(m sdk.Msg) error {
		switch msg := m.(type) {
		case *stakingtypes.MsgCreateValidator:
			// prevent new validators joining the set with
			// commission set below 5%
			if msg.Commission.Rate.LT(minCommissionRate) {
				return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "commission can't be lower than 5%")
			}
			projectedVotingPower := min.CalculateValidatorProjectedVotingPower(ctx, sdk.NewDecFromInt(msg.Value.Amount))
			if projectedVotingPower.GT(maxVotingPower) {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please delegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdk.NewDec(100)), maxVotingPower.Mul(sdk.NewDec(100)))
			}
		case *stakingtypes.MsgEditValidator:
			// if commission rate is nil, it means only
			// other fields are affected - skip
			if msg.CommissionRate == nil {
				break
			}
			if msg.CommissionRate.LT(minCommissionRate) {
				return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "commission can't be lower than 5%")
			}
		case *stakingtypes.MsgDelegate:
			val, err := min.getValidator(ctx, msg.ValidatorAddress)
			if err != nil {
				return err
			}

			projectedVotingPower := min.CalculateDelegateProjectedVotingPower(ctx, val, sdk.NewDecFromInt(msg.Amount.Amount))
			if projectedVotingPower.GT(maxVotingPower) {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please delegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdk.NewDec(100)), maxVotingPower.Mul(sdk.NewDec(100)))
			}
		case *stakingtypes.MsgBeginRedelegate:
			dstVal, err := min.getValidator(ctx, msg.ValidatorDstAddress)
			if err != nil {
				return err
			}

			var delegateAmount sdk.Dec
			if msg.ValidatorSrcAddress == msg.ValidatorDstAddress {
				// This is blocked later on by the SDK. However we may as well calculate the correct projected voting power.
				// Since this is a self redelegation, no additional tokens are delegated to this validator hence delegateAmount = 0
				delegateAmount = sdk.ZeroDec()
			} else {
				delegateAmount = sdk.NewDecFromInt(msg.Amount.Amount)
			}

			projectedVotingPower := min.CalculateRedelegateProjectedVotingPower(ctx, dstVal, delegateAmount)
			if projectedVotingPower.GT(maxVotingPower) {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please redelegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdk.NewDec(100)), maxVotingPower.Mul(sdk.NewDec(100)))
			}
		}

		return nil
	}

	validAuthz := func(execMsg *authz.MsgExec) error {
		for _, v := range execMsg.Msgs {
			var innerMsg sdk.Msg
			err := min.cdc.UnpackAny(v, &innerMsg)
			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "cannot unmarshal authz exec msgs")
			}

			err = validMsg(innerMsg)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, m := range msgs {
		if msg, ok := m.(*authz.MsgExec); ok {
			if err := validAuthz(msg); err != nil {
				return ctx, err
			}
			continue
		}

		// validate normal msgs
		err = validMsg(m)
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
