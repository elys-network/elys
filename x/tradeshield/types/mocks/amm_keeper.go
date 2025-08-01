// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	math "cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"

	mock "github.com/stretchr/testify/mock"

	osmomath "github.com/osmosis-labs/osmosis/osmomath"

	types "github.com/cosmos/cosmos-sdk/types"
)

// AmmKeeper is an autogenerated mock type for the AmmKeeper type
type AmmKeeper struct {
	mock.Mock
}

type AmmKeeper_Expecter struct {
	mock *mock.Mock
}

func (_m *AmmKeeper) EXPECT() *AmmKeeper_Expecter {
	return &AmmKeeper_Expecter{mock: &_m.Mock}
}

// CalcAmmPrice provides a mock function with given fields: ctx, denom, decimal
func (_m *AmmKeeper) CalcAmmPrice(ctx types.Context, denom string, decimal uint64) osmomath.BigDec {
	ret := _m.Called(ctx, denom, decimal)

	if len(ret) == 0 {
		panic("no return value specified for CalcAmmPrice")
	}

	var r0 osmomath.BigDec
	if rf, ok := ret.Get(0).(func(types.Context, string, uint64) osmomath.BigDec); ok {
		r0 = rf(ctx, denom, decimal)
	} else {
		r0 = ret.Get(0).(osmomath.BigDec)
	}

	return r0
}

// AmmKeeper_CalcAmmPrice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CalcAmmPrice'
type AmmKeeper_CalcAmmPrice_Call struct {
	*mock.Call
}

// CalcAmmPrice is a helper method to define mock.On call
//   - ctx types.Context
//   - denom string
//   - decimal uint64
func (_e *AmmKeeper_Expecter) CalcAmmPrice(ctx interface{}, denom interface{}, decimal interface{}) *AmmKeeper_CalcAmmPrice_Call {
	return &AmmKeeper_CalcAmmPrice_Call{Call: _e.mock.On("CalcAmmPrice", ctx, denom, decimal)}
}

func (_c *AmmKeeper_CalcAmmPrice_Call) Run(run func(ctx types.Context, denom string, decimal uint64)) *AmmKeeper_CalcAmmPrice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string), args[2].(uint64))
	})
	return _c
}

func (_c *AmmKeeper_CalcAmmPrice_Call) Return(_a0 osmomath.BigDec) *AmmKeeper_CalcAmmPrice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmmKeeper_CalcAmmPrice_Call) RunAndReturn(run func(types.Context, string, uint64) osmomath.BigDec) *AmmKeeper_CalcAmmPrice_Call {
	_c.Call.Return(run)
	return _c
}

// CalculateUSDValue provides a mock function with given fields: ctx, denom, amount
func (_m *AmmKeeper) CalculateUSDValue(ctx types.Context, denom string, amount math.Int) osmomath.BigDec {
	ret := _m.Called(ctx, denom, amount)

	if len(ret) == 0 {
		panic("no return value specified for CalculateUSDValue")
	}

	var r0 osmomath.BigDec
	if rf, ok := ret.Get(0).(func(types.Context, string, math.Int) osmomath.BigDec); ok {
		r0 = rf(ctx, denom, amount)
	} else {
		r0 = ret.Get(0).(osmomath.BigDec)
	}

	return r0
}

// AmmKeeper_CalculateUSDValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CalculateUSDValue'
type AmmKeeper_CalculateUSDValue_Call struct {
	*mock.Call
}

// CalculateUSDValue is a helper method to define mock.On call
//   - ctx types.Context
//   - denom string
//   - amount math.Int
func (_e *AmmKeeper_Expecter) CalculateUSDValue(ctx interface{}, denom interface{}, amount interface{}) *AmmKeeper_CalculateUSDValue_Call {
	return &AmmKeeper_CalculateUSDValue_Call{Call: _e.mock.On("CalculateUSDValue", ctx, denom, amount)}
}

func (_c *AmmKeeper_CalculateUSDValue_Call) Run(run func(ctx types.Context, denom string, amount math.Int)) *AmmKeeper_CalculateUSDValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string), args[2].(math.Int))
	})
	return _c
}

func (_c *AmmKeeper_CalculateUSDValue_Call) Return(_a0 osmomath.BigDec) *AmmKeeper_CalculateUSDValue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmmKeeper_CalculateUSDValue_Call) RunAndReturn(run func(types.Context, string, math.Int) osmomath.BigDec) *AmmKeeper_CalculateUSDValue_Call {
	_c.Call.Return(run)
	return _c
}

// SwapByDenom provides a mock function with given fields: ctx, msg
func (_m *AmmKeeper) SwapByDenom(ctx types.Context, msg *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error) {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for SwapByDenom")
	}

	var r0 *ammtypes.MsgSwapByDenomResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)); ok {
		return rf(ctx, msg)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *ammtypes.MsgSwapByDenom) *ammtypes.MsgSwapByDenomResponse); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ammtypes.MsgSwapByDenomResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *ammtypes.MsgSwapByDenom) error); ok {
		r1 = rf(ctx, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AmmKeeper_SwapByDenom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SwapByDenom'
type AmmKeeper_SwapByDenom_Call struct {
	*mock.Call
}

// SwapByDenom is a helper method to define mock.On call
//   - ctx types.Context
//   - msg *ammtypes.MsgSwapByDenom
func (_e *AmmKeeper_Expecter) SwapByDenom(ctx interface{}, msg interface{}) *AmmKeeper_SwapByDenom_Call {
	return &AmmKeeper_SwapByDenom_Call{Call: _e.mock.On("SwapByDenom", ctx, msg)}
}

func (_c *AmmKeeper_SwapByDenom_Call) Run(run func(ctx types.Context, msg *ammtypes.MsgSwapByDenom)) *AmmKeeper_SwapByDenom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*ammtypes.MsgSwapByDenom))
	})
	return _c
}

func (_c *AmmKeeper_SwapByDenom_Call) Return(_a0 *ammtypes.MsgSwapByDenomResponse, _a1 error) *AmmKeeper_SwapByDenom_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AmmKeeper_SwapByDenom_Call) RunAndReturn(run func(types.Context, *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)) *AmmKeeper_SwapByDenom_Call {
	_c.Call.Return(run)
	return _c
}

// NewAmmKeeper creates a new instance of AmmKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAmmKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *AmmKeeper {
	mock := &AmmKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
