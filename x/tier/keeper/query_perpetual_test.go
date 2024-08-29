package keeper_test

import (
	"reflect"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tier/types"
	"github.com/golang/mock/gomock"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestPerpetualQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeeper := NewMockKeeper(ctrl) // Use gomock to create the mock
	ctx := sdk.Context{}
	wctx := sdk.WrapSDKContext(ctx)
	//keeper, ctx := keepertest.MembershiptierKeeper(t)

	mockKeeper.EXPECT().RetrievePerpetualTotal(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx sdk.Context, user string) (sdk.Dec, sdk.Dec, sdk.Dec) {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
		},
	).AnyTimes()

	//keeper.RetrievePerpetualTotal = mockKeeper.RetrievePerpetualTotal

	tests := []struct {
		desc     string
		request  *types.QueryPerpetualRequest
		response *types.QueryPerpetualResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryPerpetualRequest{
				User: "creator",
			},
			response: &types.QueryPerpetualResponse{TotalValue: sdk.ZeroDec(), TotalBorrows: sdk.ZeroDec()},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := mockKeeper.Perpetual(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

// Define the mock interface in the same file or import it if already defined
type MockKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockKeeperMockRecorder
}

func NewMockKeeper(ctrl *gomock.Controller) *MockKeeper {
	mock := &MockKeeper{ctrl: ctrl}
	mock.recorder = &MockKeeperMockRecorder{mock}
	return mock
}

type MockKeeperMockRecorder struct {
	mock *MockKeeper
}

// Implement the methods you need to mock
func (m *MockKeeper) EXPECT() *MockKeeperMockRecorder {
	return m.recorder
}

func (m *MockKeeper) RetrievePerpetualTotal(ctx sdk.Context, user string) (sdk.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrievePerpetualTotal", ctx, user)
	ret0, _ := ret[0].(sdk.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockKeeperMockRecorder) RetrievePerpetualTotal(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrievePerpetualTotal", reflect.TypeOf((*MockKeeper)(nil).RetrievePerpetualTotal), ctx, user)
}
