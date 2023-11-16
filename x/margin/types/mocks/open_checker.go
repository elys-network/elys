// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	ammtypes "github.com/elys-network/elys/x/amm/types"
	margintypes "github.com/elys-network/elys/x/margin/types"

	mock "github.com/stretchr/testify/mock"

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

// CheckLongAssets provides a mock function with given fields: ctx, collateralAsset, borrowAsset
func (_m *OpenChecker) CheckLongAssets(ctx types.Context, collateralAsset string, borrowAsset string) error {
	ret := _m.Called(ctx, collateralAsset, borrowAsset)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, string) error); ok {
		r0 = rf(ctx, collateralAsset, borrowAsset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_CheckLongAssets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckLongAssets'
type OpenChecker_CheckLongAssets_Call struct {
	*mock.Call
}

// CheckLongAssets is a helper method to define mock.On call
//   - ctx types.Context
//   - collateralAsset string
//   - borrowAsset string
func (_e *OpenChecker_Expecter) CheckLongAssets(ctx interface{}, collateralAsset interface{}, borrowAsset interface{}) *OpenChecker_CheckLongAssets_Call {
	return &OpenChecker_CheckLongAssets_Call{Call: _e.mock.On("CheckLongAssets", ctx, collateralAsset, borrowAsset)}
}

func (_c *OpenChecker_CheckLongAssets_Call) Run(run func(ctx types.Context, collateralAsset string, borrowAsset string)) *OpenChecker_CheckLongAssets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *OpenChecker_CheckLongAssets_Call) Return(_a0 error) *OpenChecker_CheckLongAssets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckLongAssets_Call) RunAndReturn(run func(types.Context, string, string) error) *OpenChecker_CheckLongAssets_Call {
	_c.Call.Return(run)
	return _c
}

// CheckMaxOpenPositions provides a mock function with given fields: ctx
func (_m *OpenChecker) CheckMaxOpenPositions(ctx types.Context) error {
	ret := _m.Called(ctx)

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

// CheckSamePosition provides a mock function with given fields: ctx, msg
func (_m *OpenChecker) CheckSamePosition(ctx types.Context, msg *margintypes.MsgOpen) *margintypes.MTP {
	ret := _m.Called(ctx, msg)

	var r0 *margintypes.MTP
	if rf, ok := ret.Get(0).(func(types.Context, *margintypes.MsgOpen) *margintypes.MTP); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*margintypes.MTP)
		}
	}

	return r0
}

// OpenChecker_CheckSamePosition_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckSamePosition'
type OpenChecker_CheckSamePosition_Call struct {
	*mock.Call
}

// CheckSamePosition is a helper method to define mock.On call
//   - ctx types.Context
//   - msg *margintypes.MsgOpen
func (_e *OpenChecker_Expecter) CheckSamePosition(ctx interface{}, msg interface{}) *OpenChecker_CheckSamePosition_Call {
	return &OpenChecker_CheckSamePosition_Call{Call: _e.mock.On("CheckSamePosition", ctx, msg)}
}

func (_c *OpenChecker_CheckSamePosition_Call) Run(run func(ctx types.Context, msg *margintypes.MsgOpen)) *OpenChecker_CheckSamePosition_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*margintypes.MsgOpen))
	})
	return _c
}

func (_c *OpenChecker_CheckSamePosition_Call) Return(_a0 *margintypes.MTP) *OpenChecker_CheckSamePosition_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckSamePosition_Call) RunAndReturn(run func(types.Context, *margintypes.MsgOpen) *margintypes.MTP) *OpenChecker_CheckSamePosition_Call {
	_c.Call.Return(run)
	return _c
}

