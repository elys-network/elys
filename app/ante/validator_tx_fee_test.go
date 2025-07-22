package ante

import (
	"math"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	parametertypes "github.com/elys-network/elys/v6/x/parameter/types"
)

// Mock implementation of sdk.FeeTx for testing
type MockFeeTx struct {
	msgs     []sdk.Msg
	fee      sdk.Coins
	gas      uint64
	feePayer sdk.AccAddress
}

func (tx MockFeeTx) GetMsgs() []sdk.Msg { return tx.msgs }
func (tx MockFeeTx) GetMsgsV2() ([]proto.Message, error) {
	protoMsgs := make([]proto.Message, len(tx.msgs))
	for i, msg := range tx.msgs {
		if protoMsg, ok := msg.(proto.Message); ok {
			protoMsgs[i] = protoMsg
		} else {
			return nil, nil
		}
	}
	return protoMsgs, nil
}
func (tx MockFeeTx) ValidateBasic() error                             { return nil }
func (tx MockFeeTx) GetFee() sdk.Coins                               { return tx.fee }
func (tx MockFeeTx) GetGas() uint64                                  { return tx.gas }
func (tx MockFeeTx) FeePayer() []byte                                { return tx.feePayer }
func (tx MockFeeTx) FeeGranter() []byte                              { return nil }
func (tx MockFeeTx) GetSigners() ([][]byte, error)                   { return [][]byte{tx.feePayer}, nil }
func (tx MockFeeTx) GetPubKeys() ([]cryptotypes.PubKey, error)       { return nil, nil }
func (tx MockFeeTx) GetSignaturesV2() ([]signing.SignatureV2, error) { return nil, nil }

// Mock transaction that doesn't implement FeeTx
type MockInvalidTx struct{}

func (tx MockInvalidTx) GetMsgs() []sdk.Msg                             { return nil }
func (tx MockInvalidTx) GetMsgsV2() ([]proto.Message, error)            { return nil, nil }
func (tx MockInvalidTx) ValidateBasic() error                           { return nil }
func (tx MockInvalidTx) GetSigners() ([][]byte, error)                  { return nil, nil }
func (tx MockInvalidTx) GetPubKeys() ([]cryptotypes.PubKey, error)      { return nil, nil }
func (tx MockInvalidTx) GetSignaturesV2() ([]signing.SignatureV2, error) { return nil, nil }

// Helper function to create a mock transaction
func createTestTx(from sdk.AccAddress, fee sdk.Coins, gas uint64) sdk.Tx {
	privKey := secp256k1.GenPrivKey()
	recipientAddr := sdk.AccAddress(privKey.PubKey().Address())
	
	msg := banktypes.NewMsgSend(from, recipientAddr, sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(100))))
	
	return &MockFeeTx{
		msgs:     []sdk.Msg{msg},
		fee:      fee,
		gas:      gas,
		feePayer: from,
	}
}

func TestCheckTxFeeWithValidatorMinGasPrices(t *testing.T) {
	// Create test addresses
	regularAddr := sdk.AccAddress("regular_address____")
	
	// Create a mock gasless address that matches the expected format
	// We'll temporarily add it to the GaslessAddrs for testing
	mockGaslessAddr := sdk.AccAddress("gasless_address_____")
	originalGaslessAddrs := GaslessAddrs
	GaslessAddrs = append(GaslessAddrs, mockGaslessAddr.String())
	defer func() { GaslessAddrs = originalGaslessAddrs }() // Restore original after test

	for _, tc := range []struct {
		desc              string
		fromAddr          sdk.AccAddress
		fee               sdk.Coins
		gas               uint64
		minGasPrices      sdk.DecCoins
		isCheckTx         bool
		expectError       bool
		expectedPriority  int64
		errorContains     string
	}{
		{
			desc:             "Regular transaction with sufficient fees",
			fromAddr:         regularAddr,
			fee:              sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(50000))),
			gas:              50000,
			minGasPrices:     sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
			isCheckTx:        true,
			expectError:      false,
			expectedPriority: 1, // 50000/50000 = 1
		},
		{
			desc:          "Regular transaction with insufficient fees",
			fromAddr:      regularAddr,
			fee:           sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(10000))),
			gas:           50000,
			minGasPrices:  sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
			isCheckTx:     true,
			expectError:   true,
			errorContains: "insufficient fees",
		},
		{
			desc:             "Gasless transaction with no fees",
			fromAddr:         mockGaslessAddr,
			fee:              sdk.NewCoins(),
			gas:              50000,
			minGasPrices:     sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
			isCheckTx:        true,
			expectError:      false,
			expectedPriority: math.MaxInt64, // Should get max priority
		},
		{
			desc:             "Gasless transaction with fees (still gets max priority)",
			fromAddr:         mockGaslessAddr,
			fee:              sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(10000))),
			gas:              50000,
			minGasPrices:     sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
			isCheckTx:        true,
			expectError:      false,
			expectedPriority: math.MaxInt64, // Should still get max priority
		},
		{
			desc:             "Gasless transaction bypasses high min gas prices",
			fromAddr:         mockGaslessAddr,
			fee:              sdk.NewCoins(),
			gas:              50000,
			minGasPrices:     sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1000))),
			isCheckTx:        true,
			expectError:      false,
			expectedPriority: math.MaxInt64, // Should get max priority
		},
		{
			desc:             "Non-CheckTx mode doesn't validate fees",
			fromAddr:         regularAddr,
			fee:              sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(10000))),
			gas:              50000,
			minGasPrices:     sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
			isCheckTx:        false,
			expectError:      false,
			expectedPriority: 0, // 10000/50000 = 0 (rounded down)
		},
		{
			desc:             "Multiple coin fees - minimum priority used",
			fromAddr:         regularAddr,
			fee:              sdk.NewCoins(
				sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(50000)),
				sdk.NewCoin("uusdc", sdkmath.NewInt(25000)),
			),
			gas:              50000,
			minGasPrices:     sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
			isCheckTx:        true,
			expectError:      false,
			expectedPriority: 0, // min(50000/50000, 25000/50000) = min(1, 0) = 0
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := sdk.NewContext(nil, cmtproto.Header{}, false, log.NewNopLogger())
			ctx = ctx.WithIsCheckTx(tc.isCheckTx).WithMinGasPrices(tc.minGasPrices)
			tx := createTestTx(tc.fromAddr, tc.fee, tc.gas)
			
			resultFee, priority, err := CheckTxFeeWithValidatorMinGasPrices(ctx, tx)
			
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.fee, resultFee)
				require.Equal(t, tc.expectedPriority, priority)
			}
		})
	}
}

