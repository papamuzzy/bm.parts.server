package xtime

import (
	"time"
)

func FromStr(str string) time.Time {
	//location, _ := time.LoadLocation("Europe/Berlin")
	location, _ := time.LoadLocation("Local")
	ret, _ := time.ParseInLocation(time.DateTime, str, location)

	return ret
}

func ToStr(t time.Time) string {
	return t.Format(time.DateTime)
}
