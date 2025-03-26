// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is an autogenerated mock type for the BankKeeper type
type BankKeeper struct {
	mock.Mock
}

type BankKeeper_Expecter struct {
	mock *mock.Mock
}

func (_m *BankKeeper) EXPECT() *BankKeeper_Expecter {
	return &BankKeeper_Expecter{mock: &_m.Mock}
}

// BurnCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *BankKeeper) BurnCoins(ctx context.Context, moduleName string, amt types.Coins) error {
	ret := _m.Called(ctx, moduleName, amt)

	if len(ret) == 0 {
		panic("no return value specified for BurnCoins")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.Coins) error); ok {
		r0 = rf(ctx, moduleName, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BankKeeper_BurnCoins_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BurnCoins'
type BankKeeper_BurnCoins_Call struct {
	*mock.Call
}

// BurnCoins is a helper method to define mock.On call
//   - ctx context.Context
//   - moduleName string
//   - amt types.Coins
func (_e *BankKeeper_Expecter) BurnCoins(ctx interface{}, moduleName interface{}, amt interface{}) *BankKeeper_BurnCoins_Call {
	return &BankKeeper_BurnCoins_Call{Call: _e.mock.On("BurnCoins", ctx, moduleName, amt)}
}

func (_c *BankKeeper_BurnCoins_Call) Run(run func(ctx context.Context, moduleName string, amt types.Coins)) *BankKeeper_BurnCoins_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(types.Coins))
	})
	return _c
}

func (_c *BankKeeper_BurnCoins_Call) Return(_a0 error) *BankKeeper_BurnCoins_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_BurnCoins_Call) RunAndReturn(run func(context.Context, string, types.Coins) error) *BankKeeper_BurnCoins_Call {
	_c.Call.Return(run)
	return _c
}

// GetBalance provides a mock function with given fields: ctx, addr, denom
func (_m *BankKeeper) GetBalance(ctx context.Context, addr types.AccAddress, denom string) types.Coin {
	ret := _m.Called(ctx, addr, denom)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 types.Coin
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress, string) types.Coin); ok {
		r0 = rf(ctx, addr, denom)
	} else {
		r0 = ret.Get(0).(types.Coin)
	}

	return r0
}

// BankKeeper_GetBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBalance'
type BankKeeper_GetBalance_Call struct {
	*mock.Call
}

// GetBalance is a helper method to define mock.On call
//   - ctx context.Context
//   - addr types.AccAddress
//   - denom string
func (_e *BankKeeper_Expecter) GetBalance(ctx interface{}, addr interface{}, denom interface{}) *BankKeeper_GetBalance_Call {
	return &BankKeeper_GetBalance_Call{Call: _e.mock.On("GetBalance", ctx, addr, denom)}
}

func (_c *BankKeeper_GetBalance_Call) Run(run func(ctx context.Context, addr types.AccAddress, denom string)) *BankKeeper_GetBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.AccAddress), args[2].(string))
	})
	return _c
}

func (_c *BankKeeper_GetBalance_Call) Return(_a0 types.Coin) *BankKeeper_GetBalance_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_GetBalance_Call) RunAndReturn(run func(context.Context, types.AccAddress, string) types.Coin) *BankKeeper_GetBalance_Call {
	_c.Call.Return(run)
	return _c
}

// IterateAllDenomMetaData provides a mock function with given fields: ctx, cb
func (_m *BankKeeper) IterateAllDenomMetaData(ctx context.Context, cb func(banktypes.Metadata) bool) {
	_m.Called(ctx, cb)
}

// BankKeeper_IterateAllDenomMetaData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IterateAllDenomMetaData'
type BankKeeper_IterateAllDenomMetaData_Call struct {
	*mock.Call
}

// IterateAllDenomMetaData is a helper method to define mock.On call
//   - ctx context.Context
//   - cb func(banktypes.Metadata) bool
func (_e *BankKeeper_Expecter) IterateAllDenomMetaData(ctx interface{}, cb interface{}) *BankKeeper_IterateAllDenomMetaData_Call {
	return &BankKeeper_IterateAllDenomMetaData_Call{Call: _e.mock.On("IterateAllDenomMetaData", ctx, cb)}
}