// CheckShortAssets provides a mock function with given fields: ctx, collateralAsset, borrowAsset
func (_m *OpenChecker) CheckShortAssets(ctx types.Context, collateralAsset string, borrowAsset string) error {
	ret := _m.Called(ctx, collateralAsset, borrowAsset)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, string) error); ok {
		r0 = rf(ctx, collateralAsset, borrowAsset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OpenChecker_CheckShortAssets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckShortAssets'
type OpenChecker_CheckShortAssets_Call struct {
	*mock.Call
}

// CheckShortAssets is a helper method to define mock.On call
//   - ctx types.Context
//   - collateralAsset string
//   - borrowAsset string
func (_e *OpenChecker_Expecter) CheckShortAssets(ctx interface{}, collateralAsset interface{}, borrowAsset interface{}) *OpenChecker_CheckShortAssets_Call {
	return &OpenChecker_CheckShortAssets_Call{Call: _e.mock.On("CheckShortAssets", ctx, collateralAsset, borrowAsset)}
}

func (_c *OpenChecker_CheckShortAssets_Call) Run(run func(ctx types.Context, collateralAsset string, borrowAsset string)) *OpenChecker_CheckShortAssets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *OpenChecker_CheckShortAssets_Call) Return(_a0 error) *OpenChecker_CheckShortAssets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckShortAssets_Call) RunAndReturn(run func(types.Context, string, string) error) *OpenChecker_CheckShortAssets_Call {
	_c.Call.Return(run)
	return _c
}

// CheckUserAuthorization provides a mock function with given fields: ctx, msg
func (_m *OpenChecker) CheckUserAuthorization(ctx types.Context, msg *margintypes.MsgOpen) error {
	ret := _m.Called(ctx, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *margintypes.MsgOpen) error); ok {
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
//   - msg *margintypes.MsgOpen
func (_e *OpenChecker_Expecter) CheckUserAuthorization(ctx interface{}, msg interface{}) *OpenChecker_CheckUserAuthorization_Call {
	return &OpenChecker_CheckUserAuthorization_Call{Call: _e.mock.On("CheckUserAuthorization", ctx, msg)}
}

func (_c *OpenChecker_CheckUserAuthorization_Call) Run(run func(ctx types.Context, msg *margintypes.MsgOpen)) *OpenChecker_CheckUserAuthorization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*margintypes.MsgOpen))
	})
	return _c
}

func (_c *OpenChecker_CheckUserAuthorization_Call) Return(_a0 error) *OpenChecker_CheckUserAuthorization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_CheckUserAuthorization_Call) RunAndReturn(run func(types.Context, *margintypes.MsgOpen) error) *OpenChecker_CheckUserAuthorization_Call {
	_c.Call.Return(run)
	return _c
}

// EmitOpenEvent provides a mock function with given fields: ctx, mtp
func (_m *OpenChecker) EmitOpenEvent(ctx types.Context, mtp *margintypes.MTP) {
	_m.Called(ctx, mtp)
}

// OpenChecker_EmitOpenEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EmitOpenEvent'
type OpenChecker_EmitOpenEvent_Call struct {
	*mock.Call
}

// EmitOpenEvent is a helper method to define mock.On call
//   - ctx types.Context
//   - mtp *margintypes.MTP
func (_e *OpenChecker_Expecter) EmitOpenEvent(ctx interface{}, mtp interface{}) *OpenChecker_EmitOpenEvent_Call {
	return &OpenChecker_EmitOpenEvent_Call{Call: _e.mock.On("EmitOpenEvent", ctx, mtp)}
}

func (_c *OpenChecker_EmitOpenEvent_Call) Run(run func(ctx types.Context, mtp *margintypes.MTP)) *OpenChecker_EmitOpenEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*margintypes.MTP))
	})
	return _c
}

func (_c *OpenChecker_EmitOpenEvent_Call) Return() *OpenChecker_EmitOpenEvent_Call {
	_c.Call.Return()
	return _c
}

func (_c *OpenChecker_EmitOpenEvent_Call) RunAndReturn(run func(types.Context, *margintypes.MTP)) *OpenChecker_EmitOpenEvent_Call {
	_c.Call.Return(run)
	return _c
}

