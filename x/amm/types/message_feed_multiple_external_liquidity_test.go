package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestMsgFeedMultipleExternalLiquidity_ValidateBasic(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		msg            types.MsgFeedMultipleExternalLiquidity
		expectedErrMsg string
	}{
		{
			"Invalid address",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: "elys11",
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "tokenA", Depth: sdkmath.LegacyNewDec(1000), Amount: sdkmath.LegacyNewDec(500)},
							{Asset: "tokenB", Depth: sdkmath.LegacyNewDec(2000), Amount: sdkmath.LegacyNewDec(1000)},
						},
					},
				},
			},
			"invalid sender address",
		},
		{
			"valid message",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: sample.AccAddress(),
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "tokenA", Depth: sdkmath.LegacyNewDec(1000), Amount: sdkmath.LegacyNewDec(500)},
							{Asset: "tokenB", Depth: sdkmath.LegacyNewDec(2000), Amount: sdkmath.LegacyNewDec(1000)},
						},
					},
				},
			},
			"",
		},
		{
			"invalid asset denom",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: sample.AccAddress(),
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "", Depth: sdkmath.LegacyNewDec(1000), Amount: sdkmath.LegacyNewDec(500)},
						},
					},
				},
			},
			"asset cannot be empty",
		},
		{
			"negative depth",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: sample.AccAddress(),
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "tokenA", Depth: sdkmath.LegacyNewDec(-1000), Amount: sdkmath.LegacyNewDec(500)},
						},
					},
				},
			},
			"depth cannot be negative or nil",
		},
		{
			"negative amount",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: sample.AccAddress(),
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "tokenA", Depth: sdkmath.LegacyNewDec(1000), Amount: sdkmath.LegacyNewDec(-500)},
						},
					},
				},
			},
			"depth amount cannot be negative or nil",
		},
		{
			"nil depth",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: sample.AccAddress(),
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "tokenA", Depth: sdkmath.LegacyDec{}, Amount: sdkmath.LegacyNewDec(500)},
						},
					},
				},
			},
			"depth cannot be negative or nil",
		},
		{
			"nil amount",
			types.MsgFeedMultipleExternalLiquidity{
				Sender: sample.AccAddress(),
				Liquidity: []types.ExternalLiquidity{
					{
						AmountDepthInfo: []types.AssetAmountDepth{
							{Asset: "tokenA", Depth: sdkmath.LegacyNewDec(1000), Amount: sdkmath.LegacyDec{}},
						},
					},
				},
			},
			"depth amount cannot be negative or nil",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
