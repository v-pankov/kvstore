package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/delete/mocks"
)

func Test_Processor_Process(t *testing.T) {
	var (
		stubCtx = context.Background()
		stubErr = errors.New("stub error")
	)

	type (
		testCaseGiveCallsRepositoryDeleteItemByKey struct {
			returns error
			called  bool
		}
		testCaseGiveCallsRepository struct {
			deleteItemByKey testCaseGiveCallsRepositoryDeleteItemByKey
		}
		testCaseGiveCalls struct {
			Repository testCaseGiveCallsRepository
		}
		testCaseGive struct {
			request *Request
			calls   testCaseGiveCalls
		}
		testCaseWant struct {
			response *Response
			err      error
		}
		testCase struct {
			name string
			give testCaseGive
			want testCaseWant
		}
	)
	for _, tc := range []testCase{
		{
			"success",
			testCaseGive{
				request: &Request{},
				calls: testCaseGiveCalls{
					Repository: testCaseGiveCallsRepository{
						deleteItemByKey: testCaseGiveCallsRepositoryDeleteItemByKey{
							returns: nil,
							called:  true,
						},
					},
				},
			},
			testCaseWant{
				response: &Response{},
			},
		},
		{
			"error",
			testCaseGive{
				request: &Request{},
				calls: testCaseGiveCalls{
					Repository: testCaseGiveCallsRepository{
						deleteItemByKey: testCaseGiveCallsRepositoryDeleteItemByKey{
							returns: stubErr,
							called:  true,
						},
					},
				},
			},
			testCaseWant{
				response: nil,
				err:      stubErr,
			},
		},
	} {
		repositoryMock := mocks.NewRepository(t)

		if tc.give.calls.Repository.deleteItemByKey.called {
			repositoryMock.On("DeleteItemByKey", stubCtx, item.Key(tc.give.request.Key)).Return(
				tc.give.calls.Repository.deleteItemByKey.returns,
			)
		}

		processor := Processor{
			Gateways: Gateways{
				Repository: repositoryMock,
			},
		}

		gotResponse, gotErr := processor.Process(stubCtx, tc.give.request)

		if tc.want.err == nil {
			require.NoError(t, gotErr)
			require.Equal(t, tc.want.response, gotResponse)
		}

		if tc.want.err != nil {
			require.ErrorIs(t, gotErr, tc.want.err)
			require.Nil(t, gotResponse)
		}
	}
}
