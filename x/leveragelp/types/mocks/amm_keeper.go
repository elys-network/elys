// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	ammtypes "github.com/elys-network/elys/x/amm/types"

	math "cosmossdk.io/math"

	mock "github.com/stretchr/testify/mock"

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

// CalcInAmtGivenOut provides a mock function with given fields: ctx, poolId, oracle, snapshot, tokensOut, tokenInDenom, swapFee
func (_m *AmmKeeper) CalcInAmtGivenOut(ctx types.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut types.Coins, tokenInDenom string, swapFee math.LegacyDec) (types.Coin, error) {
	ret := _m.Called(ctx, poolId, oracle, snapshot, tokensOut, tokenInDenom, swapFee)

	var r0 types.Coin
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) (types.Coin, error)); ok {
		return rf(ctx, poolId, oracle, snapshot, tokensOut, tokenInDenom, swapFee)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) types.Coin); ok {
		r0 = rf(ctx, poolId, oracle, snapshot, tokensOut, tokenInDenom, swapFee)
	} else {
		r0 = ret.Get(0).(types.Coin)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) error); ok {
		r1 = rf(ctx, poolId, oracle, snapshot, tokensOut, tokenInDenom, swapFee)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AmmKeeper_CalcInAmtGivenOut_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CalcInAmtGivenOut'
type AmmKeeper_CalcInAmtGivenOut_Call struct {
	*mock.Call
}

// CalcInAmtGivenOut is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
//   - oracle ammtypes.OracleKeeper
//   - snapshot *ammtypes.Pool
//   - tokensOut types.Coins
//   - tokenInDenom string
//   - swapFee math.LegacyDec
func (_e *AmmKeeper_Expecter) CalcInAmtGivenOut(ctx interface{}, poolId interface{}, oracle interface{}, snapshot interface{}, tokensOut interface{}, tokenInDenom interface{}, swapFee interface{}) *AmmKeeper_CalcInAmtGivenOut_Call {
	return &AmmKeeper_CalcInAmtGivenOut_Call{Call: _e.mock.On("CalcInAmtGivenOut", ctx, poolId, oracle, snapshot, tokensOut, tokenInDenom, swapFee)}
}

func (_c *AmmKeeper_CalcInAmtGivenOut_Call) Run(run func(ctx types.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut types.Coins, tokenInDenom string, swapFee math.LegacyDec)) *AmmKeeper_CalcInAmtGivenOut_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64), args[2].(ammtypes.OracleKeeper), args[3].(*ammtypes.Pool), args[4].(types.Coins), args[5].(string), args[6].(math.LegacyDec))
	})
	return _c
}

func (_c *AmmKeeper_CalcInAmtGivenOut_Call) Return(tokenIn types.Coin, err error) *AmmKeeper_CalcInAmtGivenOut_Call {
	_c.Call.Return(tokenIn, err)
	return _c
}

func (_c *AmmKeeper_CalcInAmtGivenOut_Call) RunAndReturn(run func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) (types.Coin, error)) *AmmKeeper_CalcInAmtGivenOut_Call {
	_c.Call.Return(run)
	return _c
}

