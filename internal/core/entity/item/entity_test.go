package item

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdrpkv/kvstore/internal/core/entity/item/mocks"
)

func Test_Key_Validate(t *testing.T) {
	var (
		stubKey               = "stubkey"
		stubRunesNotInStubKey = "\t"
	)

	type (
		testCaseGiveCallsKeyMaxLen struct {
			Returns int
			Called  bool
		}
		testCaseGiveCallsForbiddenRunes struct {
			Returns string
			Called  bool
		}
		testCaseGiveCalls struct {
			keyMaxLen      testCaseGiveCallsKeyMaxLen
			forbiddenRunes testCaseGiveCallsForbiddenRunes
		}
		testCaseGive struct {
			key   string
			calls testCaseGiveCalls
		}
		testCase struct {
			name string
			give testCaseGive
			want error
		}
	)
	for _, tc := range []testCase{
		{
			"success",
			testCaseGive{
				key: stubKey,
				calls: testCaseGiveCalls{
					keyMaxLen: testCaseGiveCallsKeyMaxLen{
						Returns: len(stubKey),
						Called:  true,
					},
					forbiddenRunes: testCaseGiveCallsForbiddenRunes{
						Returns: stubRunesNotInStubKey,
						Called:  true,
					},
				},
			},
			nil,
		},
		{
			"key max length error",
			testCaseGive{
				key: stubKey,
				calls: testCaseGiveCalls{
					keyMaxLen: testCaseGiveCallsKeyMaxLen{
						Returns: len(stubKey) - 1,
						Called:  true,
					},
				},
			},
			ErrKeyLength{
				Max: len(stubKey) - 1,
				Got: len(stubKey),
			},
		},
		{
			"forbidden runes error",
			testCaseGive{
				key: stubKey + stubRunesNotInStubKey,
				calls: testCaseGiveCalls{
					keyMaxLen: testCaseGiveCallsKeyMaxLen{
						Returns: len(stubKey) + len(stubRunesNotInStubKey),
						Called:  true,
					},
					forbiddenRunes: testCaseGiveCallsForbiddenRunes{
						Returns: stubRunesNotInStubKey,
						Called:  true,
					},
				},
			},
			ErrKeyForbiddenChar{
				Rune:  []rune(stubRunesNotInStubKey)[0],
				Index: len(stubKey),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			keyValidatorMock := mocks.NewKeyValidator(t)

			if tc.give.calls.keyMaxLen.Called {
				keyValidatorMock.On("KeyMaxLen").Return(
					tc.give.calls.keyMaxLen.Returns,
				)
			}

			if tc.give.calls.forbiddenRunes.Called {
				keyValidatorMock.On("ForbiddenRunes").Return(
					tc.give.calls.forbiddenRunes.Returns,
				)
			}

			gotErr := Key(tc.give.key).Validate(keyValidatorMock)

			if tc.want == nil {
				require.NoError(t, gotErr)
			}

			if tc.want != nil {
				require.ErrorIs(t, gotErr, tc.want)
			}
		})
	}
}
