// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	math "cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"

	mock "github.com/stretchr/testify/mock"

	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CloseLongChecker is an autogenerated mock type for the CloseLongChecker type
type CloseLongChecker struct {
	mock.Mock
}

type CloseLongChecker_Expecter struct {
	mock *mock.Mock
}

func (_m *CloseLongChecker) EXPECT() *CloseLongChecker_Expecter {
	return &CloseLongChecker_Expecter{mock: &_m.Mock}
}

// EstimateAndRepay provides a mock function with given fields: ctx, mtp, pool, ammPool, amount, baseCurrency
func (_m *CloseLongChecker) EstimateAndRepay(ctx types.Context, mtp perpetualtypes.MTP, pool perpetualtypes.Pool, ammPool ammtypes.Pool, amount math.Int, baseCurrency string) (math.Int, error) {
	ret := _m.Called(ctx, mtp, pool, ammPool, amount, baseCurrency)

	var r0 math.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, perpetualtypes.MTP, perpetualtypes.Pool, ammtypes.Pool, math.Int, string) (math.Int, error)); ok {
		return rf(ctx, mtp, pool, ammPool, amount, baseCurrency)
	}
	if rf, ok := ret.Get(0).(func(types.Context, perpetualtypes.MTP, perpetualtypes.Pool, ammtypes.Pool, math.Int, string) math.Int); ok {
		r0 = rf(ctx, mtp, pool, ammPool, amount, baseCurrency)
	} else {
		r0 = ret.Get(0).(math.Int)
	}

	if rf, ok := ret.Get(1).(func(types.Context, perpetualtypes.MTP, perpetualtypes.Pool, ammtypes.Pool, math.Int, string) error); ok {
		r1 = rf(ctx, mtp, pool, ammPool, amount, baseCurrency)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CloseLongChecker_EstimateAndRepay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EstimateAndRepay'
type CloseLongChecker_EstimateAndRepay_Call struct {
	*mock.Call
}

// EstimateAndRepay is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp perpetualtypes.MTP
//   - pool perpetualtypes.Pool
//   - ammPool ammtypes.Pool
//   - amount math.Int
//   - baseCurrency string
func (_e *CloseLongChecker_Expecter) EstimateAndRepay(ctx interface{}, mtp interface{}, pool interface{}, ammPool interface{}, amount interface{}, baseCurrency interface{}) *CloseLongChecker_EstimateAndRepay_Call {
	return &CloseLongChecker_EstimateAndRepay_Call{Call: _e.mock.On("EstimateAndRepay", ctx, mtp, pool, ammPool, amount, baseCurrency)}
}

func (_c *CloseLongChecker_EstimateAndRepay_Call) Run(run func(ctx types.Context, mtp perpetualtypes.MTP, pool perpetualtypes.Pool, ammPool ammtypes.Pool, amount math.Int, baseCurrency string)) *CloseLongChecker_EstimateAndRepay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(perpetualtypes.MTP), args[2].(perpetualtypes.Pool), args[3].(ammtypes.Pool), args[4].(math.Int), args[5].(string))
	})
	return _c
}

func (_c *CloseLongChecker_EstimateAndRepay_Call) Return(_a0 math.Int, _a1 error) *CloseLongChecker_EstimateAndRepay_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CloseLongChecker_EstimateAndRepay_Call) RunAndReturn(run func(types.Context, perpetualtypes.MTP, perpetualtypes.Pool, ammtypes.Pool, math.Int, string) (math.Int, error)) *CloseLongChecker_EstimateAndRepay_Call {
	_c.Call.Return(run)
	return _c
}

// GetAmmPool provides a mock function with given fields: ctx, poolId, tradingAsset
func (_m *CloseLongChecker) GetAmmPool(ctx types.Context, poolId uint64, tradingAsset string) (ammtypes.Pool, error) {
	ret := _m.Called(ctx, poolId, tradingAsset)

	var r0 ammtypes.Pool
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64, string) (ammtypes.Pool, error)); ok {
		return rf(ctx, poolId, tradingAsset)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64, string) ammtypes.Pool); ok {
		r0 = rf(ctx, poolId, tradingAsset)
	} else {
		r0 = ret.Get(0).(ammtypes.Pool)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64, string) error); ok {
		r1 = rf(ctx, poolId, tradingAsset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CloseLongChecker_GetAmmPool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAmmPool'
type CloseLongChecker_GetAmmPool_Call struct {
	*mock.Call
}

// GetAmmPool is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
//   - tradingAsset string
func (_e *CloseLongChecker_Expecter) GetAmmPool(ctx interface{}, poolId interface{}, tradingAsset interface{}) *CloseLongChecker_GetAmmPool_Call {
	return &CloseLongChecker_GetAmmPool_Call{Call: _e.mock.On("GetAmmPool", ctx, poolId, tradingAsset)}
}

func (_c *CloseLongChecker_GetAmmPool_Call) Run(run func(ctx types.Context, poolId uint64, tradingAsset string)) *CloseLongChecker_GetAmmPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64), args[2].(string))
	})
	return _c
}