// CalcOutAmtGivenIn provides a mock function with given fields: ctx, poolId, oracle, snapshot, tokensIn, tokenOutDenom, swapFee
func (_m *AmmKeeper) CalcOutAmtGivenIn(ctx types.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn types.Coins, tokenOutDenom string, swapFee math.LegacyDec) (types.Coin, error) {
	ret := _m.Called(ctx, poolId, oracle, snapshot, tokensIn, tokenOutDenom, swapFee)

	var r0 types.Coin
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) (types.Coin, error)); ok {
		return rf(ctx, poolId, oracle, snapshot, tokensIn, tokenOutDenom, swapFee)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) types.Coin); ok {
		r0 = rf(ctx, poolId, oracle, snapshot, tokensIn, tokenOutDenom, swapFee)
	} else {
		r0 = ret.Get(0).(types.Coin)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) error); ok {
		r1 = rf(ctx, poolId, oracle, snapshot, tokensIn, tokenOutDenom, swapFee)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AmmKeeper_CalcOutAmtGivenIn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CalcOutAmtGivenIn'
type AmmKeeper_CalcOutAmtGivenIn_Call struct {
	*mock.Call
}

// CalcOutAmtGivenIn is a helper method to define mock.On call
//   - ctx types.Context
//   - poolId uint64
//   - oracle ammtypes.OracleKeeper
//   - snapshot *ammtypes.Pool
//   - tokensIn types.Coins
//   - tokenOutDenom string
//   - swapFee math.LegacyDec
func (_e *AmmKeeper_Expecter) CalcOutAmtGivenIn(ctx interface{}, poolId interface{}, oracle interface{}, snapshot interface{}, tokensIn interface{}, tokenOutDenom interface{}, swapFee interface{}) *AmmKeeper_CalcOutAmtGivenIn_Call {
	return &AmmKeeper_CalcOutAmtGivenIn_Call{Call: _e.mock.On("CalcOutAmtGivenIn", ctx, poolId, oracle, snapshot, tokensIn, tokenOutDenom, swapFee)}
}

func (_c *AmmKeeper_CalcOutAmtGivenIn_Call) Run(run func(ctx types.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn types.Coins, tokenOutDenom string, swapFee math.LegacyDec)) *AmmKeeper_CalcOutAmtGivenIn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64), args[2].(ammtypes.OracleKeeper), args[3].(*ammtypes.Pool), args[4].(types.Coins), args[5].(string), args[6].(math.LegacyDec))
	})
	return _c
}

func (_c *AmmKeeper_CalcOutAmtGivenIn_Call) Return(_a0 types.Coin, _a1 error) *AmmKeeper_CalcOutAmtGivenIn_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AmmKeeper_CalcOutAmtGivenIn_Call) RunAndReturn(run func(types.Context, uint64, ammtypes.OracleKeeper, *ammtypes.Pool, types.Coins, string, math.LegacyDec) (types.Coin, error)) *AmmKeeper_CalcOutAmtGivenIn_Call {
	_c.Call.Return(run)
	return _c
}

// ExitPool provides a mock function with given fields: ctx, sender, poolId, shareInAmount, tokenOutMins, tokenOutDenom
func (_m *AmmKeeper) ExitPool(ctx types.Context, sender types.AccAddress, poolId uint64, shareInAmount math.Int, tokenOutMins types.Coins, tokenOutDenom string) (types.Coins, error) {
	ret := _m.Called(ctx, sender, poolId, shareInAmount, tokenOutMins, tokenOutDenom)

	var r0 types.Coins
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, string) (types.Coins, error)); ok {
		return rf(ctx, sender, poolId, shareInAmount, tokenOutMins, tokenOutDenom)
	}
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, string) types.Coins); ok {
		r0 = rf(ctx, sender, poolId, shareInAmount, tokenOutMins, tokenOutDenom)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Coins)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, string) error); ok {
		r1 = rf(ctx, sender, poolId, shareInAmount, tokenOutMins, tokenOutDenom)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AmmKeeper_ExitPool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExitPool'
type AmmKeeper_ExitPool_Call struct {
	*mock.Call
}

// ExitPool is a helper method to define mock.On call
//   - ctx types.Context
//   - sender types.AccAddress
//   - poolId uint64
//   - shareInAmount math.Int
//   - tokenOutMins types.Coins
//   - tokenOutDenom string
func (_e *AmmKeeper_Expecter) ExitPool(ctx interface{}, sender interface{}, poolId interface{}, shareInAmount interface{}, tokenOutMins interface{}, tokenOutDenom interface{}) *AmmKeeper_ExitPool_Call {
	return &AmmKeeper_ExitPool_Call{Call: _e.mock.On("ExitPool", ctx, sender, poolId, shareInAmount, tokenOutMins, tokenOutDenom)}
}

func (_c *AmmKeeper_ExitPool_Call) Run(run func(ctx types.Context, sender types.AccAddress, poolId uint64, shareInAmount math.Int, tokenOutMins types.Coins, tokenOutDenom string)) *AmmKeeper_ExitPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(types.AccAddress), args[2].(uint64), args[3].(math.Int), args[4].(types.Coins), args[5].(string))
	})
	return _c
}

