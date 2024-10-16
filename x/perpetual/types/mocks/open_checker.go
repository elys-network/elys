// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	ammtypes "github.com/elys-network/elys/x/amm/types"
	mock "github.com/stretchr/testify/mock"

	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// OpenChecker is an autogenerated mock type for the OpenChecker type
type OpenChecker struct {
	mock.Mock
}

type OpenChecker_Expecter struct {
	mock *mock.Mock
}

func (_m *OpenChecker) EXPECT() *OpenChecker_Expecter {
	return &OpenChecker_Expecter{mock: &_m.Mock}
}

// CheckMaxOpenPositions provides a mock function with given fields: ctx
func (_m *OpenChecker) CheckMaxOpenPositions(ctx types.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CheckMaxOpenPositions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_CheckMaxOpenPositions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckMaxOpenPositions'
type OpenChecker_CheckMaxOpenPositions_Call struct {
	*mock.Call
}

// CheckMaxOpenPositions is a helper method to define mock.On call
//   - ctx types.Context
func (_e *OpenChecker_Expecter) CheckMaxOpenPositions(ctx interface{}) *OpenChecker_CheckMaxOpenPositions_Call {
	return &OpenChecker_CheckMaxOpenPositions_Call{Call: _e.mock.On("CheckMaxOpenPositions", ctx)}
}

func (_c *OpenChecker_CheckMaxOpenPositions_Call) Run(run func(ctx types.Context)) *OpenChecker_CheckMaxOpenPositions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context))
	})
	return _c
}

func (_c *OpenChecker_CheckMaxOpenPositions_Call) Return(_a0 error) *OpenChecker_CheckMaxOpenPositions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckMaxOpenPositions_Call) RunAndReturn(run func(types.Context) error) *OpenChecker_CheckMaxOpenPositions_Call {
	_c.Call.Return(run)
	return _c
}

// CheckPoolHealth provides a mock function with given fields: ctx, poolId
func (_m *OpenChecker) CheckPoolHealth(ctx types.Context, poolId uint64) error {
	ret := _m.Called(ctx, poolId)

	if len(ret) == 0 {
		panic("no return value specified for CheckPoolHealth")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64) error); ok {
		r0 = rf(ctx, poolId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_CheckPoolHealth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckPoolHealth'
type OpenChecker_CheckPoolHealth_Call struct {
	*mock.Call
}

// CheckPoolHealth is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
func (_e *OpenChecker_Expecter) CheckPoolHealth(ctx interface{}, poolId interface{}) *OpenChecker_CheckPoolHealth_Call {
	return &OpenChecker_CheckPoolHealth_Call{Call: _e.mock.On("CheckPoolHealth", ctx, poolId)}
}

func (_c *OpenChecker_CheckPoolHealth_Call) Run(run func(ctx types.Context, poolId uint64)) *OpenChecker_CheckPoolHealth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64))
	})
	return _c
}

func (_c *OpenChecker_CheckPoolHealth_Call) Return(_a0 error) *OpenChecker_CheckPoolHealth_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckPoolHealth_Call) RunAndReturn(run func(types.Context, uint64) error) *OpenChecker_CheckPoolHealth_Call {
	_c.Call.Return(run)
	return _c
}

// CheckSameAssetPosition provides a mock function with given fields: ctx, msg
func (_m *OpenChecker) CheckSameAssetPosition(ctx types.Context, msg *perpetualtypes.MsgOpen) *perpetualtypes.MTP {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for CheckSameAssetPosition")
	}

	var r0 *perpetualtypes.MTP
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualtypes.MsgOpen) *perpetualtypes.MTP); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualtypes.MTP)
		}
	}

	return r0
}

// OpenChecker_CheckSameAssetPosition_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckSameAssetPosition'
type OpenChecker_CheckSameAssetPosition_Call struct {
	*mock.Call
}

// CheckSameAssetPosition is a helper method to define mock.On call
//   - ctx types.Context
//   - msg *perpetualtypes.MsgOpen
func (_e *OpenChecker_Expecter) CheckSameAssetPosition(ctx interface{}, msg interface{}) *OpenChecker_CheckSameAssetPosition_Call {
	return &OpenChecker_CheckSameAssetPosition_Call{Call: _e.mock.On("CheckSameAssetPosition", ctx, msg)}
}

