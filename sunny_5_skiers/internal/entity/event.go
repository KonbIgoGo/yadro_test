package entity

import (
	"errors"
	"fmt"
	"time"
)

type Event struct {
	Timestamp       time.Time
	Kind            int64
	CompetitorID    int64
	AdditionalParam string
	Comment         string
}

func DisqualificationEvent(competitorID int64, timestamp time.Time) *Event {
	return &Event{
		Timestamp:    timestamp,
		Kind:         32,
		CompetitorID: competitorID,
		Comment:      fmt.Sprintf("The competitor(%d) is disqualified", competitorID),
	}
}

func FinishEvent(competitorID int64, timestamp time.Time) *Event {
	return &Event{
		Timestamp:    timestamp,
		Kind:         33,
		CompetitorID: competitorID,
		Comment:      fmt.Sprintf("The competitor(%d) is finished", competitorID),
	}
}

var (
	ErrUnexpectedKind = errors.New("unexpected event kind")
)
