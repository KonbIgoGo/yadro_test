package processor

import "biathlon/internal/entity"

//go:generate mockgen -source=./interface.go -destination=../mocks/proc_mock.go -package=mocks
type Processor interface {
	Process(event *entity.Event) error
	GetLog() []string
	GetResult() []string
}
