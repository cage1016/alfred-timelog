package utils_test

import (
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/cate1016/alfred-timelog/utils"
)

func Test_IdOf(t *testing.T) {
	tests := []struct {
		id     int
		expect string
	}{
		{0, "A"},
		{26, "AA"},
		{27, "AB"},
		{500, "SG"},
	}

	for _, tt := range tests {
		id := utils.IdOf(tt.id)
		assert.Equal(t, tt.expect, id)
	}
}

func Test_GetWeekDay(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	tests := []struct {
		year   int
		month  int
		day    int
		expect []string
	}{
		{2021, 9, 16, []string{"09/13, Monday", "09/14, Tuesday", "09/15, Wednesday", "09/16, Thursday", "09/17, Friday", "09/18, Saturday", "09/19, Sunday"}},
	}

	for _, tt := range tests {
		d := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
		week := utils.GetWeekDay(d, loc)
		assert.DeepEqual(t, tt.expect, week)
	}
}
