package utils

import (
	"fmt"
	"time"

	"github.com/snabb/isoweek"
)

func IdOf(i int) string {
	var b string
	if i >= 26 {
		b = IdOf((i / 26 >> 0) - 1)
	} else {
		b = ""
	}
	return b + string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i%26>>0])
}

func GetWeekFileName(t time.Time) string {
	wyear, week := isoweek.FromDate(t.Year(), t.Month(), t.Day())
	return fmt.Sprintf("%dw%d", wyear, week)
}

func GetWeekDay(t time.Time, loc *time.Location) []string {
	wyear, week := isoweek.FromDate(t.Year(), t.Month(), t.Day())
	c := isoweek.StartTime(wyear, week, loc)
	buf := make([]string, 7)
	for i := 0; i < 7; i++ {
		buf[i] = c.AddDate(0, 0, -1+i).Format("01/02, Monday")
	}
	return buf
}

func GetTimeRange() []string {
	buf := []string{
		"5",
		"5:30",
		"6",
		"6:30",
		"7",
		"7:30",
		"8",
		"8:30",
		"9",
		"9:30",
		"10",
		"10:30",
		"11",
		"11:30",
		"12PM",
		"12:30",
		"1",
		"1:30",
		"2",
		"2:30",
		"3",
		"3:30",
		"4",
		"4:30",
		"5",
		"5:30",
		"6",
		"6:30",
		"7",
		"7:30",
		"8",
		"8:30",
		"9",
		"9:30",
		"10",
		"10:30",
		"11",
		"11:30",
		"12AM",
		"12:30",
		"1",
		"1:30",
		"2",
		"2:30",
		"3",
		"3:30",
		"4",
		"4:30",
	}
	return buf
}