func (_c *AmmKeeper_ExitPool_Call) Return(exitCoins types.Coins, err error) *AmmKeeper_ExitPool_Call {
	_c.Call.Return(exitCoins, err)
	return _c
}

func (_c *AmmKeeper_ExitPool_Call) RunAndReturn(run func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, string) (types.Coins, error)) *AmmKeeper_ExitPool_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllPool provides a mock function with given fields: _a0
func (_m *AmmKeeper) GetAllPool(_a0 types.Context) []ammtypes.Pool {
	ret := _m.Called(_a0)

	var r0 []ammtypes.Pool
	if rf, ok := ret.Get(0).(func(types.Context) []ammtypes.Pool); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ammtypes.Pool)
		}
	}

	return r0
}

// AmmKeeper_GetAllPool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllPool'
type AmmKeeper_GetAllPool_Call struct {
	*mock.Call
}

// GetAllPool is a helper method to define mock.On call
//   - _a0 types.Context
func (_e *AmmKeeper_Expecter) GetAllPool(_a0 interface{}) *AmmKeeper_GetAllPool_Call {
	return &AmmKeeper_GetAllPool_Call{Call: _e.mock.On("GetAllPool", _a0)}
}

func (_c *AmmKeeper_GetAllPool_Call) Run(run func(_a0 types.Context)) *AmmKeeper_GetAllPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context))
	})
	return _c
}

func (_c *AmmKeeper_GetAllPool_Call) Return(_a0 []ammtypes.Pool) *AmmKeeper_GetAllPool_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmmKeeper_GetAllPool_Call) RunAndReturn(run func(types.Context) []ammtypes.Pool) *AmmKeeper_GetAllPool_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllPoolIdsWithDenom provides a mock function with given fields: _a0, _a1
func (_m *AmmKeeper) GetAllPoolIdsWithDenom(_a0 types.Context, _a1 string) []uint64 {
	ret := _m.Called(_a0, _a1)

	var r0 []uint64
	if rf, ok := ret.Get(0).(func(types.Context, string) []uint64); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint64)
		}
	}

	return r0
}

// AmmKeeper_GetAllPoolIdsWithDenom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllPoolIdsWithDenom'
type AmmKeeper_GetAllPoolIdsWithDenom_Call struct {
	*mock.Call
}

// GetAllPoolIdsWithDenom is a helper method to define mock.On call
//   - _a0 types.Context
//   - _a1 string
func (_e *AmmKeeper_Expecter) GetAllPoolIdsWithDenom(_a0 interface{}, _a1 interface{}) *AmmKeeper_GetAllPoolIdsWithDenom_Call {
	return &AmmKeeper_GetAllPoolIdsWithDenom_Call{Call: _e.mock.On("GetAllPoolIdsWithDenom", _a0, _a1)}
}

func (_c *AmmKeeper_GetAllPoolIdsWithDenom_Call) Run(run func(_a0 types.Context, _a1 string)) *AmmKeeper_GetAllPoolIdsWithDenom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(string))
	})
	return _c
}

func (_c *AmmKeeper_GetAllPoolIdsWithDenom_Call) Return(_a0 []uint64) *AmmKeeper_GetAllPoolIdsWithDenom_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmmKeeper_GetAllPoolIdsWithDenom_Call) RunAndReturn(run func(types.Context, string) []uint64) *AmmKeeper_GetAllPoolIdsWithDenom_Call {
	_c.Call.Return(run)
	return _c
}

