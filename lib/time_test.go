package lib

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_GetNewStartEndUnix(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")

	type res struct {
		SheetName string
		Wds       []string
		Start     int64
		End       int64
	}

	type fields struct {
		c func(now time.Time) (string, []string, int64, int64)
	}

	type args struct {
		n map[time.Time]res
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
	}{
		{
			name: "test1",
			prepare: func(f *fields) {
				f.c = GetNewStartEndUnix
			},
			args: args{
				n: map[time.Time]res{
					time.Date(2022, 9, 26, 3, 24, 22, 0, loc): {
						"2022w38",
						[]string{"09/19, Monday", "09/20, Tuesday", "09/21, Wednesday", "09/22, Thursday", "09/23, Friday", "09/24, Saturday", "09/25, Sunday"},
						1663534800,
						1664139600,
					},
					time.Date(2022, 1, 19, 5, 0, 0, 0, loc): {
						"2022w03",
						[]string{"01/17, Monday", "01/18, Tuesday", "01/19, Wednesday", "01/20, Thursday", "01/21, Friday", "01/22, Saturday", "01/23, Sunday"},
						1642366800,
						1642971600,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			for k, v := range tt.args.n {
				sheetName, wds, start, end := f.c(k)
				assert.Equal(t, v.SheetName, sheetName)
				assert.DeepEqual(t, v.Wds, wds)
				assert.Equal(t, v.Start, start)
				assert.Equal(t, v.End, end)
			}
		})
	}
}
