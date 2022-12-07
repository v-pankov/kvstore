package delete

/*
func Test_Processor_Process(t *testing.T) {
	var (
		stubCtx = context.Background()
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
	for _, tc := range []testCase{} {
		repositoryMock := mocks.NewUseCaseItemDeleteRepository(t)

		if tc.give.calls.Repository.deleteItemByKey.called {
			repositoryMock.On("DeleteItemByKey").Return(
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
*/
