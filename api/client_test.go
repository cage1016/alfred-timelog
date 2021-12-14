package api

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_GetWeekDay(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	tests := []struct {
		year   int
		month  int
		day    int
		expect []string
	}{
		{2021, 9, 16, []string{"09/13, Monday", "09/14, Tuesday", "09/15, Wednesday", "09/16, Thursday", "09/17, Friday", "09/18, Saturday", "09/19, Sunday"}},
		{2021, 12, 14, []string{"12/13, Monday", "12/14, Tuesday", "12/15, Wednesday", "12/16, Thursday", "12/17, Friday", "12/18, Saturday", "12/19, Sunday"}},
	}

	for _, tt := range tests {
		d := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
		week := GetWeekDay(d, loc)
		assert.DeepEqual(t, tt.expect, week)
	}
}