func (_c *OpenChecker_CheckSameAssetPosition_Call) Run(run func(ctx types.Context, msg *perpetualtypes.MsgOpen)) *OpenChecker_CheckSameAssetPosition_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*perpetualtypes.MsgOpen))
	})
	return _c
}

func (_c *OpenChecker_CheckSameAssetPosition_Call) Return(_a0 *perpetualtypes.MTP) *OpenChecker_CheckSameAssetPosition_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckSameAssetPosition_Call) RunAndReturn(run func(types.Context, *perpetualtypes.MsgOpen) *perpetualtypes.MTP) *OpenChecker_CheckSameAssetPosition_Call {
	_c.Call.Return(run)
	return _c
}

// CheckUserAuthorization provides a mock function with given fields: ctx, msg
func (_m *OpenChecker) CheckUserAuthorization(ctx types.Context, msg *perpetualtypes.MsgOpen) error {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for CheckUserAuthorization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualtypes.MsgOpen) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_CheckUserAuthorization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckUserAuthorization'
type OpenChecker_CheckUserAuthorization_Call struct {
	*mock.Call
}

// CheckUserAuthorization is a helper method to define mock.On call
//   - ctx types.Context
//   - msg *perpetualtypes.MsgOpen
func (_e *OpenChecker_Expecter) CheckUserAuthorization(ctx interface{}, msg interface{}) *OpenChecker_CheckUserAuthorization_Call {
	return &OpenChecker_CheckUserAuthorization_Call{Call: _e.mock.On("CheckUserAuthorization", ctx, msg)}
}

func (_c *OpenChecker_CheckUserAuthorization_Call) Run(run func(ctx types.Context, msg *perpetualtypes.MsgOpen)) *OpenChecker_CheckUserAuthorization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*perpetualtypes.MsgOpen))
	})
	return _c
}

func (_c *OpenChecker_CheckUserAuthorization_Call) Return(_a0 error) *OpenChecker_CheckUserAuthorization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckUserAuthorization_Call) RunAndReturn(run func(types.Context, *perpetualtypes.MsgOpen) error) *OpenChecker_CheckUserAuthorization_Call {
	_c.Call.Return(run)
	return _c
}

// EmitOpenEvent provides a mock function with given fields: ctx, mtp
func (_m *OpenChecker) EmitOpenEvent(ctx types.Context, mtp *perpetualtypes.MTP) {
	_m.Called(ctx, mtp)
}

// OpenChecker_EmitOpenEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EmitOpenEvent'
type OpenChecker_EmitOpenEvent_Call struct {
	*mock.Call
}

// EmitOpenEvent is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp *perpetualtypes.MTP
func (_e *OpenChecker_Expecter) EmitOpenEvent(ctx interface{}, mtp interface{}) *OpenChecker_EmitOpenEvent_Call {
	return &OpenChecker_EmitOpenEvent_Call{Call: _e.mock.On("EmitOpenEvent", ctx, mtp)}
}

func (_c *OpenChecker_EmitOpenEvent_Call) Run(run func(ctx types.Context, mtp *perpetualtypes.MTP)) *OpenChecker_EmitOpenEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*perpetualtypes.MTP))
	})
	return _c
}

func (_c *OpenChecker_EmitOpenEvent_Call) Return() *OpenChecker_EmitOpenEvent_Call {
	_c.Call.Return()
	return _c
}

func (_c *OpenChecker_EmitOpenEvent_Call) RunAndReturn(run func(types.Context, *perpetualtypes.MTP)) *OpenChecker_EmitOpenEvent_Call {
	_c.Call.Return(run)
	return _c
}

