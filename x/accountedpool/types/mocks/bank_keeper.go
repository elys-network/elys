// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"
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

// BlockedAddr provides a mock function with given fields: addr
func (_m *BankKeeper) BlockedAddr(addr types.AccAddress) bool {
	ret := _m.Called(addr)

	if len(ret) == 0 {
		panic("no return value specified for BlockedAddr")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.AccAddress) bool); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// BankKeeper_BlockedAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockedAddr'
type BankKeeper_BlockedAddr_Call struct {
	*mock.Call
}

// BlockedAddr is a helper method to define mock.On call
//   - addr types.AccAddress
func (_e *BankKeeper_Expecter) BlockedAddr(addr interface{}) *BankKeeper_BlockedAddr_Call {
	return &BankKeeper_BlockedAddr_Call{Call: _e.mock.On("BlockedAddr", addr)}
}

func (_c *BankKeeper_BlockedAddr_Call) Run(run func(addr types.AccAddress)) *BankKeeper_BlockedAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.AccAddress))
	})
	return _c
}

func (_c *BankKeeper_BlockedAddr_Call) Return(_a0 bool) *BankKeeper_BlockedAddr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_BlockedAddr_Call) RunAndReturn(run func(types.AccAddress) bool) *BankKeeper_BlockedAddr_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllBalances provides a mock function with given fields: ctx, addr
func (_m *BankKeeper) GetAllBalances(ctx context.Context, addr types.AccAddress) types.Coins {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for GetAllBalances")
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

// BankKeeper_GetAllBalances_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllBalances'
type BankKeeper_GetAllBalances_Call struct {
	*mock.Call
}

// GetAllBalances is a helper method to define mock.On call
//   - ctx context.Context
//   - addr types.AccAddress
func (_e *BankKeeper_Expecter) GetAllBalances(ctx interface{}, addr interface{}) *BankKeeper_GetAllBalances_Call {
	return &BankKeeper_GetAllBalances_Call{Call: _e.mock.On("GetAllBalances", ctx, addr)}
}

func (_c *BankKeeper_GetAllBalances_Call) Run(run func(ctx context.Context, addr types.AccAddress)) *BankKeeper_GetAllBalances_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.AccAddress))
	})
	return _c
}

func (_c *BankKeeper_GetAllBalances_Call) Return(_a0 types.Coins) *BankKeeper_GetAllBalances_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_GetAllBalances_Call) RunAndReturn(run func(context.Context, types.AccAddress) types.Coins) *BankKeeper_GetAllBalances_Call {
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

// HasBalance provides a mock function with given fields: ctx, addr, amt
func (_m *BankKeeper) HasBalance(ctx context.Context, addr types.AccAddress, amt types.Coin) bool {
	ret := _m.Called(ctx, addr, amt)

	if len(ret) == 0 {
		panic("no return value specified for HasBalance")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress, types.Coin) bool); ok {
		r0 = rf(ctx, addr, amt)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// BankKeeper_HasBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasBalance'
type BankKeeper_HasBalance_Call struct {
	*mock.Call
}

// HasBalance is a helper method to define mock.On call
//   - ctx context.Context
//   - addr types.AccAddress
//   - amt types.Coin
func (_e *BankKeeper_Expecter) HasBalance(ctx interface{}, addr interface{}, amt interface{}) *BankKeeper_HasBalance_Call {
	return &BankKeeper_HasBalance_Call{Call: _e.mock.On("HasBalance", ctx, addr, amt)}
}

func (_c *BankKeeper_HasBalance_Call) Run(run func(ctx context.Context, addr types.AccAddress, amt types.Coin)) *BankKeeper_HasBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.AccAddress), args[2].(types.Coin))
	})
	return _c
}

func (_c *BankKeeper_HasBalance_Call) Return(_a0 bool) *BankKeeper_HasBalance_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_HasBalance_Call) RunAndReturn(run func(context.Context, types.AccAddress, types.Coin) bool) *BankKeeper_HasBalance_Call {
	_c.Call.Return(run)
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

// SendCoinsFromModuleToAccount provides a mock function with given fields: ctx, senderModule, recipientAddr, amt
func (_m *BankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientAddr, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromModuleToAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.AccAddress, types.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BankKeeper_SendCoinsFromModuleToAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendCoinsFromModuleToAccount'
type BankKeeper_SendCoinsFromModuleToAccount_Call struct {
	*mock.Call
}

// SendCoinsFromModuleToAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - senderModule string
//   - recipientAddr types.AccAddress
//   - amt types.Coins
func (_e *BankKeeper_Expecter) SendCoinsFromModuleToAccount(ctx interface{}, senderModule interface{}, recipientAddr interface{}, amt interface{}) *BankKeeper_SendCoinsFromModuleToAccount_Call {
	return &BankKeeper_SendCoinsFromModuleToAccount_Call{Call: _e.mock.On("SendCoinsFromModuleToAccount", ctx, senderModule, recipientAddr, amt)}
}

func (_c *BankKeeper_SendCoinsFromModuleToAccount_Call) Run(run func(ctx context.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins)) *BankKeeper_SendCoinsFromModuleToAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(types.AccAddress), args[3].(types.Coins))
	})
	return _c
}

func (_c *BankKeeper_SendCoinsFromModuleToAccount_Call) Return(_a0 error) *BankKeeper_SendCoinsFromModuleToAccount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_SendCoinsFromModuleToAccount_Call) RunAndReturn(run func(context.Context, string, types.AccAddress, types.Coins) error) *BankKeeper_SendCoinsFromModuleToAccount_Call {
	_c.Call.Return(run)
	return _c
}

// SendCoinsFromModuleToModule provides a mock function with given fields: ctx, senderModule, recipientModule, amt
func (_m *BankKeeper) SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt types.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientModule, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromModuleToModule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, types.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BankKeeper_SendCoinsFromModuleToModule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendCoinsFromModuleToModule'
type BankKeeper_SendCoinsFromModuleToModule_Call struct {
	*mock.Call
}

// SendCoinsFromModuleToModule is a helper method to define mock.On call
//   - ctx context.Context
//   - senderModule string
//   - recipientModule string
//   - amt types.Coins
func (_e *BankKeeper_Expecter) SendCoinsFromModuleToModule(ctx interface{}, senderModule interface{}, recipientModule interface{}, amt interface{}) *BankKeeper_SendCoinsFromModuleToModule_Call {
	return &BankKeeper_SendCoinsFromModuleToModule_Call{Call: _e.mock.On("SendCoinsFromModuleToModule", ctx, senderModule, recipientModule, amt)}
}

func (_c *BankKeeper_SendCoinsFromModuleToModule_Call) Run(run func(ctx context.Context, senderModule string, recipientModule string, amt types.Coins)) *BankKeeper_SendCoinsFromModuleToModule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(types.Coins))
	})
	return _c
}

func (_c *BankKeeper_SendCoinsFromModuleToModule_Call) Return(_a0 error) *BankKeeper_SendCoinsFromModuleToModule_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BankKeeper_SendCoinsFromModuleToModule_Call) RunAndReturn(run func(context.Context, string, string, types.Coins) error) *BankKeeper_SendCoinsFromModuleToModule_Call {
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