// GetPool provides a mock function with given fields: _a0, _a1
func (_m *AmmKeeper) GetPool(_a0 types.Context, _a1 uint64) (ammtypes.Pool, bool) {
	ret := _m.Called(_a0, _a1)

	var r0 ammtypes.Pool
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, uint64) (ammtypes.Pool, bool)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(types.Context, uint64) ammtypes.Pool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(ammtypes.Pool)
	}

	if rf, ok := ret.Get(1).(func(types.Context, uint64) bool); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// AmmKeeper_GetPool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPool'
type AmmKeeper_GetPool_Call struct {
	*mock.Call
}

// GetPool is a helper method to define mock.On call
//   - _a0 types.Context
//   - _a1 uint64
func (_e *AmmKeeper_Expecter) GetPool(_a0 interface{}, _a1 interface{}) *AmmKeeper_GetPool_Call {
	return &AmmKeeper_GetPool_Call{Call: _e.mock.On("GetPool", _a0, _a1)}
}

func (_c *AmmKeeper_GetPool_Call) Run(run func(_a0 types.Context, _a1 uint64)) *AmmKeeper_GetPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(uint64))
	})
	return _c
}

func (_c *AmmKeeper_GetPool_Call) Return(_a0 ammtypes.Pool, _a1 bool) *AmmKeeper_GetPool_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AmmKeeper_GetPool_Call) RunAndReturn(run func(types.Context, uint64) (ammtypes.Pool, bool)) *AmmKeeper_GetPool_Call {
	_c.Call.Return(run)
	return _c
}

// GetPoolSnapshotOrSet provides a mock function with given fields: ctx, pool
func (_m *AmmKeeper) GetPoolSnapshotOrSet(ctx types.Context, pool ammtypes.Pool) ammtypes.Pool {
	ret := _m.Called(ctx, pool)

	var r0 ammtypes.Pool
	if rf, ok := ret.Get(0).(func(types.Context, ammtypes.Pool) ammtypes.Pool); ok {
		r0 = rf(ctx, pool)
	} else {
		r0 = ret.Get(0).(ammtypes.Pool)
	}

	return r0
}

// AmmKeeper_GetPoolSnapshotOrSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPoolSnapshotOrSet'
type AmmKeeper_GetPoolSnapshotOrSet_Call struct {
	*mock.Call
}

// GetPoolSnapshotOrSet is a helper method to define mock.On call
//   - ctx types.Context
//   - pool ammtypes.Pool
func (_e *AmmKeeper_Expecter) GetPoolSnapshotOrSet(ctx interface{}, pool interface{}) *AmmKeeper_GetPoolSnapshotOrSet_Call {
	return &AmmKeeper_GetPoolSnapshotOrSet_Call{Call: _e.mock.On("GetPoolSnapshotOrSet", ctx, pool)}
}

func (_c *AmmKeeper_GetPoolSnapshotOrSet_Call) Run(run func(ctx types.Context, pool ammtypes.Pool)) *AmmKeeper_GetPoolSnapshotOrSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(ammtypes.Pool))
	})
	return _c
}

func (_c *AmmKeeper_GetPoolSnapshotOrSet_Call) Return(val ammtypes.Pool) *AmmKeeper_GetPoolSnapshotOrSet_Call {
	_c.Call.Return(val)
	return _c
}

func (_c *AmmKeeper_GetPoolSnapshotOrSet_Call) RunAndReturn(run func(types.Context, ammtypes.Pool) ammtypes.Pool) *AmmKeeper_GetPoolSnapshotOrSet_Call {
	_c.Call.Return(run)
	return _c
}

// IterateLiquidityPools provides a mock function with given fields: _a0, _a1
func (_m *AmmKeeper) IterateLiquidityPools(_a0 types.Context, _a1 func(ammtypes.Pool) bool) {
	_m.Called(_a0, _a1)
}

// AmmKeeper_IterateLiquidityPools_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IterateLiquidityPools'
type AmmKeeper_IterateLiquidityPools_Call struct {
	*mock.Call
}

