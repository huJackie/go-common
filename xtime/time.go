package xtime

import (
	"time"
)

const _format = "2006-01-02"

type Time time.Time

func (t Time) format() string {
	f := time.Time(t)
	if f.IsZero() {
		return ""
	}
	return f.Format(_format)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.format() + `"`), nil
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.format()), nil
}

func (t *Time) UnmarshalJSON(text []byte) (err error) {
	now, err := time.ParseInLocation(`"`+_format+`"`, string(text), time.Local)
	*t = Time(now)
	return
}

func (t *Time) UnmarshalText(text []byte) (err error) {
	now, err := time.ParseInLocation(_format, string(text), time.Local)
	*t = Time(now)
	return
}

func (t Time) String() string {
	return time.Time(t).String()
}

const _jsonFormat = "2006-01-02 15:04:05"

type JsonTime time.Time

func (t JsonTime) format() string {
	f := time.Time(t)
	if f.IsZero() {
		return ""
	}
	return f.Format(_jsonFormat)
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.format() + `"`), nil
}

func (t JsonTime) MarshalText() ([]byte, error) {
	return []byte(t.format()), nil
}

func (t *JsonTime) UnmarshalJSON(text []byte) (err error) {
	now, err := time.ParseInLocation(`"`+_jsonFormat+`"`, string(text), time.Local)
	*t = JsonTime(now)
	return
}

func (t *JsonTime) UnmarshalText(text []byte) (err error) {
	now, err := time.ParseInLocation(_jsonFormat, string(text), time.Local)
	*t = JsonTime(now)
	return
}

func (t JsonTime) String() string {
	return time.Time(t).String()
}
