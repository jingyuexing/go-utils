package utils

import (
	"strings"
	"time"
)

type FormatTemplate string

const (
	FormatYear        FormatTemplate = "YYYY"
	FormatShortYear   FormatTemplate = "YY"
	FormatMonth       FormatTemplate = "MM"
	FormatShortMonth  FormatTemplate = "M"
	FormatDay         FormatTemplate = "dd"
	FormatHour        FormatTemplate = "HH"
	FormatShortHour   FormatTemplate = "H"
	FormatSecond      FormatTemplate = "ss"
	FormatMinute      FormatTemplate = "mm"
	FormatMillisecond FormatTemplate = "ms"
	FormatWeek        FormatTemplate = "W"
	FormatShortWeek   FormatTemplate = "WW"
)

func DateTimeFormat(date time.Time, format string) string {
	formatted := format
	formatMap := map[FormatTemplate]int{
		FormatYear:        date.Year(),
		FormatShortYear:   date.Year() % 100,
		FormatMonth:       int(date.Month()),
		FormatShortMonth:  int(date.Month()),
		FormatDay:         date.Day(),
		FormatHour:        date.Hour(),
		FormatShortHour:   date.Hour() % 12,
		FormatSecond:      date.Second(),
		FormatMinute:      date.Minute(),
		FormatMillisecond: date.Nanosecond() / int(time.Millisecond),
		FormatWeek:        int(date.Weekday()),
		FormatShortWeek:   int(date.Weekday()),
	}

	for _, key := range []FormatTemplate{FormatYear, FormatMonth, FormatShortYear, FormatShortMonth, FormatDay, FormatHour, FormatShortHour, FormatSecond, FormatMinute, FormatMillisecond, FormatWeek, FormatShortWeek} {
		formatted = strings.Replace(
			formatted,
			string(key),
			strings.Join(
				PadStart(
					strings.Split(ToString(formatMap[key]), ""),
					len(string(key)),
					"0",
				),
				"",
			),
			1,
		)
	}

	return formatted
}