// IterateLiquidityPools is a helper method to define mock.On call
//   - _a0 types.Context
//   - _a1 func(ammtypes.Pool) bool
func (_e *AmmKeeper_Expecter) IterateLiquidityPools(_a0 interface{}, _a1 interface{}) *AmmKeeper_IterateLiquidityPools_Call {
	return &AmmKeeper_IterateLiquidityPools_Call{Call: _e.mock.On("IterateLiquidityPools", _a0, _a1)}
}

func (_c *AmmKeeper_IterateLiquidityPools_Call) Run(run func(_a0 types.Context, _a1 func(ammtypes.Pool) bool)) *AmmKeeper_IterateLiquidityPools_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(func(ammtypes.Pool) bool))
	})
	return _c
}

func (_c *AmmKeeper_IterateLiquidityPools_Call) Return() *AmmKeeper_IterateLiquidityPools_Call {
	_c.Call.Return()
	return _c
}

func (_c *AmmKeeper_IterateLiquidityPools_Call) RunAndReturn(run func(types.Context, func(ammtypes.Pool) bool)) *AmmKeeper_IterateLiquidityPools_Call {
	_c.Call.Return(run)
	return _c
}

// JoinPoolNoSwap provides a mock function with given fields: ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining
func (_m *AmmKeeper) JoinPoolNoSwap(ctx types.Context, sender types.AccAddress, poolId uint64, shareOutAmount math.Int, tokenInMaxs types.Coins, noRemaining bool) (types.Coins, math.Int, error) {
	ret := _m.Called(ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining)

	var r0 types.Coins
	var r1 math.Int
	var r2 error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, bool) (types.Coins, math.Int, error)); ok {
		return rf(ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining)
	}
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, bool) types.Coins); ok {
		r0 = rf(ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Coins)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, bool) math.Int); ok {
		r1 = rf(ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining)
	} else {
		r1 = ret.Get(1).(math.Int)
	}

	if rf, ok := ret.Get(2).(func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, bool) error); ok {
		r2 = rf(ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// AmmKeeper_JoinPoolNoSwap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'JoinPoolNoSwap'
type AmmKeeper_JoinPoolNoSwap_Call struct {
	*mock.Call
}

// JoinPoolNoSwap is a helper method to define mock.On call
//   - ctx types.Context
//   - sender types.AccAddress
//   - poolId uint64
//   - shareOutAmount math.Int
//   - tokenInMaxs types.Coins
//   - noRemaining bool
func (_e *AmmKeeper_Expecter) JoinPoolNoSwap(ctx interface{}, sender interface{}, poolId interface{}, shareOutAmount interface{}, tokenInMaxs interface{}, noRemaining interface{}) *AmmKeeper_JoinPoolNoSwap_Call {
	return &AmmKeeper_JoinPoolNoSwap_Call{Call: _e.mock.On("JoinPoolNoSwap", ctx, sender, poolId, shareOutAmount, tokenInMaxs, noRemaining)}
}

func (_c *AmmKeeper_JoinPoolNoSwap_Call) Run(run func(ctx types.Context, sender types.AccAddress, poolId uint64, shareOutAmount math.Int, tokenInMaxs types.Coins, noRemaining bool)) *AmmKeeper_JoinPoolNoSwap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].(types.AccAddress), args[2].(uint64), args[3].(math.Int), args[4].(types.Coins), args[5].(bool))
	})
	return _c
}

func (_c *AmmKeeper_JoinPoolNoSwap_Call) Return(tokenIn types.Coins, sharesOut math.Int, err error) *AmmKeeper_JoinPoolNoSwap_Call {
	_c.Call.Return(tokenIn, sharesOut, err)
	return _c
}

func (_c *AmmKeeper_JoinPoolNoSwap_Call) RunAndReturn(run func(types.Context, types.AccAddress, uint64, math.Int, types.Coins, bool) (types.Coins, math.Int, error)) *AmmKeeper_JoinPoolNoSwap_Call {
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
