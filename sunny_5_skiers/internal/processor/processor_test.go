package processor

import (
	"biathlon/config"
	"biathlon/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testCase struct {
	eventList   []*entity.Event
	errExpected bool
}

func TestProcessor(t *testing.T) {
	l, _ := zap.NewProduction()

	t.Parallel()
	t.Run("add competitor test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         1,
						CompetitorID: 1,
					},
				},
				errExpected: false,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Now(),
						Kind:         1,
						CompetitorID: 1,
					},
				},
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			proc := New(&config.Config{}, l)
			var err error
			for _, e := range tc.eventList {
				procErr := proc.Process(e)
				if procErr != nil {
					err = procErr
				}
			}
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("schedule start test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Now(),
						Kind:            2,
						CompetitorID:    1,
						AdditionalParam: "10:00:00.000",
					},
				},
				errExpected: false,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:       time.Now(),
						Kind:            2,
						CompetitorID:    0,
						AdditionalParam: "10:00:00.000",
					},
				},
				errExpected: true,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Now(),
						Kind:            2,
						CompetitorID:    0,
						AdditionalParam: "incorrect timestamp",
					},
				},
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			proc := New(&config.Config{}, l)
			var err error
			for _, e := range tc.eventList {
				procErr := proc.Process(e)
				if procErr != nil {
					err = procErr
				}
			}
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("start test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Time{},
						Kind:            2,
						CompetitorID:    1,
						AdditionalParam: "10:00:00.000",
					},
					{
						Timestamp:    time.Time{}.Add(time.Duration(time.Hour * 10)),
						Kind:         4,
						CompetitorID: 1,
					},
				},
				errExpected: false,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         4,
						CompetitorID: 0,
					},
				},
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			proc := New(&config.Config{StartDelta: "00:01:30"}, l)
			var err error
			for _, e := range tc.eventList {
				procErr := proc.Process(e)
				if procErr != nil {
					err = procErr
				}
			}
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("hit target test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Time{},
						Kind:            2,
						CompetitorID:    1,
						AdditionalParam: "10:00:00.000",
					},
					{
						Timestamp:    time.Time{}.Add(time.Hour*10).AddDate(-1, 0, 0),
						Kind:         4,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Time{}.Add(time.Duration(time.Hour*10)).AddDate(-1, 0, 0),
						Kind:         6,
						CompetitorID: 1,
					},
				},
				errExpected: false,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         6,
						CompetitorID: 0,
					},
				},
				errExpected: true,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Now(),
						Kind:         6,
						CompetitorID: 1,
					},
				},
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			proc := New(&config.Config{StartDelta: "00:01:30"}, l)
			var err error
			for _, e := range tc.eventList {
				procErr := proc.Process(e)
				if procErr != nil {
					err = procErr
				}
			}
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("penalty lap test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Time{},
						Kind:            2,
						CompetitorID:    1,
						AdditionalParam: "10:00:00.000",
					},
					{
						Timestamp:    time.Time{}.Add(time.Hour*10).AddDate(-1, 0, 0),
						Kind:         4,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Time{}.Add(time.Duration(time.Hour*10)).AddDate(-1, 0, 0),
						Kind:         8,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Time{}.Add(time.Duration(time.Hour*11)).AddDate(-1, 0, 0),
						Kind:         9,
						CompetitorID: 1,
					},
				},
				errExpected: false,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         8,
						CompetitorID: 0,
					},
				},
				errExpected: true,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         9,
						CompetitorID: 0,
					},
				},
				errExpected: true,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Now(),
						Kind:         8,
						CompetitorID: 1,
					},
				},
				errExpected: true,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Now(),
						Kind:         9,
						CompetitorID: 1,
					},
				},
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			proc := New(&config.Config{StartDelta: "00:01:30"}, l)
			var err error
			for _, e := range tc.eventList {
				procErr := proc.Process(e)
				if procErr != nil {
					err = procErr
				}
			}
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("finish test", func(t *testing.T) {
		t.Parallel()
		tcs := []testCase{
			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Time{},
						Kind:            2,
						CompetitorID:    1,
						AdditionalParam: "10:00:00.000",
					},
					{
						Timestamp:    time.Time{}.Add(time.Hour*10).AddDate(-1, 0, 0),
						Kind:         4,
						CompetitorID: 1,
					},
					{
						Timestamp:    time.Now(),
						Kind:         10,
						CompetitorID: 1,
					},
				},
				errExpected: false,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Now(),
						Kind:         10,
						CompetitorID: 0,
					},
				},
				errExpected: true,
			},

			{
				eventList: []*entity.Event{
					{
						Timestamp:    time.Time{},
						Kind:         1,
						CompetitorID: 1,
					},
					{
						Timestamp:       time.Time{},
						Kind:            2,
						CompetitorID:    1,
						AdditionalParam: "10:00:00.000",
					},
					{
						Timestamp:    time.Now(),
						Kind:         10,
						CompetitorID: 1,
					},
				},
				errExpected: true,
			},
		}

		for _, tc := range tcs {
			proc := New(&config.Config{StartDelta: "00:01:30"}, l)
			var err error
			for _, e := range tc.eventList {
				procErr := proc.Process(e)
				if procErr != nil {
					err = procErr
				}
			}
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})
}
