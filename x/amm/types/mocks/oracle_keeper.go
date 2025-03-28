// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	elystypes "github.com/elys-network/elys/types"
	mock "github.com/stretchr/testify/mock"

	oracletypes "github.com/elys-network/elys/x/oracle/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// OracleKeeper is an autogenerated mock type for the OracleKeeper type
type OracleKeeper struct {
	mock.Mock
}

type OracleKeeper_Expecter struct {
	mock *mock.Mock
}

func (_m *OracleKeeper) EXPECT() *OracleKeeper_Expecter {
	return &OracleKeeper_Expecter{mock: &_m.Mock}
}

// GetAssetInfo provides a mock function with given fields: ctx, denom
func (_m *OracleKeeper) GetAssetInfo(ctx types.Context, denom string) (oracletypes.AssetInfo, bool) {
	ret := _m.Called(ctx, denom)

	if len(ret) == 0 {
		panic("no return value specified for GetAssetInfo")
	}

	var r0 oracletypes.AssetInfo
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (oracletypes.AssetInfo, bool)); ok {
		return rf(ctx, denom)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) oracletypes.AssetInfo); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(oracletypes.AssetInfo)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, denom)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// OracleKeeper_GetAssetInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAssetInfo'
type OracleKeeper_GetAssetInfo_Call struct {
	*mock.Call
}

// GetAssetInfo is a helper method to define mock.On call
//   - ctx types.Context
//   - denom string
func (_e *OracleKeeper_Expecter) GetAssetInfo(ctx interface{}, denom interface{}) *OracleKeeper_GetAssetInfo_Call {
	return &OracleKeeper_GetAssetInfo_Call{Call: _e.mock.On("GetAssetInfo", ctx, denom)}
}

func (_c *OracleKeeper_GetAssetInfo_Call) Run(run func(ctx types.Context, denom string)) *OracleKeeper_GetAssetInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *OracleKeeper_GetAssetInfo_Call) Return(val oracletypes.AssetInfo, found bool) *OracleKeeper_GetAssetInfo_Call {
	_c.Call.Return(val, found)
	return _c
}

func (_c *OracleKeeper_GetAssetInfo_Call) RunAndReturn(run func(types.Context, string) (oracletypes.AssetInfo, bool)) *OracleKeeper_GetAssetInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetAssetPrice provides a mock function with given fields: ctx, asset
func (_m *OracleKeeper) GetAssetPrice(ctx types.Context, asset string) (oracletypes.Price, bool) {
	ret := _m.Called(ctx, asset)

	if len(ret) == 0 {
		panic("no return value specified for GetAssetPrice")
	}

	var r0 oracletypes.Price
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (oracletypes.Price, bool)); ok {
		return rf(ctx, asset)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) oracletypes.Price); ok {
		r0 = rf(ctx, asset)
	} else {
		r0 = ret.Get(0).(oracletypes.Price)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, asset)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// OracleKeeper_GetAssetPrice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAssetPrice'
type OracleKeeper_GetAssetPrice_Call struct {
	*mock.Call
}

// GetAssetPrice is a helper method to define mock.On call
//   - ctx types.Context
//   - asset string
func (_e *OracleKeeper_Expecter) GetAssetPrice(ctx interface{}, asset interface{}) *OracleKeeper_GetAssetPrice_Call {
	return &OracleKeeper_GetAssetPrice_Call{Call: _e.mock.On("GetAssetPrice", ctx, asset)}
}

func (_c *OracleKeeper_GetAssetPrice_Call) Run(run func(ctx types.Context, asset string)) *OracleKeeper_GetAssetPrice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *OracleKeeper_GetAssetPrice_Call) Return(_a0 oracletypes.Price, _a1 bool) *OracleKeeper_GetAssetPrice_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OracleKeeper_GetAssetPrice_Call) RunAndReturn(run func(types.Context, string) (oracletypes.Price, bool)) *OracleKeeper_GetAssetPrice_Call {
	_c.Call.Return(run)
	return _c
}