func (_c *CloseLongChecker_GetAmmPool_Call) Return(_a0 ammtypes.Pool, _a1 error) *CloseLongChecker_GetAmmPool_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CloseLongChecker_GetAmmPool_Call) RunAndReturn(run func(types.Context, uint64, string) (ammtypes.Pool, error)) *CloseLongChecker_GetAmmPool_Call {
	_c.Call.Return(run)
	return _c
}

// GetMTP provides a mock function with given fields: ctx, mtpAddress, id
func (_m *CloseLongChecker) GetMTP(ctx types.Context, mtpAddress string, id uint64) (perpetualtypes.MTP, error) {
	ret := _m.Called(ctx, mtpAddress, id)

	var r0 perpetualtypes.MTP
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, string, uint64) (perpetualtypes.MTP, error)); ok {
		return rf(ctx, mtpAddress, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string, uint64) perpetualtypes.MTP); ok {
		r0 = rf(ctx, mtpAddress, id)
	} else {
		r0 = ret.Get(0).(perpetualtypes.MTP)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string, uint64) error); ok {
		r1 = rf(ctx, mtpAddress, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CloseLongChecker_GetMTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMTP'
type CloseLongChecker_GetMTP_Call struct {
	*mock.Call
}

// GetMTP is a helper method to define mock.On call
//   - ctx types.Context
//   - mtpAddress string
//   - id uint64
func (_e *CloseLongChecker_Expecter) GetMTP(ctx interface{}, mtpAddress interface{}, id interface{}) *CloseLongChecker_GetMTP_Call {
	return &CloseLongChecker_GetMTP_Call{Call: _e.mock.On("GetMTP", ctx, mtpAddress, id)}
}

func (_c *CloseLongChecker_GetMTP_Call) Run(run func(ctx types.Context, mtpAddress string, id uint64)) *CloseLongChecker_GetMTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string), args[2].(uint64))
	})
	return _c
}

func (_c *CloseLongChecker_GetMTP_Call) Return(_a0 perpetualtypes.MTP, _a1 error) *CloseLongChecker_GetMTP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CloseLongChecker_GetMTP_Call) RunAndReturn(run func(types.Context, string, uint64) (perpetualtypes.MTP, error)) *CloseLongChecker_GetMTP_Call {
	_c.Call.Return(run)
	return _c
}

// GetPool provides a mock function with given fields: ctx, poolId
func (_m *CloseLongChecker) GetPool(ctx types.Context, poolId uint64) (perpetualtypes.Pool, bool) {
	ret := _m.Called(ctx, poolId)

	var r0 perpetualtypes.Pool
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, uint64) (perpetualtypes.Pool, bool)); ok {
		return rf(ctx, poolId)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64) perpetualtypes.Pool); ok {
		r0 = rf(ctx, poolId)
	} else {
		r0 = ret.Get(0).(perpetualtypes.Pool)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64) bool); ok {
		r1 = rf(ctx, poolId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// CloseLongChecker_GetPool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPool'
type CloseLongChecker_GetPool_Call struct {
	*mock.Call
}

// GetPool is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
func (_e *CloseLongChecker_Expecter) GetPool(ctx interface{}, poolId interface{}) *CloseLongChecker_GetPool_Call {
	return &CloseLongChecker_GetPool_Call{Call: _e.mock.On("GetPool", ctx, poolId)}
}

func (_c *CloseLongChecker_GetPool_Call) Run(run func(ctx types.Context, poolId uint64)) *CloseLongChecker_GetPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64))
	})
	return _c
}

func (_c *CloseLongChecker_GetPool_Call) Return(val perpetualtypes.Pool, found bool) *CloseLongChecker_GetPool_Call {
	_c.Call.Return(val, found)
	return _c
}

