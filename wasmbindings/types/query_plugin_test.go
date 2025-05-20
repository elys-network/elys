package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/elys-network/elys/v4/wasmbindings/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	proto "github.com/golang/protobuf/proto" //nolint:staticcheck // we're intentionally using this deprecated package to be compatible with cosmos protos

	// "github.com/osmosis-labs/osmosis/math"
	// "github.com/osmosis-labs/osmosis/v17/app/apptesting"
	// "github.com/osmosis-labs/osmosis/v17/x/gamm/pool-models/balancer"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/v4/app"
	"github.com/stretchr/testify/suite"
)

type StargateTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *simapp.ElysApp
	legacyAmino *codec.LegacyAmino
}

func (suite *StargateTestSuite) SetupTestInternal() {
	app := simapp.InitElysTestApp(true, suite.Suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(true)
	suite.app = app
}

// func (suite *StargateTestSuite) TearDownTestInternal() {
// 	os.RemoveAll(suite.HomeDir)
// }

func TestStargateTestSuite(t *testing.T) {
	suite.Run(t, new(StargateTestSuite))
}

func (suite *StargateTestSuite) AddAccounts(n int, given []sdk.AccAddress) []sdk.AccAddress {
	issueAmount := math.NewInt(10_000_000_000_000)
	var addresses []sdk.AccAddress
	if n > len(given) {
		addresses = simapp.AddTestAddrs(suite.app, suite.ctx, n-len(given), issueAmount)
		addresses = append(addresses, given...)
	} else {
		addresses = given
	}
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin("uatom", issueAmount),
			sdk.NewCoin("uelys", issueAmount),
			sdk.NewCoin("uusdc", issueAmount),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
		if err != nil {
			panic(err)
		}
	}
	return addresses
}

func (suite *StargateTestSuite) TestStargateQuerier() {
	testCases := []struct {
		name                   string
		testSetup              func()
		path                   string
		requestData            func() []byte
		responseProtoStruct    proto.Message
		expectedQuerierError   bool
		expectedUnMarshalError bool
		resendRequest          bool
		checkResponseStruct    bool
	}{
		{
			name: "happy path Params",
			path: "/elys.amm.Query/Params",
			requestData: func() []byte {
				paaramrequest := ammtypes.QueryParamsRequest{}
				bz, err := proto.Marshal(&paaramrequest)
				suite.Require().NoError(err)
				return bz
			},
			responseProtoStruct: &ammtypes.QueryParamsResponse{},
		},
		{
			name: "test query using iterator",
			testSetup: func() {
				types.SetWhitelistedQuery("/cosmos.bank.v1beta1.Query/AllBalances", &banktypes.QueryAllBalancesResponse{})
			},
			path: "/cosmos.bank.v1beta1.Query/AllBalances",
			requestData: func() []byte {
				accAddr := suite.AddAccounts(1, nil)[0]
				bankrequest := banktypes.QueryAllBalancesRequest{
					Address: accAddr.String(),
				}
				bz, err := proto.Marshal(&bankrequest)
				suite.Require().NoError(err)
				return bz
			},
			responseProtoStruct: &banktypes.QueryAllBalancesResponse{},
		},
		{
			name: "edge case: resending request",
			testSetup: func() {
				types.SetWhitelistedQuery("/cosmos.bank.v1beta1.Query/AllBalances", &banktypes.QueryAllBalancesResponse{})
			},
			path: "/cosmos.bank.v1beta1.Query/AllBalances",
			requestData: func() []byte {
				accAddr := suite.AddAccounts(1, nil)[0]
				bankrequest := banktypes.QueryAllBalancesRequest{
					Address: accAddr.String(),
				}
				bz, err := proto.Marshal(&bankrequest)
				suite.Require().NoError(err)
				return bz
			},
			responseProtoStruct: &banktypes.QueryAllBalancesResponse{},
			resendRequest:       true,
		},
		{
			name: "error in grpc querier",
			// set up whitelist with wrong data
			testSetup: func() {
				types.SetWhitelistedQuery("/cosmos.bank.v1beta1.Query/AllBalances", &banktypes.QueryAllBalancesRequest{})
			},
			path: "/cosmos.bank.v1beta1.Query/AllBalances",
			requestData: func() []byte {
				bankrequest := banktypes.QueryAllBalancesRequest{}
				bz, err := proto.Marshal(&bankrequest)
				suite.Require().NoError(err)
				return bz
			},
			responseProtoStruct:  &banktypes.QueryAllBalancesRequest{},
			expectedQuerierError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTestInternal()
			if tc.testSetup != nil {
				tc.testSetup()
			}

			stargateQuerier := types.StargateQuerier(*suite.app.GRPCQueryRouter(), suite.app.AppCodec())
			stargateRequest := &wasmvmtypes.StargateQuery{
				Path: tc.path,
				Data: tc.requestData(),
			}
			stargateResponse, err := stargateQuerier(suite.ctx, stargateRequest)
			if tc.expectedQuerierError {
				suite.Require().Error(err)
				return
			}
			if tc.checkResponseStruct {
				expectedResponse, err := proto.Marshal(tc.responseProtoStruct)
				suite.Require().NoError(err)
				expJsonResp, err := types.ConvertProtoToJSONMarshal(tc.responseProtoStruct, expectedResponse, suite.app.AppCodec())
				suite.Require().NoError(err)
				suite.Require().Equal(expJsonResp, stargateResponse)
			}

			suite.Require().NoError(err)

			protoResponse, ok := tc.responseProtoStruct.(proto.Message)
			suite.Require().True(ok)

			// test correctness by unmarshalling json response into proto struct
			err = suite.app.AppCodec().UnmarshalJSON(stargateResponse, protoResponse)
			if tc.expectedUnMarshalError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(protoResponse)
			}

			if tc.resendRequest {
				stargateQuerier = types.StargateQuerier(*suite.app.GRPCQueryRouter(), suite.app.AppCodec())
				stargateRequest = &wasmvmtypes.StargateQuery{
					Path: tc.path,
					Data: tc.requestData(),
				}
				resendResponse, err := stargateQuerier(suite.ctx, stargateRequest)
				suite.Require().NoError(err)
				suite.Require().Equal(stargateResponse, resendResponse)
			}
		})
	}
}

func (suite *StargateTestSuite) TestConvertProtoToJsonMarshal() {
	testCases := []struct {
		name                  string
		queryPath             string
		protoResponseStruct   proto.Message
		originalResponse      string
		expectedProtoResponse proto.Message
		expectedError         bool
	}{
		{
			name:                "successful conversion from proto response to json marshalled response",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &banktypes.QueryAllBalancesResponse{},
			expectedProtoResponse: &banktypes.QueryAllBalancesResponse{
				Balances: sdk.NewCoins(sdk.NewCoin("bar", math.NewInt(30))),
				Pagination: &query.PageResponse{
					NextKey: []byte("foo"),
				},
			},
		},
		{
			name:                "invalid proto response struct",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &ammtypes.QueryGetPoolResponse{},
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTestInternal()

			originalVersionBz, err := hex.DecodeString(tc.originalResponse)
			suite.Require().NoError(err)

			jsonMarshalledResponse, err := types.ConvertProtoToJSONMarshal(tc.protoResponseStruct, originalVersionBz, suite.app.AppCodec())
			if tc.expectedError {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)

			// check response by json marshalling proto response into json response manually
			jsonMarshalExpectedResponse, err := suite.app.AppCodec().MarshalJSON(tc.expectedProtoResponse)
			suite.Require().NoError(err)
			suite.Require().Equal(jsonMarshalledResponse, jsonMarshalExpectedResponse)
		})
	}
}