// GetMaxOpenPositions provides a mock function with given fields: ctx
func (_m *OpenChecker) GetMaxOpenPositions(ctx types.Context) uint64 {
	ret := _m.Called(ctx)

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

// GetTradingAsset provides a mock function with given fields: collateralAsset, borrowAsset, baseCurrency
func (_m *OpenChecker) GetTradingAsset(collateralAsset string, borrowAsset string, baseCurrency string) string {
	ret := _m.Called(collateralAsset, borrowAsset, baseCurrency)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(collateralAsset, borrowAsset, baseCurrency)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OpenChecker_GetTradingAsset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTradingAsset'
type OpenChecker_GetTradingAsset_Call struct {
	*mock.Call
}

// GetTradingAsset is a helper method to define mock.On call
//   - collateralAsset string
//   - borrowAsset string
//   - baseCurrency string
func (_e *OpenChecker_Expecter) GetTradingAsset(collateralAsset interface{}, borrowAsset interface{}, baseCurrency interface{}) *OpenChecker_GetTradingAsset_Call {
	return &OpenChecker_GetTradingAsset_Call{Call: _e.mock.On("GetTradingAsset", collateralAsset, borrowAsset, baseCurrency)}
}

func (_c *OpenChecker_GetTradingAsset_Call) Run(run func(collateralAsset string, borrowAsset string, baseCurrency string)) *OpenChecker_GetTradingAsset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *OpenChecker_GetTradingAsset_Call) Return(_a0 string) *OpenChecker_GetTradingAsset_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_GetTradingAsset_Call) RunAndReturn(run func(string, string, string) string) *OpenChecker_GetTradingAsset_Call {
	_c.Call.Return(run)
	return _c
}

// OpenLong provides a mock function with given fields: ctx, poolId, msg
func (_m *OpenChecker) OpenLong(ctx types.Context, poolId uint64, msg *margintypes.MsgOpen) (*margintypes.MTP, error) {
	ret := _m.Called(ctx, poolId, msg)

	var r0 *margintypes.MTP
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64, *margintypes.MsgOpen) (*margintypes.MTP, error)); ok {
		return rf(ctx, poolId, msg)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64, *margintypes.MsgOpen) *margintypes.MTP); ok {
		r0 = rf(ctx, poolId, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*margintypes.MTP)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64, *margintypes.MsgOpen) error); ok {
		r1 = rf(ctx, poolId, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenChecker_OpenLong_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OpenLong'
type OpenChecker_OpenLong_Call struct {
	*mock.Call
}

// OpenLong is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
//   - msg *margintypes.MsgOpen
func (_e *OpenChecker_Expecter) OpenLong(ctx interface{}, poolId interface{}, msg interface{}) *OpenChecker_OpenLong_Call {
	return &OpenChecker_OpenLong_Call{Call: _e.mock.On("OpenLong", ctx, poolId, msg)}
}

func (_c *OpenChecker_OpenLong_Call) Run(run func(ctx types.Context, poolId uint64, msg *margintypes.MsgOpen)) *OpenChecker_OpenLong_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64), args[2].(*margintypes.MsgOpen))
	})
	return _c
}

func (_c *OpenChecker_OpenLong_Call) Return(_a0 *margintypes.MTP, _a1 error) *OpenChecker_OpenLong_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OpenChecker_OpenLong_Call) RunAndReturn(run func(types.Context, uint64, *margintypes.MsgOpen) (*margintypes.MTP, error)) *OpenChecker_OpenLong_Call {
	_c.Call.Return(run)
	return _c
}

// OpenShort provides a mock function with given fields: ctx, poolId, msg
func (_m *OpenChecker) OpenShort(ctx types.Context, poolId uint64, msg *margintypes.MsgOpen) (*margintypes.MTP, error) {
	ret := _m.Called(ctx, poolId, msg)

	var r0 *margintypes.MTP
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64, *margintypes.MsgOpen) (*margintypes.MTP, error)); ok {
		return rf(ctx, poolId, msg)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64, *margintypes.MsgOpen) *margintypes.MTP); ok {
		r0 = rf(ctx, poolId, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*margintypes.MTP)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64, *margintypes.MsgOpen) error); ok {
		r1 = rf(ctx, poolId, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenChecker_OpenShort_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OpenShort'
type OpenChecker_OpenShort_Call struct {
	*mock.Call
}

// OpenShort is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
//   - msg *margintypes.MsgOpen
func (_e *OpenChecker_Expecter) OpenShort(ctx interface{}, poolId interface{}, msg interface{}) *OpenChecker_OpenShort_Call {
	return &OpenChecker_OpenShort_Call{Call: _e.mock.On("OpenShort", ctx, poolId, msg)}
}

func (_c *OpenChecker_OpenShort_Call) Run(run func(ctx types.Context, poolId uint64, msg *margintypes.MsgOpen)) *OpenChecker_OpenShort_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64), args[2].(*margintypes.MsgOpen))
	})
	return _c
}

