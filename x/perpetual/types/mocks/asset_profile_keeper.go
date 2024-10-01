// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import (
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// AssetProfileKeeper is an autogenerated mock type for the AssetProfileKeeper type
type AssetProfileKeeper struct {
	mock.Mock
}

type AssetProfileKeeper_Expecter struct {
	mock *mock.Mock
}

func (_m *AssetProfileKeeper) EXPECT() *AssetProfileKeeper_Expecter {
	return &AssetProfileKeeper_Expecter{mock: &_m.Mock}
}

// GetEntry provides a mock function with given fields: ctx, baseDenom
func (_m *AssetProfileKeeper) GetEntry(ctx types.Context, baseDenom string) (assetprofiletypes.Entry, bool) {
	ret := _m.Called(ctx, baseDenom)

	if len(ret) == 0 {
		panic("no return value specified for GetEntry")
	}

	var r0 assetprofiletypes.Entry
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (assetprofiletypes.Entry, bool)); ok {
		return rf(ctx, baseDenom)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) assetprofiletypes.Entry); ok {
		r0 = rf(ctx, baseDenom)
	} else {
		r0 = ret.Get(0).(assetprofiletypes.Entry)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, baseDenom)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// AssetProfileKeeper_GetEntry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEntry'
type AssetProfileKeeper_GetEntry_Call struct {
	*mock.Call
}

// GetEntry is a helper method to define mock.On call
//   - ctx types.Context
//   - baseDenom string
func (_e *AssetProfileKeeper_Expecter) GetEntry(ctx interface{}, baseDenom interface{}) *AssetProfileKeeper_GetEntry_Call {
	return &AssetProfileKeeper_GetEntry_Call{Call: _e.mock.On("GetEntry", ctx, baseDenom)}
}

func (_c *AssetProfileKeeper_GetEntry_Call) Run(run func(ctx types.Context, baseDenom string)) *AssetProfileKeeper_GetEntry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *AssetProfileKeeper_GetEntry_Call) Return(val assetprofiletypes.Entry, found bool) *AssetProfileKeeper_GetEntry_Call {
	_c.Call.Return(val, found)
	return _c
}

func (_c *AssetProfileKeeper_GetEntry_Call) RunAndReturn(run func(types.Context, string) (assetprofiletypes.Entry, bool)) *AssetProfileKeeper_GetEntry_Call {
	_c.Call.Return(run)
	return _c
}

// GetEntryByDenom provides a mock function with given fields: ctx, denom
func (_m *AssetProfileKeeper) GetEntryByDenom(ctx types.Context, denom string) (assetprofiletypes.Entry, bool) {
	ret := _m.Called(ctx, denom)

	if len(ret) == 0 {
		panic("no return value specified for GetEntryByDenom")
	}

	var r0 assetprofiletypes.Entry
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (assetprofiletypes.Entry, bool)); ok {
		return rf(ctx, denom)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) assetprofiletypes.Entry); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(assetprofiletypes.Entry)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, denom)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// AssetProfileKeeper_GetEntryByDenom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEntryByDenom'
type AssetProfileKeeper_GetEntryByDenom_Call struct {
	*mock.Call
}

// GetEntryByDenom is a helper method to define mock.On call
//   - ctx types.Context
//   - denom string
func (_e *AssetProfileKeeper_Expecter) GetEntryByDenom(ctx interface{}, denom interface{}) *AssetProfileKeeper_GetEntryByDenom_Call {
	return &AssetProfileKeeper_GetEntryByDenom_Call{Call: _e.mock.On("GetEntryByDenom", ctx, denom)}
}

func (_c *AssetProfileKeeper_GetEntryByDenom_Call) Run(run func(ctx types.Context, denom string)) *AssetProfileKeeper_GetEntryByDenom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *AssetProfileKeeper_GetEntryByDenom_Call) Return(val assetprofiletypes.Entry, found bool) *AssetProfileKeeper_GetEntryByDenom_Call {
	_c.Call.Return(val, found)
	return _c
}

func (_c *AssetProfileKeeper_GetEntryByDenom_Call) RunAndReturn(run func(types.Context, string) (assetprofiletypes.Entry, bool)) *AssetProfileKeeper_GetEntryByDenom_Call {
	_c.Call.Return(run)
	return _c
}

// NewAssetProfileKeeper creates a new instance of AssetProfileKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAssetProfileKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *AssetProfileKeeper {
	mock := &AssetProfileKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