func (_c *BankKeeper_IterateAllDenomMetaData_Call) Run(run func(ctx context.Context, cb func(banktypes.Metadata) bool)) *BankKeeper_IterateAllDenomMetaData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(banktypes.Metadata) bool))
	})
	return _c
}

func (_c *BankKeeper_IterateAllDenomMetaData_Call) Return() *BankKeeper_IterateAllDenomMetaData_Call {
	_c.Call.Return()
	return _c
}

func (_c *BankKeeper_IterateAllDenomMetaData_Call) RunAndReturn(run func(context.Context, func(banktypes.Metadata) bool)) *BankKeeper_IterateAllDenomMetaData_Call {
	_c.Run(run)
	return _c
}

// SendCoinsFromAccountToModule provides a mock function with given fields: ctx, senderAddr, recipientModule, amt
func (_m *BankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
	ret := _m.Called(ctx, senderAddr, recipientModule, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromAccountToModule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress, string, types.Coins) error); ok {
		r0 = rf(ctx, senderAddr, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BankKeeper_SendCoinsFromAccountToModule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendCoinsFromAccountToModule'
type BankKeeper_SendCoinsFromAccountToModule_Call struct {
	*mock.Call
}

// SendCoinsFromAccountToModule is a helper method to define mock.On call
//   - ctx context.Context
//   - senderAddr types.AccAddress
//   - recipientModule string
//   - amt types.Coins
func (_e *BankKeeper_Expecter) SendCoinsFromAccountToModule(ctx interface{}, senderAddr interface{}, recipientModule interface{}, amt interface{}) *BankKeeper_SendCoinsFromAccountToModule_Call {
	return &BankKeeper_SendCoinsFromAccountToModule_Call{Call: _e.mock.On("SendCoinsFromAccountToModule", ctx, senderAddr, recipientModule, amt)}
}

func (_c *BankKeeper_SendCoinsFromAccountToModule_Call) Run(run func(ctx context.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins)) *BankKeeper_SendCoinsFromAccountToModule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.AccAddress), args[2].(string), args[3].(types.Coins))
	})
	return _c
}

func (_c *BankKeeper_SendCoinsFromAccountToModule_Call) Return(_a0 error) *BankKeeper_SendCoinsFromAccountToModule_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_SendCoinsFromAccountToModule_Call) RunAndReturn(run func(context.Context, types.AccAddress, string, types.Coins) error) *BankKeeper_SendCoinsFromAccountToModule_Call {
	_c.Call.Return(run)
	return _c
}

// SpendableCoins provides a mock function with given fields: ctx, addr
func (_m *BankKeeper) SpendableCoins(ctx context.Context, addr types.AccAddress) types.Coins {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for SpendableCoins")
	}

	var r0 types.Coins
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress) types.Coins); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Coins)
		}
	}

	return r0
}

// BankKeeper_SpendableCoins_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SpendableCoins'
type BankKeeper_SpendableCoins_Call struct {
	*mock.Call
}

// SpendableCoins is a helper method to define mock.On call
//   - ctx context.Context
//   - addr types.AccAddress
func (_e *BankKeeper_Expecter) SpendableCoins(ctx interface{}, addr interface{}) *BankKeeper_SpendableCoins_Call {
	return &BankKeeper_SpendableCoins_Call{Call: _e.mock.On("SpendableCoins", ctx, addr)}
}

func (_c *BankKeeper_SpendableCoins_Call) Run(run func(ctx context.Context, addr types.AccAddress)) *BankKeeper_SpendableCoins_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.AccAddress))
	})
	return _c
}

func (_c *BankKeeper_SpendableCoins_Call) Return(_a0 types.Coins) *BankKeeper_SpendableCoins_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_SpendableCoins_Call) RunAndReturn(run func(context.Context, types.AccAddress) types.Coins) *BankKeeper_SpendableCoins_Call {
	_c.Call.Return(run)
	return _c
}

// NewBankKeeper creates a new instance of BankKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBankKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *BankKeeper {
	mock := &BankKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