// GetAssetPriceFromDenom provides a mock function with given fields: ctx, denom
func (_m *OracleKeeper) GetAssetPriceFromDenom(ctx types.Context, denom string) (elystypes.Dec34, uint64) {
	ret := _m.Called(ctx, denom)

	if len(ret) == 0 {
		panic("no return value specified for GetAssetPriceFromDenom")
	}

	var r0 elystypes.Dec34
	var r1 uint64
	if rf, ok := ret.Get(0).(func(types.Context, string) (elystypes.Dec34, uint64)); ok {
		return rf(ctx, denom)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) elystypes.Dec34); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(elystypes.Dec34)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) uint64); ok {
		r1 = rf(ctx, denom)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	return r0, r1
}

// OracleKeeper_GetAssetPriceFromDenom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAssetPriceFromDenom'
type OracleKeeper_GetAssetPriceFromDenom_Call struct {
	*mock.Call
}

// GetAssetPriceFromDenom is a helper method to define mock.On call
//   - ctx types.Context
//   - denom string
func (_e *OracleKeeper_Expecter) GetAssetPriceFromDenom(ctx interface{}, denom interface{}) *OracleKeeper_GetAssetPriceFromDenom_Call {
	return &OracleKeeper_GetAssetPriceFromDenom_Call{Call: _e.mock.On("GetAssetPriceFromDenom", ctx, denom)}
}

func (_c *OracleKeeper_GetAssetPriceFromDenom_Call) Run(run func(ctx types.Context, denom string)) *OracleKeeper_GetAssetPriceFromDenom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *OracleKeeper_GetAssetPriceFromDenom_Call) Return(_a0 elystypes.Dec34, _a1 uint64) *OracleKeeper_GetAssetPriceFromDenom_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OracleKeeper_GetAssetPriceFromDenom_Call) RunAndReturn(run func(types.Context, string) (elystypes.Dec34, uint64)) *OracleKeeper_GetAssetPriceFromDenom_Call {
	_c.Call.Return(run)
	return _c
}

// GetPriceFeeder provides a mock function with given fields: ctx, feeder
func (_m *OracleKeeper) GetPriceFeeder(ctx types.Context, feeder types.AccAddress) (oracletypes.PriceFeeder, bool) {
	ret := _m.Called(ctx, feeder)

	if len(ret) == 0 {
		panic("no return value specified for GetPriceFeeder")
	}

	var r0 oracletypes.PriceFeeder
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) (oracletypes.PriceFeeder, bool)); ok {
		return rf(ctx, feeder)
	}
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) oracletypes.PriceFeeder); ok {
		r0 = rf(ctx, feeder)
	} else {
		r0 = ret.Get(0).(oracletypes.PriceFeeder)
	}

	if rf, ok := ret.Get(1).(func(types.Context, types.AccAddress) bool); ok {
		r1 = rf(ctx, feeder)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// OracleKeeper_GetPriceFeeder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPriceFeeder'
type OracleKeeper_GetPriceFeeder_Call struct {
	*mock.Call
}

// GetPriceFeeder is a helper method to define mock.On call
//   - ctx types.Context
//   - feeder types.AccAddress
func (_e *OracleKeeper_Expecter) GetPriceFeeder(ctx interface{}, feeder interface{}) *OracleKeeper_GetPriceFeeder_Call {
	return &OracleKeeper_GetPriceFeeder_Call{Call: _e.mock.On("GetPriceFeeder", ctx, feeder)}
}

func (_c *OracleKeeper_GetPriceFeeder_Call) Run(run func(ctx types.Context, feeder types.AccAddress)) *OracleKeeper_GetPriceFeeder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(types.AccAddress))
	})
	return _c
}

func (_c *OracleKeeper_GetPriceFeeder_Call) Return(val oracletypes.PriceFeeder, found bool) *OracleKeeper_GetPriceFeeder_Call {
	_c.Call.Return(val, found)
	return _c
}

func (_c *OracleKeeper_GetPriceFeeder_Call) RunAndReturn(run func(types.Context, types.AccAddress) (oracletypes.PriceFeeder, bool)) *OracleKeeper_GetPriceFeeder_Call {
	_c.Call.Return(run)
	return _c
}

// NewOracleKeeper creates a new instance of OracleKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOracleKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *OracleKeeper {
	mock := &OracleKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
