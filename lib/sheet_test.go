package lib_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/cage1016/alfred-timelog/alfred"
	"github.com/cage1016/alfred-timelog/lib"
	template "github.com/cage1016/alfred-timelog/templates"
	"github.com/xuri/excelize/v2"
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
		id := lib.IdOf(tt.id)
		assert.Equal(t, tt.expect, id)
	}
}

func Test_InitOrUpdate(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")

	type res struct {
		res alfred.Timelog
		err error
	}

	type fields struct {
		x func(now time.Time) (alfred.Timelog, error)
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
			name: "new timelog",
			prepare: func(f *fields) {
				dir, _ := ioutil.TempDir("", "timelog")
				f.x = func(now time.Time) (alfred.Timelog, error) {
					return lib.InitOrUpdate(now, dir)
				}
			},
			args: args{
				n: map[time.Time]res{
					time.Date(2021, 1, 6, 0, 0, 0, 0, loc): {
						res: alfred.Timelog{
							Dir:           "/tmp",
							FileName:      "2021w01.xlsx",
							WeekStartUnix: 1609459200,
							WeekEndUnix:   1610064000,
						},
						err: nil,
					},
				},
			},
		},
		{
			name: "use existing timelog",
			prepare: func(f *fields) {
				dir, _ := ioutil.TempDir("", "timelog")
				te := template.NewEngine()
				data := te.MustAssetString("xlsx/template.xlsx")
				ff, _ := excelize.OpenReader(bytes.NewReader([]byte(data)))
				ff.SaveAs(filepath.Join(dir, "2021w01.xlsx"))

				f.x = func(now time.Time) (alfred.Timelog, error) {
					return lib.InitOrUpdate(now, dir)
				}
			},
			args: args{
				n: map[time.Time]res{
					time.Date(2021, 1, 6, 0, 0, 0, 0, loc): {
						res: alfred.Timelog{
							Dir:           "/tmp",
							FileName:      "2021w01.xlsx",
							WeekStartUnix: 1609459200,
							WeekEndUnix:   1610064000,
						},
						err: nil,
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
				res, err := f.x(k)
				if err != nil {
					assert.Equal(t, v.res.FileName, res.FileName)
					assert.Equal(t, v.res.WeekStartUnix, res.WeekStartUnix)
					assert.Equal(t, v.res.WeekEndUnix, res.WeekEndUnix)
				}
			}
		})
	}
}
