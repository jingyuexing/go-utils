package utils

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

type FormatTemplate string
type AddUnits string

const (
	FormatYear        FormatTemplate = "YYYY"
	FormatShortYear   FormatTemplate = "YY"
	FormatMonth       FormatTemplate = "MM"
	FormatShortMonth  FormatTemplate = "M"
	FormatDay         FormatTemplate = "dd"
	FormatUpperDay    FormatTemplate = "DD"
	FormatShortDay    FormatTemplate = "d"
	FormatHour        FormatTemplate = "HH"
	FormatShortHour   FormatTemplate = "H"
	FormatSecond      FormatTemplate = "ss"
	FormatShortSecond FormatTemplate = "s"
	FormatMinute      FormatTemplate = "mm"
	FormatShortMinute FormatTemplate = "m"
	FormatMillisecond FormatTemplate = "ms"
	FormatWeek        FormatTemplate = "W"
	FormatShortWeek   FormatTemplate = "WW"
)

const (
	AddYearUnit           AddUnits = "year"
	AddShortUnit          AddUnits = "Y"
	AddDaysUnit           AddUnits = "day"
	AddDaysShortUnit      AddUnits = "D"
	AddDaysShortLowerUnit AddUnits = "d"
)

func DateTimeFormat(date time.Time, format string) string {
	formatted := format
	formatMap := map[FormatTemplate]int{
		FormatYear:        date.Year(),
		FormatShortYear:   date.Year() % 100,
		FormatMonth:       int(date.Month()),
		FormatShortMonth:  int(date.Month()),
		FormatDay:         date.Day(),
		FormatShortDay:    date.Day(),
		FormatUpperDay:    date.Day(),
		FormatHour:        date.Hour(),
		FormatShortHour:   date.Hour() % 12,
		FormatSecond:      date.Second(),
		FormatShortSecond: date.Second(),
		FormatMinute:      date.Minute(),
		FormatShortMinute: date.Minute(),
		FormatMillisecond: date.Nanosecond() / int(time.Millisecond),
		FormatWeek:        int(date.Weekday()),
		FormatShortWeek:   int(date.Weekday()),
	}

	for _, key := range []FormatTemplate{FormatYear, FormatMonth, FormatShortYear, FormatShortMonth, FormatDay, FormatUpperDay, FormatShortDay, FormatHour, FormatShortHour, FormatSecond, FormatShortSecond, FormatMinute, FormatShortMinute, FormatMillisecond, FormatWeek, FormatShortWeek} {
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

type FormatCallback func(int) string

type DateTime struct {
	time         time.Time
	Year         int
	Date         int
	Month        int
	Day          int
	Week         int
	Hour         int
	Minute       int
	Second       int
	Milliseconds int
	Nanosecond   int
	DateFormat   string
	TimeFormat   string
	monthFormat  FormatCallback
	weekFormat   FormatCallback
}

func NewDateTime() DateTime {
	return From(time.Now())
}

func From(time time.Time) DateTime {
    d := DateTime{
        time:       time,
        DateFormat: "YYYY-MM-DD",
        TimeFormat: "HH:mm:ss.ms",
    }
    d.Year = d.time.Year()
    d.Month = int(d.time.Month())
    d.Date = d.time.Day()
    d.Day = d.time.Day()
    d.Hour = d.time.Hour()
    d.Minute = d.time.Minute()
    d.Second = d.time.Second()
    d.Nanosecond = d.time.Nanosecond()
    return d
}

func (dt DateTime) SetYear(year int, month int, day int, hour int, minute int, second int, nanosecond int) DateTime {
	d := DateTime{
		time:        time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, time.UTC),
		DateFormat:  dt.DateFormat,
		TimeFormat:  dt.TimeFormat,
		monthFormat: dt.monthFormat,
		weekFormat:  dt.weekFormat,
	}
	d.Year = d.time.Year()
	d.Month = int(d.time.Month())
	d.Day = d.time.Day()
	d.Hour = d.time.Hour()
	d.Minute = d.time.Minute()
	d.Second = d.time.Second()
	d.Nanosecond = d.time.Nanosecond()
	return d
}

func (dt DateTime) SetMonth(month int, day int, hour int, minute int, second int, nanosecond int) DateTime {
	return dt.SetYear(dt.Year, month, day, hour, minute, second, nanosecond)
}

func (dt DateTime) SetDate(day int, hour int, minute int, second int, nanosecond int) DateTime {
	return dt.SetYear(dt.Year, dt.Month, day, hour, minute, second, nanosecond)
}

func (dt DateTime) SetHour(hour int, minute int, second int, nanoseconds int) DateTime {
	return dt.SetYear(dt.Year, dt.Month, dt.Day, hour, minute, second, nanoseconds)
}

func (dt DateTime) SetMinute(minute int, second int, nanoseconds int) DateTime {
	return dt.SetHour(dt.Hour, minute, second, nanoseconds)
}

func (dt DateTime) SetSecond(second int, nanoseconds int) DateTime {
	return dt.SetMinute(dt.Minute, second, nanoseconds)
}

func (dt DateTime) SetMilliseconds(milliseconds int, nanosecond int) DateTime {
	return dt.SetHour(dt.Hour, dt.Minute, 0, milliseconds*int(time.Millisecond)+nanosecond)
}

func (dt DateTime) SetNanosecond(nanoseconds int) DateTime {
	return dt.SetSecond(dt.Second, nanoseconds)
}

func (dt DateTime) Format(format string) string {
	return DateTimeFormat(dt.time, format)
}

func (dt DateTime) String() string {
	return dt.Format("YYYY-MM-DDTHH:mm:ss.msZ")
}

func (dt DateTime) SetWeekFormatFunc(format func(week int) string) {
	dt.weekFormat = format
}

func (dt DateTime) SetMonthFormatFunc(format func(month int) string) {
	dt.monthFormat = format
}

func (dt DateTime) MonthToString() string {
	if dt.weekFormat == nil {
		return dt.time.Month().String()
	}
	return dt.monthFormat(dt.Month)
}

func (dt DateTime) WeekToString() string {
	if dt.weekFormat == nil {
		return dt.time.Weekday().String()
	}
	return dt.weekFormat(int(dt.WeekDay()))
}

func (dt DateTime) TimeToString() string {
	return dt.Format(dt.TimeFormat)
}

func (dt DateTime) Today() string {
	return dt.Format(dt.DateFormat)
}

func (dt DateTime) SetTime(sec int64, ns int64) DateTime {
	d := DateTime{
		time:        time.Unix(sec, ns),
		DateFormat:  dt.DateFormat,
		TimeFormat:  dt.TimeFormat,
		monthFormat: dt.monthFormat,
		weekFormat:  dt.weekFormat,
	}
	d.Year = d.time.Year()
	d.Month = int(d.time.Month())
	d.Day = d.time.Day()
	d.Date = d.Day
	d.Hour = d.time.Hour()
	d.Minute = d.time.Minute()
	d.Second = d.time.Second()
	d.Nanosecond = d.time.Nanosecond()
	return d
}

func (dt DateTime) RawTime() *time.Time {
    return &dt.time
}

func (dt DateTime) LocaleCallBack(t string, call func(t string) string) string {
	return dt.Format(call(t))
}

func (dt DateTime) IsLeapYear() bool {
	return ((dt.Year%4 == 0) && (dt.Year%100 != 0)) || (dt.Year%400 == 0)
}

func (dt DateTime) DayOfYear() int {
	return dt.time.YearDay()
}

func (dt DateTime) MinuteOfDay() int {
	return dt.Hour*60 + dt.Minute
}

func (dt DateTime) CurrentYearDays() int {
	days := 365
	if dt.IsLeapYear() {
		return days + 1
	}
	return days
}

func (dt DateTime) AddDays(days int) DateTime {
	return dt.Add(days, "day")
}

func (dt DateTime) AddMonth(month int) DateTime {
	return dt.Add(month, "month")
}
func (dt DateTime) AddWeek(week int) DateTime {
	return dt.Add(week, "week")
}

func (d DateTime) SubtractDays(days int) DateTime {
	return d.AddDays(-days)
}

func (d DateTime) WeekDay() int {
	return int(d.time.Weekday())
}

func (d DateTime) CountDays() int {
    start := time.Date(d.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(d.Year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
    return int(end.Sub(start).Hours() / 24)
}

func (d DateTime) Progress() float32 {
    return float32(d.time.YearDay()) / float32(d.CountDays())
}

func (d DateTime) WeekOfYear() int {
	_, week := d.time.ISOWeek()
	return week
}

func (d DateTime) Add(num int, unit AddUnits) DateTime {
	result := DateTime{}
	switch unit {
	case "Year", "year", "Y":
		result = d.SetYear(d.Year+num, d.Month, d.Day, d.Hour, d.Minute, d.Second, d.Nanosecond)
	case "month", "M":
		result = d.SetMonth(d.Month+num, d.Day, d.Hour, d.Minute, d.Second, d.Nanosecond)
	case "day", "D":
		result = d.SetDate(d.Day+num, d.Hour, d.Minute, d.Second, d.Nanosecond)
	case "week", "W":
		result = d.SetDate(d.Day+(num*7), d.Hour, d.Minute, d.Second, d.Nanosecond)
	case "hour", "H", "h":
		result = d.SetHour(d.Hour+num, d.Minute, d.Second, d.Nanosecond)
	case "minute", "m":
		result = d.SetMinute(d.Minute+num, d.Second, d.Nanosecond)
	case "second", "s":
		result = d.SetSecond(d.Second+num, d.Nanosecond)
	case "milliseconds", "ms":
		result = d.SetMilliseconds(d.Milliseconds+num, d.Nanosecond)
	case "nanosecond", "ns":
		result = d.SetNanosecond(d.Nanosecond + num)
	}
	return result
}

func (dt DateTime) Time() int64 {
	return dt.time.UnixNano()
}

func (dt DateTime) Now() DateTime {
	now := time.Now()
	return dt.SetTime(now.Unix(), int64(now.Nanosecond()))
}

func (dt DateTime) UTCOffset() int64 {
	localUnix := dt.time.Local().Unix()
	utcUnix := dt.time.Unix()
	return localUnix - utcUnix
}

func (dt DateTime) IsToday(date DateTime) bool {
	return dt.Year == date.Year && dt.Month == date.Month && dt.Day == date.Day
}

func (dt DateTime) Max(date ...DateTime) DateTime {
	var max DateTime
	for _, v := range date {
		if max.Time() < v.Time() {
			max = v
		}
	}
	return max
}

func (dt DateTime) Min(date ...DateTime) DateTime {
	min := date[0]
	for _, v := range date {
		if v.Time() > min.Time() {
			min = v
		}
	}
	return min
}

// the current time is before the specified time
func (dt DateTime) IsBefore(d DateTime) bool {
	return d.Time() > dt.Time()
}

// the current time is after the specified time
func (dt DateTime) IsAfter(d DateTime) bool {
	return dt.Time() < d.Time()
}

func parseFormatTemplate(formatTemplate string) []string {
	formatParts := make([]string, 0)
	i := 0
	for i < len(formatTemplate) {
		switch {
		case strings.HasPrefix(formatTemplate[i:], string(FormatYear)):
			formatParts = append(formatParts, string(FormatYear))
			i += 4
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortDay)):
			formatParts = append(formatParts, string(FormatShortYear))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatMonth)):
			formatParts = append(formatParts, string(FormatMonth))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortMonth)):
			formatParts = append(formatParts, string(FormatShortMonth))
			i++
		case strings.HasPrefix(formatTemplate[i:], string(FormatDay)) || strings.HasPrefix(formatTemplate[i:], string(FormatUpperDay)):
			formatParts = append(formatParts, string(FormatDay))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortDay)):
			formatParts = append(formatParts, string(FormatShortDay))
			i++
		case strings.HasPrefix(formatTemplate[i:], string(FormatHour)):
			formatParts = append(formatParts, string(FormatHour))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortHour)):
			formatParts = append(formatParts, string(FormatShortHour))
			i++
		case strings.HasPrefix(formatTemplate[i:], string(FormatMinute)):
			formatParts = append(formatParts, string(FormatMinute))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortHour)):
			formatParts = append(formatParts, string(FormatShortHour))
			i++
		case strings.HasPrefix(formatTemplate[i:], string(FormatSecond)):
			formatParts = append(formatParts, string(FormatSecond))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortSecond)):
			formatParts = append(formatParts, string(FormatShortSecond))
			i++
		case strings.HasPrefix(formatTemplate[i:], string(FormatMillisecond)):
			formatParts = append(formatParts, string(FormatMillisecond))
			i += 2
		case strings.HasPrefix(formatTemplate[i:], string(FormatWeek)):
			formatParts = append(formatParts, string(FormatWeek))
			i++
		case strings.HasPrefix(formatTemplate[i:], string(FormatShortWeek)):
			formatParts = append(formatParts, string(FormatShortWeek))
			i += 2
		default:
			i++
		}
	}
	return formatParts
}
func extractDateParts(date string, formatParts []string) map[string]string {
	dateParts := make(map[string]string)
	current := 0

	for _, part := range formatParts {
		val := ""
		for current < len(date) && (unicode.IsDigit(rune(date[current])) || rune(date[current]) == '.') {
			val += string(date[current])
			current++
		}
		dateParts[part] = val
		current++
	}

	return dateParts
}

func (dt DateTime) Parse(date string, formatTemplate string) *DateTime {
	datetime := NewDateTime()
    if formatTemplate == "" {
        formatTemplate = dt.DateFormat + "T" + dt.TimeFormat + "Z"
    }
	formatParts := parseFormatTemplate(formatTemplate)
	dateParts := extractDateParts(date, formatParts)

	if len(dateParts) == 0 {
		return nil
	}

	for key, value := range dateParts {
		val, _ := strconv.Atoi(value)
		switch key {
		case "YYYY":
			datetime.Year = val
		case "YY":
			datetime.Year = val
		case "MM", "M":
			datetime.Month = val
		case "dd", "DD", "d":
			datetime.Day = val
		case "HH", "H":
			datetime.Hour = val
		case "W", "WW":
			datetime.Week = val
		case "mm", "m":
			datetime.Minute = val
		case "ss", "s":
			datetime.Second = val
		case "ms":
			datetime.Milliseconds = val
		}
	}

	return &datetime
}
