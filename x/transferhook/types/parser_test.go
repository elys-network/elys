package types_test

import (
	"encoding/json"
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/transferhook/types"
)

func getAmmMemo(address, action string, routes []ammtypes.SwapAmountInRoute) string {
	routesText, err := json.Marshal(routes)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`
		{
			"transferhook": {
				"receiver": "%[1]s",
				"amm": {
					"action": "%[2]s",
					"routes": %s
				} 
			}
		}`, address, action, routesText)
}

// Helper function to check the routingInfo with a switch statement
// This isn't the most efficient way to check the type  (require.TypeOf could be used instead)
// but it better aligns with how the routing info is checked in module_ibc
func checkModuleRoutingInfoType(routingInfo types.ModuleRoutingInfo, expectedType string) bool {
	switch routingInfo.(type) {
	case types.AmmPacketMetadata:
		return expectedType == "amm"
	default:
		return false
	}
}

func TestParsePacketMetadata(t *testing.T) {
	validAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()
	invalidAddress := "invalid_address"
	validAmmAction := "Swap"

	validParsedAmmPacketMetadata := types.AmmPacketMetadata{
		Action: validAmmAction,
		Routes: []ammtypes.SwapAmountInRoute{
			{
				PoolId:        1,
				TokenOutDenom: "uelys",
			},
		},
	}

	testCases := []struct {
		name                string
		metadata            string
		parsedAmm           *types.AmmPacketMetadata
		expectedNilMetadata bool
		expectedErr         string
	}{
		{
			name:      "valid amm memo",
			metadata:  getAmmMemo(validAddress, validAmmAction, validParsedAmmPacketMetadata.Routes),
			parsedAmm: &validParsedAmmPacketMetadata,
		},
		{
			name:                "normal IBC transfer",
			metadata:            validAddress, // normal address - not transferhook JSON
			expectedNilMetadata: true,
		},
		{
			name:                "empty memo",
			metadata:            "",
			expectedNilMetadata: true,
		},
		{
			name:                "empty JSON memo",
			metadata:            "{}",
			expectedNilMetadata: true,
		},
		{
			name:                "different module specified",
			metadata:            `{ "other_module": { } }`,
			expectedNilMetadata: true,
		},
		{
			name:        "empty receiver address",
			metadata:    `{ "transferhook": { } }`,
			expectedErr: "receiver address must be specified when using transferhook",
		},
		{
			name:        "invalid receiver address",
			metadata:    `{ "transferhook": { "receiver": "invalid_address" } }`,
			expectedErr: "receiver address must be specified when using transferhook",
		},
		{
			name:        "invalid amm address",
			metadata:    getAmmMemo(invalidAddress, validAmmAction, validParsedAmmPacketMetadata.Routes),
			expectedErr: "receiver address must be specified when using transferhook",
		},
		{
			name:        "invalid amm action",
			metadata:    getAmmMemo(validAddress, "bad_action", validParsedAmmPacketMetadata.Routes),
			expectedErr: "unsupported amm action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedData, actualErr := types.ParsePacketMetadata(tc.metadata)

			if tc.expectedErr == "" {
				require.NoError(t, actualErr)
				if tc.expectedNilMetadata {
					require.Nil(t, parsedData, "parsed data response should be nil")
				} else {
					if tc.parsedAmm != nil {
						checkModuleRoutingInfoType(parsedData.RoutingInfo, "amm")
						routingInfo, ok := parsedData.RoutingInfo.(types.AmmPacketMetadata)
						require.True(t, ok, "routing info should be amm")
						require.Equal(t, *tc.parsedAmm, routingInfo, "parsed amm value")
					}
				}
			} else {
				require.ErrorContains(t, actualErr, types.ErrInvalidPacketMetadata.Error(), "expected error type for %s", tc.name)
				require.ErrorContains(t, actualErr, tc.expectedErr, "expected error for %s", tc.name)
			}
		})
	}
}

func TestValidateAmmPacketMetadata(t *testing.T) {
	validAction := "Swap"

	testCases := []struct {
		name        string
		metadata    *types.AmmPacketMetadata
		expectedErr string
	}{
		{
			name: "valid Metadata data",
			metadata: &types.AmmPacketMetadata{
				Action: validAction,
			},
		},
		{
			name: "invalid action",
			metadata: &types.AmmPacketMetadata{
				Action: "bad_action",
			},
			expectedErr: "unsupported amm action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualErr := tc.metadata.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, actualErr, "no error expected for %s", tc.name)
			} else {
				require.ErrorContains(t, actualErr, tc.expectedErr, "error expected for %s", tc.name)
			}
		})
	}
}
