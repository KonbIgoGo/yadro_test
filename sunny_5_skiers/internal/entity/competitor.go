package entity

import (
	"biathlon/internal/util"
	"errors"
	"fmt"
	"strings"
	"time"
)

type LapData struct {
	StartLap  time.Time
	FinishLap time.Time
	Size      int
}

func (l *LapData) toString() string {
	dur := l.FinishLap.Sub(l.StartLap)
	return fmt.Sprintf("{%s, %.3f}", util.FormatDuration(dur), util.GetAverageSpeed(dur, l.Size))
}

type Competitior struct {
	ID                 int64
	Penalty            int
	TotalTargets       int
	HitedTargets       int
	LapCounter         int
	MainLapsData       []LapData
	PenaltyLapData     []LapData
	FinishRaceTime     time.Time
	ScheduledStartTime time.Time
	Status             string
}

func (c *Competitior) getTotalTime() string {
	if c.Status != "Finished" {
		return fmt.Sprintf("[%s]", c.Status)
	}
	return util.GetTimeDiffString(c.FinishRaceTime, c.ScheduledStartTime)
}

func (c *Competitior) getLapsData() string {
	var res []string = make([]string, len(c.MainLapsData))
	for i, l := range c.MainLapsData {
		res[i] = l.toString()
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ","))
}

func (c *Competitior) getPenaltySumDuration() time.Duration {
	var total time.Duration
	for _, l := range c.PenaltyLapData {
		total += l.FinishLap.Sub(l.StartLap)
	}
	return total
}

func (c *Competitior) GetResult() string {
	sum := c.getPenaltySumDuration()
	return fmt.Sprintf("%s %d %s {%s, %.3f} %d/%d",
		c.getTotalTime(),
		c.ID,
		c.getLapsData(),
		util.FormatDuration(sum),
		util.GetAverageSpeed(sum, c.Penalty),
		c.HitedTargets,
		c.TotalTargets)
}

var (
	ErrCompetitorNotFound     = errors.New("competitor not found")
	ErrCompetitorAlreadyExist = errors.New("competitor already exist")
	ErrCompetitorDisqualified = errors.New("competitior disqualified")
)