// GetMaxOpenPositions provides a mock function with given fields: ctx
func (_m *OpenChecker) GetMaxOpenPositions(ctx types.Context) uint64 {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetMaxOpenPositions")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func(types.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// OpenChecker_GetMaxOpenPositions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMaxOpenPositions'
type OpenChecker_GetMaxOpenPositions_Call struct {
	*mock.Call
}

// GetMaxOpenPositions is a helper method to define mock.On call
//   - ctx types.Context
func (_e *OpenChecker_Expecter) GetMaxOpenPositions(ctx interface{}) *OpenChecker_GetMaxOpenPositions_Call {
	return &OpenChecker_GetMaxOpenPositions_Call{Call: _e.mock.On("GetMaxOpenPositions", ctx)}
}

func (_c *OpenChecker_GetMaxOpenPositions_Call) Run(run func(ctx types.Context)) *OpenChecker_GetMaxOpenPositions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context))
	})
	return _c
}

func (_c *OpenChecker_GetMaxOpenPositions_Call) Return(_a0 uint64) *OpenChecker_GetMaxOpenPositions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_GetMaxOpenPositions_Call) RunAndReturn(run func(types.Context) uint64) *OpenChecker_GetMaxOpenPositions_Call {
	_c.Call.Return(run)
	return _c
}

// GetOpenMTPCount provides a mock function with given fields: ctx
func (_m *OpenChecker) GetOpenMTPCount(ctx types.Context) uint64 {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetOpenMTPCount")
	}

	var r0 uint64
	if rf, ok := ret.Get(0).(func(types.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// OpenChecker_GetOpenMTPCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOpenMTPCount'
type OpenChecker_GetOpenMTPCount_Call struct {
	*mock.Call
}

// GetOpenMTPCount is a helper method to define mock.On call
//   - ctx types.Context
func (_e *OpenChecker_Expecter) GetOpenMTPCount(ctx interface{}) *OpenChecker_GetOpenMTPCount_Call {
	return &OpenChecker_GetOpenMTPCount_Call{Call: _e.mock.On("GetOpenMTPCount", ctx)}
}

func (_c *OpenChecker_GetOpenMTPCount_Call) Run(run func(ctx types.Context)) *OpenChecker_GetOpenMTPCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context))
	})
	return _c
}

func (_c *OpenChecker_GetOpenMTPCount_Call) Return(_a0 uint64) *OpenChecker_GetOpenMTPCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_GetOpenMTPCount_Call) RunAndReturn(run func(types.Context) uint64) *OpenChecker_GetOpenMTPCount_Call {
	_c.Call.Return(run)
	return _c
}

// OpenDefineAssets provides a mock function with given fields: ctx, poolId, msg, baseCurrency, isBroker
func (_m *OpenChecker) OpenDefineAssets(ctx types.Context, poolId uint64, msg *perpetualtypes.MsgOpen, baseCurrency string, isBroker bool) (*perpetualtypes.MTP, error) {
	ret := _m.Called(ctx, poolId, msg, baseCurrency, isBroker)

	if len(ret) == 0 {
		panic("no return value specified for OpenDefineAssets")
	}

	var r0 *perpetualtypes.MTP
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64, *perpetualtypes.MsgOpen, string, bool) (*perpetualtypes.MTP, error)); ok {
		return rf(ctx, poolId, msg, baseCurrency, isBroker)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64, *perpetualtypes.MsgOpen, string, bool) *perpetualtypes.MTP); ok {
		r0 = rf(ctx, poolId, msg, baseCurrency, isBroker)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*perpetualtypes.MTP)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64, *perpetualtypes.MsgOpen, string, bool) error); ok {
		r1 = rf(ctx, poolId, msg, baseCurrency, isBroker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenChecker_OpenDefineAssets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OpenDefineAssets'
type OpenChecker_OpenDefineAssets_Call struct {
	*mock.Call
}

// OpenDefineAssets is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
//   - msg *perpetualtypes.MsgOpen
//   - baseCurrency string
//   - isBroker bool
func (_e *OpenChecker_Expecter) OpenDefineAssets(ctx interface{}, poolId interface{}, msg interface{}, baseCurrency interface{}, isBroker interface{}) *OpenChecker_OpenDefineAssets_Call {
	return &OpenChecker_OpenDefineAssets_Call{Call: _e.mock.On("OpenDefineAssets", ctx, poolId, msg, baseCurrency, isBroker)}
}

func (_c *OpenChecker_OpenDefineAssets_Call) Run(run func(ctx types.Context, poolId uint64, msg *perpetualtypes.MsgOpen, baseCurrency string, isBroker bool)) *OpenChecker_OpenDefineAssets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64), args[2].(*perpetualtypes.MsgOpen), args[3].(string), args[4].(bool))
	})
	return _c
}

