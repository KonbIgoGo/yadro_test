package processor

import (
	"biathlon/config"
	"biathlon/internal/entity"
	"biathlon/internal/util"
	"slices"
	"time"

	"go.uber.org/zap"
)

var _ Processor = (*processorImpl)(nil)

type processorImpl struct {
	competitorList map[int64]*entity.Competitior
	events         []*entity.Event
	cfg            *config.Config
	logger         *zap.Logger
}

func (p *processorImpl) parseStartDelta() (time.Duration, error) {
	t, err := time.Parse("15:04:05", p.cfg.StartDelta)
	if err != nil {
		return 0, err
	}

	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond()), nil
}

func New(cfg *config.Config, logger *zap.Logger) *processorImpl {
	return &processorImpl{
		competitorList: make(map[int64]*entity.Competitior),
		events:         make([]*entity.Event, 0),
		cfg:            cfg,
		logger:         logger,
	}
}

func (p *processorImpl) Process(event *entity.Event) error {
	switch event.Kind {
	case 1:
		_, ok := p.competitorList[event.CompetitorID]
		if ok {
			err := entity.ErrCompetitorAlreadyExist
			p.logger.Error("failed to register competitor", zap.Error(err))
			return err
		}
		p.competitorList[event.CompetitorID] = &entity.Competitior{
			ID:             event.CompetitorID,
			Status:         "NotStarted",
			TotalTargets:   p.cfg.Laps * 5,
			HitedTargets:   0,
			LapCounter:     0,
			MainLapsData:   make([]entity.LapData, 0),
			PenaltyLapData: make([]entity.LapData, 0),
			Penalty:        p.cfg.Laps * p.cfg.PenaltyLen * p.cfg.FiringLines * 5,
		}
		p.events = append(p.events, event)
	case 2:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		time, err := util.ConvertToTimestamp(event.AdditionalParam)
		if err != nil {
			p.logger.Error("failed to convert additional param to timestamp", zap.Error(err))
			return err
		}

		competitor.ScheduledStartTime = time
		p.events = append(p.events, event)
	case 3:
		p.events = append(p.events, event)
	case 4:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		timeDelta, err := p.parseStartDelta()
		if err != nil {
			p.logger.Error("failed to convert start delta to duration", zap.Error(err))
			return err
		}

		delta := event.Timestamp.Sub(competitor.ScheduledStartTime)
		if delta < 0 || delta > timeDelta {
			competitor.Status = "NotStarted"
			p.events = append(p.events, entity.DisqualificationEvent(competitor.ID, event.Timestamp))
			return nil
		}

		competitor.Status = "Started"
		competitor.MainLapsData = append(competitor.MainLapsData, entity.LapData{StartLap: event.Timestamp, Size: p.cfg.LapLen})
		p.events = append(p.events, event)

	case 5:
		p.events = append(p.events, event)
	case 6:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		if competitor.Status != "Started" {
			return entity.ErrCompetitorDisqualified
		}

		competitor.Penalty -= p.cfg.PenaltyLen
		competitor.HitedTargets += 1
		p.events = append(p.events, event)
	case 7:
		p.events = append(p.events, event)
	case 8:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		if competitor.Status != "Started" {
			return entity.ErrCompetitorDisqualified
		}

		competitor.PenaltyLapData = append(competitor.PenaltyLapData, entity.LapData{StartLap: event.Timestamp, Size: p.cfg.PenaltyLen})
		p.events = append(p.events, event)
	case 9:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		if competitor.Status != "Started" {
			return entity.ErrCompetitorDisqualified
		}

		competitor.PenaltyLapData[len(competitor.PenaltyLapData)-1].FinishLap = event.Timestamp
		p.events = append(p.events, event)
	case 10:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		if competitor.Status != "Started" {
			return entity.ErrCompetitorDisqualified
		}

		if len(competitor.PenaltyLapData) != 0 &&
			competitor.PenaltyLapData[len(competitor.PenaltyLapData)-1].FinishLap.IsZero() {
			competitor.Status = "NotFinished"
			p.events = append(p.events, entity.DisqualificationEvent(competitor.ID, event.Timestamp))
			return nil
		}

		competitor.MainLapsData[competitor.LapCounter].FinishLap = event.Timestamp
		competitor.LapCounter += 1

		p.events = append(p.events, event)
		if competitor.LapCounter == p.cfg.Laps && competitor.Status == "Started" {
			competitor.Status = "Finished"
			competitor.FinishRaceTime = event.Timestamp
			p.events = append(p.events, entity.FinishEvent(competitor.ID, event.Timestamp))
		} else {
			competitor.MainLapsData = append(competitor.MainLapsData, entity.LapData{StartLap: event.Timestamp, Size: p.cfg.LapLen})
		}
	case 11:
		competitor, ok := p.competitorList[event.CompetitorID]
		if !ok {
			err := entity.ErrCompetitorNotFound
			p.logger.Error("failed to get competitor", zap.Error(err))
			return err
		}

		competitor.Status = "NotFinished"
		p.events = append(p.events, event)
		p.events = append(p.events, entity.DisqualificationEvent(competitor.ID, event.Timestamp))
	default:
		return entity.ErrUnexpectedKind
	}
	return nil
}

func (p *processorImpl) GetResult() []string {
	var finished []*entity.Competitior = make([]*entity.Competitior, 0)
	var disqualified []*entity.Competitior = make([]*entity.Competitior, 0)

	for _, c := range p.competitorList {
		if c.Status != "Finished" {
			disqualified = append(disqualified, c)
		} else {
			finished = append(finished, c)
		}
	}

	slices.SortStableFunc(finished, func(i, j *entity.Competitior) int {
		ti := i.FinishRaceTime.Sub(i.ScheduledStartTime)
		tj := j.FinishRaceTime.Sub(j.ScheduledStartTime)

		if ti == tj {
			return 0
		} else if ti > tj {
			return 1
		}
		return -1
	})

	var res []string = make([]string, 0, len(p.competitorList))

	for _, c := range finished {
		res = append(res, c.GetResult())
	}

	for _, c := range disqualified {
		res = append(res, c.GetResult())
	}

	return res
}

func (p *processorImpl) GetLog() []string {
	var res []string = make([]string, len(p.events))
	for i, e := range p.events {
		res[i] = e.Comment
	}

	return res
}
