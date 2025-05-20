package keeper_test

import (
	"github.com/elys-network/elys/v4/x/stablestake/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestParams() {
	tests := []struct {
		name          string
		req           *types.QueryParamsRequest
		expectedError error
		expectedResp  *types.QueryParamsResponse
	}{
		{
			name:          "valid request",
			req:           &types.QueryParamsRequest{},
			expectedError: nil,
			expectedResp:  &types.QueryParamsResponse{Params: types.DefaultParams()},
		},
		{
			name:          "invalid request",
			req:           nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid request"),
			expectedResp:  nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			resp, err := suite.app.StablestakeKeeper.Params(suite.ctx, tt.req)
			if tt.expectedError != nil {
				require.ErrorIs(suite.T(), err, tt.expectedError)
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResp, resp)
			}
		})
	}
}
