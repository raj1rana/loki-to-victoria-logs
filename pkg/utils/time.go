package utils

import (
    "time"
)

func GetTimeRange(window time.Duration) (time.Time, time.Time) {
    end := time.Now()
    start := end.Add(-window)
    return start, end
}

func ParseTimeString(timeStr string) (time.Time, error) {
    return time.Parse("01/02/2006 15:04:05", timeStr)
}