func (_c *CloseLongChecker_GetPool_Call) RunAndReturn(run func(types.Context, uint64) (perpetualtypes.Pool, bool)) *CloseLongChecker_GetPool_Call {
	_c.Call.Return(run)
	return _c
}

// SettleBorrowInterest provides a mock function with given fields: ctx, mtp, pool, ammPool
func (_m *CloseLongChecker) SettleBorrowInterest(ctx types.Context, mtp *perpetualtypes.MTP, pool *perpetualtypes.Pool, ammPool ammtypes.Pool) (math.Int, error) {
	ret := _m.Called(ctx, mtp, pool, ammPool)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualtypes.MTP, *perpetualtypes.Pool, ammtypes.Pool) error); ok {
		r0 = rf(ctx, mtp, pool, ammPool)
	} else {
		r0 = ret.Error(0)
	}

	return types.ZeroInt(), r0
}

// CloseLongChecker_SettleBorrowInterest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SettleBorrowInterest'
type CloseLongChecker_SettleBorrowInterest_Call struct {
	*mock.Call
}

// SettleBorrowInterest is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp *perpetualtypes.MTP
//   - pool *perpetualtypes.Pool
//   - ammPool ammtypes.Pool
func (_e *CloseLongChecker_Expecter) SettleBorrowInterest(ctx interface{}, mtp interface{}, pool interface{}, ammPool interface{}) *CloseLongChecker_SettleBorrowInterest_Call {
	return &CloseLongChecker_SettleBorrowInterest_Call{Call: _e.mock.On("SettleBorrowInterest", ctx, mtp, pool, ammPool)}
}

func (_c *CloseLongChecker_SettleBorrowInterest_Call) Run(run func(ctx types.Context, mtp *perpetualtypes.MTP, pool *perpetualtypes.Pool, ammPool ammtypes.Pool)) *CloseLongChecker_SettleBorrowInterest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*perpetualtypes.MTP), args[2].(*perpetualtypes.Pool), args[3].(ammtypes.Pool))
	})
	return _c
}

func (_c *CloseLongChecker_SettleBorrowInterest_Call) Return(_a0 error) *CloseLongChecker_SettleBorrowInterest_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CloseLongChecker_SettleBorrowInterest_Call) RunAndReturn(run func(types.Context, *perpetualtypes.MTP, *perpetualtypes.Pool, ammtypes.Pool) error) *CloseLongChecker_SettleBorrowInterest_Call {
	_c.Call.Return(run)
	return _c
}

// TakeOutCustody provides a mock function with given fields: ctx, mtp, pool, amount
func (_m *CloseLongChecker) TakeOutCustody(ctx types.Context, mtp perpetualtypes.MTP, pool *perpetualtypes.Pool, amount math.Int) error {
	ret := _m.Called(ctx, mtp, pool, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, perpetualtypes.MTP, *perpetualtypes.Pool, math.Int) error); ok {
		r0 = rf(ctx, mtp, pool, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CloseLongChecker_TakeOutCustody_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TakeOutCustody'
type CloseLongChecker_TakeOutCustody_Call struct {
	*mock.Call
}

// TakeOutCustody is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp perpetualtypes.MTP
//   - pool *perpetualtypes.Pool
//   - amount math.Int
func (_e *CloseLongChecker_Expecter) TakeOutCustody(ctx interface{}, mtp interface{}, pool interface{}, amount interface{}) *CloseLongChecker_TakeOutCustody_Call {
	return &CloseLongChecker_TakeOutCustody_Call{Call: _e.mock.On("TakeOutCustody", ctx, mtp, pool, amount)}
}

func (_c *CloseLongChecker_TakeOutCustody_Call) Run(run func(ctx types.Context, mtp perpetualtypes.MTP, pool *perpetualtypes.Pool, amount math.Int)) *CloseLongChecker_TakeOutCustody_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(perpetualtypes.MTP), args[2].(*perpetualtypes.Pool), args[3].(math.Int))
	})
	return _c
}

func (_c *CloseLongChecker_TakeOutCustody_Call) Return(_a0 error) *CloseLongChecker_TakeOutCustody_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CloseLongChecker_TakeOutCustody_Call) RunAndReturn(run func(types.Context, perpetualtypes.MTP, *perpetualtypes.Pool, math.Int) error) *CloseLongChecker_TakeOutCustody_Call {
	_c.Call.Return(run)
	return _c
}

// NewCloseLongChecker creates a new instance of CloseLongChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCloseLongChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *CloseLongChecker {
	mock := &CloseLongChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
