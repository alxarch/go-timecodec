package tc_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/alxarch/go-timecodec"
)

func Test_LayoutCodec(t *testing.T) {
	d := time.Now()
	layout := "2006-01-02"
	expect := d.Format(layout)
	c := tc.LayoutCodec(layout)
	if actual := c.MarshalTime(d); actual != expect {
		t.Error("Invalid date string %s", expect)
	}
	tm, err := c.UnmarshalTime(expect)
	if err != nil {
		t.Error(err)
	}
	sexpect, _ := time.Parse(layout, expect)
	if tm != sexpect {
		t.Error("Invalid time ", tm)
	}

}
func Test_Millis(t *testing.T) {
	now := time.Now()
	ms := tc.UnixMillis(now)
	dt := now.UnixNano() - ms*1000000
	if dt < 0 {
		dt = -dt
	}
	if dt > int64(time.Millisecond) {
		t.Error("Wrong millis", dt)
	}
}
func Test_MillisCodec(t *testing.T) {
	tnow := time.Now()
	now := tc.UnixMillis(tnow)
	snow := fmt.Sprintf("%d", now)
	tm, err := tc.MillisTimeCodec.UnmarshalTime(snow)
	if err != nil {
		t.Error(err)
	}

	if tc.UnixMillis(tm) != now {
		t.Error("Invalid time ", now-tc.UnixMillis(tm))
	}
	s := tc.MillisTimeCodec.MarshalTime(tnow)
	if s != snow {
		t.Error("Invalid decoder output ", s, snow)
	}
	_, err = tc.MillisTimeCodec.UnmarshalTime("foo")
	if err == nil {
		t.Error("Invalid error")
	}
}
func Test_NewTimeCodecNoDec(t *testing.T) {
	defer func() {
		if msg := recover(); msg == nil {
			t.Error("Didn't panic without encoder")
		}
	}()
	lc := tc.LayoutCodec("")
	tc.NewTimeCodec(nil, lc.UnmarshalTime)
}

func Test_NewTimeCodecNoEnc(t *testing.T) {
	defer func() {
		if msg := recover(); msg == nil {
			t.Error("Didn't panic without encoder")
		}
	}()
	lc := tc.LayoutCodec("")
	tc.NewTimeCodec(lc.MarshalTime, nil)
}

func Test_ISOWeekCodec(t *testing.T) {
	tm, err := tc.ISOWeekCodec.UnmarshalTime("2017-09")
	if err != nil {
		t.Error(err.Error())
	}
	if _, w := tm.ISOWeek(); w != 9 {
		t.Error(fmt.Sprintf("ISOWeek invalid week %d", w))
	}
	if y, _ := tm.ISOWeek(); y != 2017 {
		t.Error(fmt.Sprintf("ISOWeek invalid year %d", y))
	}
	s := tc.ISOWeekCodec.MarshalTime(tm)
	if s != "2017-09" {
		t.Error(fmt.Sprintf("ISOWeek invalid string %s", s))

	}

	if _, err := tc.ISOWeekCodec.UnmarshalTime("2017-09a"); err != tc.InvalidISOWeekString {
		t.Error("Invalid error")
	}
	if _, err := tc.ISOWeekCodec.UnmarshalTime("2017-99"); err != tc.InvalidWeekNumberError {
		t.Error("Invalid error")
	}
}
