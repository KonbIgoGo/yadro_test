package validator

import (
	"biathlon/internal/entity"
	"biathlon/internal/util"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func (i *implementation) Validate(rawData string) error {
	event, err := i.parseEvent(rawData)
	if err != nil {
		return err
	}

	switch event.Kind {
	case 1:
		startTime, err := util.ConvertToTimestamp(i.cfg.Start)
		if err != nil {
			return errors.New("incorrect event timestamp")
		}
		event.Comment = fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)

		if startTime.Before(event.Timestamp) {
			event = entity.DisqualificationEvent(event.CompetitorID, event.Timestamp)
		}
	case 2:
		event.Comment =
			fmt.Sprintf(
				"The start time for the competitor(%d) was set by a draw to %s",
				event.CompetitorID,
				event.AdditionalParam)
	case 3:
		event.Comment = fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
	case 4:
		event.Comment = fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
	case 5:

		flNumber, err := strconv.ParseInt(event.AdditionalParam, 10, 64)
		if err != nil {
			return errors.New("incorrect firing range format")
		}

		if flNumber > int64(i.cfg.FiringLines) {
			return errors.New("number of fire line is more then the amount of firelines")
		}

		event.Comment = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.AdditionalParam)
	case 6:
		event.Comment = fmt.Sprintf("The target(%s) has been hit by competitior(%d)", event.AdditionalParam, event.CompetitorID)
	case 7:
		event.Comment = fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
	case 8:
		event.Comment = fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
	case 9:
		event.Comment = fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
	case 10:
		event.Comment = fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
	case 11:
		event.Comment = fmt.Sprintf("The competitor(%d) can`t continue: %s", event.CompetitorID, event.AdditionalParam)
	default:
		return entity.ErrUnexpectedKind
	}

	return i.processor.Process(event)
}

func (i *implementation) GetLog() {
	println("log=============================")
	for _, s := range i.processor.GetLog() {
		fmt.Println(s)
	}
	println("log=============================")
}

func (i *implementation) GetResult() {
	println("result table====================")
	for _, s := range i.processor.GetResult() {
		fmt.Println(s)
	}
	println("result table====================")
}

func (i *implementation) parseEvent(rawData string) (*entity.Event, error) {
	var res = &entity.Event{}
	splitedData := strings.Fields(rawData)

	if len(splitedData) < 3 {
		return nil, errors.New("incorrect data format")
	}

	splitedData[0] = splitedData[0][1 : len(splitedData[0])-1]
	t, err := util.ConvertToTimestamp(splitedData[0])
	if err != nil {
		i.logger.Error("failed to convert timestamp", zap.Error(err))
		return nil, err
	}

	res.Timestamp = t

	res.Kind, err = strconv.ParseInt(splitedData[1], 10, 64)
	if err != nil {
		i.logger.Error("failed to parse int", zap.Error(err))
		return nil, err
	}

	res.CompetitorID, err = strconv.ParseInt(splitedData[2], 10, 64)
	if err != nil {
		i.logger.Error("failed to parse int", zap.Error(err))
		return nil, err
	}

	if len(splitedData) > 3 {
		res.AdditionalParam = splitedData[3]
	}

	return res, nil
}
