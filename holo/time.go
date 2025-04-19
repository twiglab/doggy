package holo

import "time"

func MilliToTime(milli int64, tz int64) time.Time {
	return time.UnixMilli(milli)
}
