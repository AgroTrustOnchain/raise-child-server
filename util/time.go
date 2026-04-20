package util

import "time"

const dateFormat string = "02/01/2006"

var request_duration time.Duration = time.Hour * 72 // 3 days
var bank_request_duration time.Duration = time.Hour * 24 * 7

func ToMilliseconds(time time.Time) int64 {
	return time.UnixMilli()
}

func RawDateToTime(rawDate string) time.Time {
	parsedTime, _ := time.Parse(dateFormat, rawDate)
	return parsedTime
}

func TimeToRawDate(date time.Time) string {
	return date.Format(dateFormat)
}

func MilliSecToTime(milliseconds int64) time.Time {
	return time.UnixMilli(milliseconds)
}

func GetRequestDuration() time.Time {
	return time.Now().Add(time.Minute)
}

func GetBankTransactionDuration() time.Time {
	return time.Now().Add(time.Minute)
}

func ToStartOfDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 1, 0, 0, date.Location())
}

func ToEndOfDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 99, 0, date.Location())
}