func TestCheckTxFeeWithInvalidTransaction(t *testing.T) {
	ctx := sdk.NewContext(nil, cmtproto.Header{}, false, log.NewNopLogger())
	
	invalidTx := &MockInvalidTx{}
	
	_, _, err := CheckTxFeeWithValidatorMinGasPrices(ctx, invalidTx)
	
	require.Error(t, err)
	require.Contains(t, err.Error(), "tx must implement FeeTx")
}

func TestGaslessTransactionPriorityComparison(t *testing.T) {
	// Test that gasless transactions have higher priority than regular transactions
	ctx := sdk.NewContext(nil, cmtproto.Header{}, false, log.NewNopLogger())
	ctx = ctx.WithIsCheckTx(true).WithMinGasPrices(
		sdk.NewDecCoins(sdk.NewDecCoinFromDec(parametertypes.Elys, sdkmath.LegacyNewDec(1))),
	)
	
	// Create addresses
	regularAddr := sdk.AccAddress("regular_address____")
	mockGaslessAddr := sdk.AccAddress("gasless_address_____")
	
	// Temporarily add gasless address to whitelist
	originalGaslessAddrs := GaslessAddrs
	GaslessAddrs = append(GaslessAddrs, mockGaslessAddr.String())
	defer func() { GaslessAddrs = originalGaslessAddrs }()
	
	// Regular transaction with very high fees
	regularTx := createTestTx(regularAddr, sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(1000000))), 50000)
	
	// Gasless transaction with no fees
	gaslessTx := createTestTx(mockGaslessAddr, sdk.NewCoins(), 50000)
	
	// Check regular transaction priority
	_, regularPriority, err := CheckTxFeeWithValidatorMinGasPrices(ctx, regularTx)
	require.NoError(t, err)
	
	// Check gasless transaction priority
	_, gaslessPriority, err := CheckTxFeeWithValidatorMinGasPrices(ctx, gaslessTx)
	require.NoError(t, err)
	
	// Gasless should have higher priority
	require.Greater(t, gaslessPriority, regularPriority)
	require.Equal(t, int64(math.MaxInt64), gaslessPriority)
}

func TestGetTxPriority(t *testing.T) {
	for _, tc := range []struct {
		desc             string
		fee              sdk.Coins
		gas              int64
		expectedPriority int64
	}{
		{
			desc:             "Single coin fee",
			fee:              sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(100000))),
			gas:              50000,
			expectedPriority: 2, // 100000/50000 = 2
		},
		{
			desc:             "Multiple coin fees - minimum used",
			fee:              sdk.NewCoins(
				sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(100000)),
				sdk.NewCoin("uusdc", sdkmath.NewInt(25000)),
			),
			gas:              50000,
			expectedPriority: 0, // min(100000/50000, 25000/50000) = min(2, 0) = 0
		},
		{
			desc:             "Zero fees",
			fee:              sdk.NewCoins(),
			gas:              50000,
			expectedPriority: 0,
		},
		{
			desc:             "High gas price",
			fee:              sdk.NewCoins(sdk.NewCoin(parametertypes.Elys, sdkmath.NewInt(1000000000))),
			gas:              50000,
			expectedPriority: 20000, // 1000000000/50000 = 20000
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			priority := getTxPriority(tc.fee, tc.gas)
			require.Equal(t, tc.expectedPriority, priority)
		})
	}
} 