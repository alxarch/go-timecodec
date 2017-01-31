package tc

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type TimeEncoder interface {
	MarshalTime(time.Time) string
}
type TimeEncoderFunc func(time.Time) string
type TimeDecoderFunc func(string) (time.Time, error)
type TimeDecoder interface {
	UnmarshalTime(string) (time.Time, error)
}
type TimeCodec interface {
	TimeEncoder
	TimeDecoder
}

func NewTimeCodec(enc TimeEncoderFunc, dec TimeDecoderFunc) TimeCodec {
	if nil == enc {
		panic("Invalid TimeEncoder")
	}
	if nil == dec {
		panic("Invalid TimeDecoder")
	}
	return &timeCodecFunc{enc, dec}
}

type timeCodecFunc struct {
	enc TimeEncoderFunc
	dec TimeDecoderFunc
}

func (c *timeCodecFunc) MarshalTime(t time.Time) string {
	return c.enc(t)
}
func (c *timeCodecFunc) UnmarshalTime(value string) (time.Time, error) {
	return c.dec(value)
}

type LayoutCodec string

func (layout LayoutCodec) UnmarshalTime(value string) (time.Time, error) {
	return time.Parse(string(layout), value)
}

func (layout LayoutCodec) MarshalTime(t time.Time) string {
	return t.Format(string(layout))
}

var isoweekRx = regexp.MustCompile("^(\\d{4})-(\\d{2})$")
var (
	InvalidISOWeekString   = errors.New("Invalid ISOWeek string")
	InvalidWeekNumberError = errors.New("Invalid week number")
)

func UnixMillis(tm time.Time) int64 {
	return int64(time.Duration(tm.UnixNano()) * time.Nanosecond / time.Millisecond)
}

var ISOWeekCodec = NewTimeCodec(func(t time.Time) string {
	y, d := t.ISOWeek()
	return fmt.Sprintf("%d-%02d", y, d)

}, func(value string) (time.Time, error) {
	match := isoweekRx.FindStringSubmatch(value)
	if match == nil {
		return time.Time{}, InvalidISOWeekString
	}
	year, _ := strconv.Atoi(string(match[1]))
	week, _ := strconv.Atoi(string(match[2]))
	if !(0 < week && week <= 53) {
		return time.Time{}, InvalidWeekNumberError
	}
	t := time.Date(year, 1, 0, 0, 0, 0, 0, time.UTC)
	for t.Weekday() > time.Sunday {
		t = t.Add(-24 * time.Hour)
	}
	t = t.Add(time.Duration(week+1) * 7 * 24 * time.Hour)
	return t, nil
})

var MillisTimeCodec = NewTimeCodec(func(t time.Time) string {
	return strconv.FormatInt(UnixMillis(t), 10)
}, func(value string) (time.Time, error) {
	ms, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(0, int64(time.Duration(ms)*time.Millisecond))

	return t, err
})
