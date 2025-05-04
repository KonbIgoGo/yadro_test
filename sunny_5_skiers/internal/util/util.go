package util

import (
	"fmt"
	"time"
)

const layout = "15:04:05.000"

func ConvertToTimestamp(s string) (time.Time, error) {
	return time.Parse(layout, s)
}

func GetTimeDiffString(a, b time.Time) string {
	return FormatDuration(GetTimeDiff(a, b))
}

func GetTimeDiff(a, b time.Time) time.Duration {
	return a.Sub(b)
}

func GetAverageSpeed(d time.Duration, size int) float32 {
	return float32(d.Seconds()) / float32(size)
}

func FormatDuration(d time.Duration) string {

	if d < 0 {
		d = -d
	}
	hours := int64(d / time.Hour)
	d -= time.Duration(hours) * time.Hour
	minutes := int64(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute
	seconds := int64(d / time.Second)
	d -= time.Duration(seconds) * time.Second
	millis := int64(d / time.Millisecond)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
}
