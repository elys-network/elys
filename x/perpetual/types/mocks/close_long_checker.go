// Code generated by mockery v2.45.0. DO NOT EDIT.

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

// EstimateAndRepay provides a mock function with given fields: ctx, mtp, pool, ammPool, amount, baseCurrency
func (_m *CloseLongChecker) EstimateAndRepay(ctx types.Context, mtp perpetualtypes.MTP, pool perpetualtypes.Pool, ammPool ammtypes.Pool, amount math.Int, baseCurrency string) (math.Int, error) {
	ret := _m.Called(ctx, mtp, pool, ammPool, amount, baseCurrency)

	if len(ret) == 0 {
		panic("no return value specified for EstimateAndRepay")
	}

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

// GetAmmPool provides a mock function with given fields: ctx, poolId, tradingAsset
func (_m *CloseLongChecker) GetAmmPool(ctx types.Context, poolId uint64, tradingAsset string) (ammtypes.Pool, error) {
	ret := _m.Called(ctx, poolId, tradingAsset)

	if len(ret) == 0 {
		panic("no return value specified for GetAmmPool")
	}

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

// GetMTP provides a mock function with given fields: ctx, mtpAddress, id
func (_m *CloseLongChecker) GetMTP(ctx types.Context, mtpAddress types.AccAddress, id uint64) (perpetualtypes.MTP, error) {
	ret := _m.Called(ctx, mtpAddress, id)

	if len(ret) == 0 {
		panic("no return value specified for GetMTP")
	}

	var r0 perpetualtypes.MTP
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, uint64) (perpetualtypes.MTP, error)); ok {
		return rf(ctx, mtpAddress, id)
	}
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, uint64) perpetualtypes.MTP); ok {
		r0 = rf(ctx, mtpAddress, id)
	} else {
		r0 = ret.Get(0).(perpetualtypes.MTP)
	}

	if rf, ok := ret.Get(1).(func(types.Context, types.AccAddress, uint64) error); ok {
		r1 = rf(ctx, mtpAddress, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPool provides a mock function with given fields: ctx, poolId
func (_m *CloseLongChecker) GetPool(ctx types.Context, poolId uint64) (perpetualtypes.Pool, bool) {
	ret := _m.Called(ctx, poolId)

	if len(ret) == 0 {
		panic("no return value specified for GetPool")
	}

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

// HandleBorrowInterest provides a mock function with given fields: ctx, mtp, pool, ammPool
func (_m *CloseLongChecker) SettleBorrowInterest(ctx types.Context, mtp *perpetualtypes.MTP, pool *perpetualtypes.Pool, ammPool ammtypes.Pool) (math.Int, error) {
	ret := _m.Called(ctx, mtp, pool, ammPool)

	if len(ret) == 0 {
		panic("no return value specified for SettleBorrowInterest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualtypes.MTP, *perpetualtypes.Pool, ammtypes.Pool) error); ok {
		r0 = rf(ctx, mtp, pool, ammPool)
	} else {
		r0 = ret.Error(0)
	}

	return types.ZeroInt(), r0
}

// TakeOutCustody provides a mock function with given fields: ctx, mtp, pool, amount
func (_m *CloseLongChecker) TakeOutCustody(ctx types.Context, mtp perpetualtypes.MTP, pool *perpetualtypes.Pool, amount math.Int) error {
	ret := _m.Called(ctx, mtp, pool, amount)

	if len(ret) == 0 {
		panic("no return value specified for TakeOutCustody")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, perpetualtypes.MTP, *perpetualtypes.Pool, math.Int) error); ok {
		r0 = rf(ctx, mtp, pool, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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
