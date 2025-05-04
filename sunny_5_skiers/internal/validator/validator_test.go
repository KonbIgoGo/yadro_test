package validator

import (
	"biathlon/config"
	"biathlon/internal/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type testCase struct {
	input           string
	errExpected     bool
	errProcExpected bool
}

func TestValidator(t *testing.T) {
	t.Parallel()
	l, _ := zap.NewProduction()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	processor := mocks.NewMockProcessor(ctrl)
	validator := New(l, &config.Config{FiringLines: 2, Start: "10:00:00.000"}, processor)

	t.Run("parser test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				input:       "[09:31:49.285] 1 3",
				errExpected: false,
			},
			{
				input:       "[09:55:00.000] 2 1 10:00:00.000",
				errExpected: false,
			},
			{
				input:       "[incorrectTimestamp] 2 1 10:00:00.000",
				errExpected: true,
			},
			{
				input:       "[09:55:00.000] incorrect 1 10:00:00.000",
				errExpected: true,
			},
			{
				input:       "[09:55:00.000] 2 incorrect 10:00:00.000",
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			_, err := validator.parseEvent(tc.input)
			if tc.errExpected {
				require.Error(t, err, tc.input)
			} else {
				require.NoError(t, err, tc.input)
			}
		}
	})

	t.Run("validation test", func(t *testing.T) {
		tcs := []testCase{
			{
				input:           "[09:31:49.285] 1 1",
				errExpected:     false,
				errProcExpected: false,
			},
			{
				input:           "[incorrectTimestamp] 1 1",
				errExpected:     true,
				errProcExpected: false,
			},
			{
				input:           "[09:31:49.285] 1 1",
				errExpected:     false,
				errProcExpected: false,
			},
			{
				input:           "[09:55:00.000] 5 1 1",
				errExpected:     false,
				errProcExpected: false,
			},
			{
				input:           "[09:55:00.000] 5 1 notInt",
				errExpected:     true,
				errProcExpected: false,
			},
			{
				input:           "[09:55:00.000] 5 1 50000",
				errExpected:     true,
				errProcExpected: false,
			},
			{
				input:           "[09:31:49.285] 1 3",
				errExpected:     true,
				errProcExpected: true,
			},
		}

		for _, tc := range tcs {
			if tc.errProcExpected {
				processor.EXPECT().Process(gomock.Any()).Return(errors.New("error"))
			} else if !tc.errExpected {
				processor.EXPECT().Process(gomock.Any()).Return(nil)
			}

			err := validator.Validate(tc.input)
			if tc.errExpected {
				require.Error(t, err, tc.input)
			} else {
				require.NoError(t, err, tc.input)
			}
		}
	})
}