func (_c *OpenChecker_OpenShort_Call) Return(_a0 *margintypes.MTP, _a1 error) *OpenChecker_OpenShort_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OpenChecker_OpenShort_Call) RunAndReturn(run func(types.Context, uint64, *margintypes.MsgOpen) (*margintypes.MTP, error)) *OpenChecker_OpenShort_Call {
	_c.Call.Return(run)
	return _c
}

// PreparePools provides a mock function with given fields: ctx, tradingAsset
func (_m *OpenChecker) PreparePools(ctx types.Context, tradingAsset string) (uint64, ammtypes.Pool, margintypes.Pool, error) {
	ret := _m.Called(ctx, tradingAsset)

	var r0 uint64
	var r1 ammtypes.Pool
	var r2 margintypes.Pool
	var r3 error
	if rf, ok := ret.Get(0).(func(types.Context, string) (uint64, ammtypes.Pool, margintypes.Pool, error)); ok {
		return rf(ctx, tradingAsset)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) uint64); ok {
		r0 = rf(ctx, tradingAsset)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) ammtypes.Pool); ok {
		r1 = rf(ctx, tradingAsset)
	} else {
		r1 = ret.Get(1).(ammtypes.Pool)
	}

	if rf, ok := ret.Get(2).(func(types.Context, string) margintypes.Pool); ok {
		r2 = rf(ctx, tradingAsset)
	} else {
		r2 = ret.Get(2).(margintypes.Pool)
	}

	if rf, ok := ret.Get(3).(func(types.Context, string) error); ok {
		r3 = rf(ctx, tradingAsset)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// OpenChecker_PreparePools_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PreparePools'
type OpenChecker_PreparePools_Call struct {
	*mock.Call
}

// PreparePools is a helper method to define mock.On call
//   - ctx types.Context
//   - tradingAsset string
func (_e *OpenChecker_Expecter) PreparePools(ctx interface{}, tradingAsset interface{}) *OpenChecker_PreparePools_Call {
	return &OpenChecker_PreparePools_Call{Call: _e.mock.On("PreparePools", ctx, tradingAsset)}
}

func (_c *OpenChecker_PreparePools_Call) Run(run func(ctx types.Context, tradingAsset string)) *OpenChecker_PreparePools_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *OpenChecker_PreparePools_Call) Return(poolId uint64, ammPool ammtypes.Pool, pool margintypes.Pool, err error) *OpenChecker_PreparePools_Call {
	_c.Call.Return(poolId, ammPool, pool, err)
	return _c
}

func (_c *OpenChecker_PreparePools_Call) RunAndReturn(run func(types.Context, string) (uint64, ammtypes.Pool, margintypes.Pool, error)) *OpenChecker_PreparePools_Call {
	_c.Call.Return(run)
	return _c
}

// SetMTP provides a mock function with given fields: ctx, mtp
func (_m *OpenChecker) SetMTP(ctx types.Context, mtp *margintypes.MTP) error {
	ret := _m.Called(ctx, mtp)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *margintypes.MTP) error); ok {
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
//   - mtp *margintypes.MTP
func (_e *OpenChecker_Expecter) SetMTP(ctx interface{}, mtp interface{}) *OpenChecker_SetMTP_Call {
	return &OpenChecker_SetMTP_Call{Call: _e.mock.On("SetMTP", ctx, mtp)}
}

func (_c *OpenChecker_SetMTP_Call) Run(run func(ctx types.Context, mtp *margintypes.MTP)) *OpenChecker_SetMTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(*margintypes.MTP))
	})
	return _c
}

func (_c *OpenChecker_SetMTP_Call) Return(_a0 error) *OpenChecker_SetMTP_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OpenChecker_SetMTP_Call) RunAndReturn(run func(types.Context, *margintypes.MTP) error) *OpenChecker_SetMTP_Call {
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
