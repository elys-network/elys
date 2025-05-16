package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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

	val, err := min.sk.GetValidator(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, errorsmod.Wrapf(err, "validator not found in ante MinCommissionDecorator")
	}

	return val, nil
}

func (min MinCommissionDecorator) getTotalBondedTokens(ctx sdk.Context) sdkmath.Int {
	bondDenom, err := min.sk.BondDenom(ctx)
	if err != nil {
		panic(err)
	}
	bondedPool := min.sk.GetBondedPool(ctx)
	bondedAmount := min.bk.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount

	return bondedAmount
}

// Returns the projected voting power as a percentage (not a fraction)
func (min MinCommissionDecorator) CalculateValidatorProjectedVotingPower(ctx sdk.Context, delegateAmount sdkmath.LegacyDec) sdkmath.LegacyDec {

	bondedAmt := min.getTotalBondedTokens(ctx).ToLegacyDec()
	// If I am the first validator, then accept 100% voting power
	if bondedAmt.LTE(sdkmath.LegacyZeroDec()) {
		return sdkmath.LegacyZeroDec()
	}

	projectedTotalDelegatedTokens := bondedAmt.Add(delegateAmount)
	projectedValidatorTokens := delegateAmount

	// Ensure projectedTotalDelegatedTokens is not zero to avoid division by zero
	if projectedTotalDelegatedTokens.IsZero() {
		return sdkmath.LegacyZeroDec()
	}

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}

// Returns the projected voting power as a percentage (not a fraction)
func (min MinCommissionDecorator) CalculateDelegateProjectedVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdkmath.LegacyDec) sdkmath.LegacyDec {
	validatorTokens := validator.GetTokens().ToLegacyDec()
	bondedAmt := min.getTotalBondedTokens(ctx).ToLegacyDec()

	projectedTotalDelegatedTokens := bondedAmt.Add(delegateAmount)
	projectedValidatorTokens := validatorTokens.Add(delegateAmount)

	// Ensure projectedTotalDelegatedTokens is not zero to avoid division by zero
	if projectedTotalDelegatedTokens.IsZero() {
		return sdkmath.LegacyZeroDec()
	}

	return projectedValidatorTokens.Quo(projectedTotalDelegatedTokens)
}

// Returns the projected voting power as a percentage (not a fraction)
func (min MinCommissionDecorator) CalculateRedelegateProjectedVotingPower(ctx sdk.Context, validator stakingtypes.ValidatorI, delegateAmount sdkmath.LegacyDec) sdkmath.LegacyDec {
	validatorTokens := validator.GetTokens().ToLegacyDec()
	projectedTotalDelegatedTokens := min.getTotalBondedTokens(ctx).ToLegacyDec() // no additional delegated tokens

	projectedValidatorTokens := validatorTokens.Add(delegateAmount)

	// Ensure projectedTotalDelegatedTokens is not zero to avoid division by zero
	if projectedTotalDelegatedTokens.IsZero() {
		return sdkmath.LegacyZeroDec()
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
			projectedVotingPower := min.CalculateValidatorProjectedVotingPower(ctx, msg.Value.Amount.ToLegacyDec())
			if projectedVotingPower.GT(maxVotingPower) {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please delegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdkmath.LegacyNewDec(100)), maxVotingPower.Mul(sdkmath.LegacyNewDec(100)))
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

			projectedVotingPower := min.CalculateDelegateProjectedVotingPower(ctx, val, msg.Amount.Amount.ToLegacyDec())
			if projectedVotingPower.GT(maxVotingPower) {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please delegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdkmath.LegacyNewDec(100)), maxVotingPower.Mul(sdkmath.LegacyNewDec(100)))
			}
		case *stakingtypes.MsgBeginRedelegate:
			dstVal, err := min.getValidator(ctx, msg.ValidatorDstAddress)
			if err != nil {
				return err
			}

			var delegateAmount sdkmath.LegacyDec
			if msg.ValidatorSrcAddress == msg.ValidatorDstAddress {
				// This is blocked later on by the SDK. However we may as well calculate the correct projected voting power.
				// Since this is a self redelegation, no additional tokens are delegated to this validator hence delegateAmount = 0
				delegateAmount = sdkmath.LegacyZeroDec()
			} else {
				delegateAmount = msg.Amount.Amount.ToLegacyDec()
			}

			projectedVotingPower := min.CalculateRedelegateProjectedVotingPower(ctx, dstVal, delegateAmount)
			if projectedVotingPower.GT(maxVotingPower) {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please redelegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdkmath.LegacyNewDec(100)), maxVotingPower.Mul(sdkmath.LegacyNewDec(100)))
			}
		case *commitmenttypes.MsgStake:
			if msg.Asset == ptypes.Elys {
				val, err := min.getValidator(ctx, msg.ValidatorAddress)
				if err != nil {
					return err
				}

				projectedVotingPower := min.CalculateDelegateProjectedVotingPower(ctx, val, msg.Amount.ToLegacyDec())
				if projectedVotingPower.GT(maxVotingPower) {
					return errorsmod.Wrapf(
						sdkerrors.ErrInvalidRequest,
						"This validator has a voting power of %s%%. Delegations not allowed to a validator whose post-delegation voting power is more than %s%%. Please delegate to a validator with less bonded tokens", projectedVotingPower.Mul(sdkmath.LegacyNewDec(100)), maxVotingPower.Mul(sdkmath.LegacyNewDec(100)))
				}
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