func (_c *OpenChecker_OpenDefineAssets_Call) Return(_a0 *perpetualtypes.MTP, _a1 error) *OpenChecker_OpenDefineAssets_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OpenChecker_OpenDefineAssets_Call) RunAndReturn(run func(types.Context, uint64, *perpetualtypes.MsgOpen, string, bool) (*perpetualtypes.MTP, error)) *OpenChecker_OpenDefineAssets_Call {
	_c.Call.Return(run)
	return _c
}

// SetMTP provides a mock function with given fields: ctx, mtp
func (_m *OpenChecker) SetMTP(ctx types.Context, mtp *perpetualtypes.MTP) error {
	ret := _m.Called(ctx, mtp)

	if len(ret) == 0 {
		panic("no return value specified for SetMTP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualtypes.MTP) error); ok {
		r0 = rf(ctx, mtp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_SetMTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMTP'
type OpenChecker_SetMTP_Call struct {
	*mock.Call
}

// SetMTP is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp *perpetualtypes.MTP
func (_e *OpenChecker_Expecter) SetMTP(ctx interface{}, mtp interface{}) *OpenChecker_SetMTP_Call {
	return &OpenChecker_SetMTP_Call{Call: _e.mock.On("SetMTP", ctx, mtp)}
}

func (_c *OpenChecker_SetMTP_Call) Run(run func(ctx types.Context, mtp *perpetualtypes.MTP)) *OpenChecker_SetMTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*perpetualtypes.MTP))
	})
	return _c
}

func (_c *OpenChecker_SetMTP_Call) Return(_a0 error) *OpenChecker_SetMTP_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_SetMTP_Call) RunAndReturn(run func(types.Context, *perpetualtypes.MTP) error) *OpenChecker_SetMTP_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOpenPrice provides a mock function with given fields: ctx, mtp, ammPool, baseCurrency
func (_m *OpenChecker) UpdateOpenPrice(ctx types.Context, mtp *perpetualtypes.MTP, ammPool ammtypes.Pool, baseCurrency string) error {
	ret := _m.Called(ctx, mtp, ammPool, baseCurrency)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOpenPrice")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *perpetualtypes.MTP, ammtypes.Pool, string) error); ok {
		r0 = rf(ctx, mtp, ammPool, baseCurrency)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_UpdateOpenPrice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOpenPrice'
type OpenChecker_UpdateOpenPrice_Call struct {
	*mock.Call
}

// UpdateOpenPrice is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp *perpetualtypes.MTP
//   - ammPool ammtypes.Pool
//   - baseCurrency string
func (_e *OpenChecker_Expecter) UpdateOpenPrice(ctx interface{}, mtp interface{}, ammPool interface{}, baseCurrency interface{}) *OpenChecker_UpdateOpenPrice_Call {
	return &OpenChecker_UpdateOpenPrice_Call{Call: _e.mock.On("UpdateOpenPrice", ctx, mtp, ammPool, baseCurrency)}
}

func (_c *OpenChecker_UpdateOpenPrice_Call) Run(run func(ctx types.Context, mtp *perpetualtypes.MTP, ammPool ammtypes.Pool, baseCurrency string)) *OpenChecker_UpdateOpenPrice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*perpetualtypes.MTP), args[2].(ammtypes.Pool), args[3].(string))
	})
	return _c
}

func (_c *OpenChecker_UpdateOpenPrice_Call) Return(_a0 error) *OpenChecker_UpdateOpenPrice_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_UpdateOpenPrice_Call) RunAndReturn(run func(types.Context, *perpetualtypes.MTP, ammtypes.Pool, string) error) *OpenChecker_UpdateOpenPrice_Call {
	_c.Call.Return(run)
	return _c
}

// NewOpenChecker creates a new instance of OpenChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOpenChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *OpenChecker {
	mock := &OpenChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
