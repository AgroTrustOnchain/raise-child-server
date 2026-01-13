package util

import "time"

const dateFormat string = "25/12/2026"

var request_duration time.Duration = time.Hour * 72 // 3 days

func ToMilliseconds(time time.Time) int64 {
	return time.UnixMilli()
}

func RawDateToTime(rawDate string) time.Time {
	parsedTime, _ := time.Parse(dateFormat, rawDate)
	return parsedTime
}

func MilliSecToTime(milliseconds int64) time.Time {
	return time.UnixMilli(milliseconds)
}

func GetRequestDuration() time.Time {
	return time.Now().Add(request_duration)
}
