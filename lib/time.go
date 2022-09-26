package lib

import (
	"fmt"
	"time"

	"github.com/snabb/isoweek"
)

func GetNewStartEndUnix(now time.Time) (string, []string, int64, int64) {
	var wds []string
	var wyear, week int

	wyear, week, wds = func(t time.Time) (int, int, []string) {
		wyear, week := isoweek.FromDate(t.Year(), t.Month(), t.Day())
		c := isoweek.StartTime(wyear, week, t.Location())
		if t.Sub(c) < time.Hour*5 {
			week -= 1
			c = isoweek.StartTime(wyear, week, t.Location())
		}

		buf := make([]string, 7)
		for i := 0; i < 7; i++ {
			buf[i] = c.AddDate(0, 0, i).Format("01/02, Monday")
		}
		return wyear, week, buf
	}(now)

	c := isoweek.StartTime(wyear, week, now.Location())
	a := time.Date(c.Year(), c.Month(), c.Day(), 5, 0, 0, 0, now.Location())
	b := a.AddDate(0, 0, 7)
	start := a.Unix()
	end := b.Unix()

	sheetname := fmt.Sprintf("%dw%02d", wyear, week)
	return sheetname, wds, start, end
}
